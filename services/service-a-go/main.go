package main

import (
	"context"
	cryptorand "crypto/rand"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"os"
	"time"

	serviceapb "github.com/onewealthplace/test-devops/proto/service_a"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// memLeak holds allocated byte slices to simulate a memory leak.
var memLeak [][]byte

// serviceAServer implements serviceapb.ServiceAServer.
type serviceAServer struct {
	serviceapb.UnimplementedServiceAServer
}

// Ping implements the Ping RPC.
func (s *serviceAServer) Ping(ctx context.Context, _ *serviceapb.PingRequest) (*serviceapb.PingResponse, error) {
	// Add random latency similar to HTTP handler for observability.
	delay := time.Duration(100+rand.Intn(200)) * time.Millisecond
	time.Sleep(delay)
	return &serviceapb.PingResponse{Message: "pong from service-a-go"}, nil
}

func initTracer(ctx context.Context, serviceName string) (*sdktrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
			attribute.String("service.language", "go"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	return tp, nil
}

func main() {

	ctx := context.Background()
	tp, err := initTracer(ctx, "service-a")
	if err != nil {
		panic(fmt.Sprintf("failed to initialize tracer: %v", err))
	}
	defer func() { _ = tp.Shutdown(ctx) }()

	// start gRPC server concurrently
	go func() {
		if err := startGRPCServer(ctx, "50051"); err != nil {
			panic(fmt.Sprintf("failed to start gRPC server: %v", err))
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/process", otelhttp.NewHandler(http.HandlerFunc(processHandler), "process"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("server error: %v", err))
	}
}

func startGRPCServer(ctx context.Context, port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	serviceapb.RegisterServiceAServer(srv, &serviceAServer{})

	go func() {
		<-ctx.Done()
		srv.GracefulStop()
	}()
	return srv.Serve(lis)
}

func processHandler(w http.ResponseWriter, r *http.Request) {
	// Add random latency between 100-300 ms
	delay := time.Duration(100+rand.Intn(200)) * time.Millisecond
	time.Sleep(delay)

	// --- simulate a memory leak: keep 1 MB of random data per request ---
	leak := make([]byte, 1024*1024)
	_, _ = cryptorand.Read(leak) // fill with cryptographically secure random data
	memLeak = append(memLeak, leak)

	// Add OTEL span event with latency and leak size
	if span := trace.SpanFromContext(r.Context()); span.IsRecording() {
		span.AddEvent("processing")
	}

}
