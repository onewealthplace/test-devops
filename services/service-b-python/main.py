import os
import grpc
from fastapi import FastAPI, HTTPException
from opentelemetry import trace
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor
from opentelemetry.instrumentation.grpc import GrpcInstrumentorClient
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.trace.export import BatchSpanProcessor

# -----------------------------------------------------------------------------
# OpenTelemetry setup
# -----------------------------------------------------------------------------
resource = Resource.create({"service.name": "service-b-python"})
trace.set_tracer_provider(TracerProvider(resource=resource))
span_processor = BatchSpanProcessor(
    OTLPSpanExporter(
        endpoint=os.getenv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318/v1/traces"),
        insecure=True,
    )
)
trace.get_tracer_provider().add_span_processor(span_processor)

# Instrument gRPC client & FastAPI
GrpcInstrumentorClient().instrument()

app = FastAPI(title="Service B (Python)")
FastAPIInstrumentor.instrument_app(app)

# -----------------------------------------------------------------------------
# gRPC client configuration
# -----------------------------------------------------------------------------
SERVICE_A_GRPC_ADDR = os.getenv("SERVICE_A_GRPC_ADDR", "service-a-go:50051")

# Try importing generated stubs. They will be provided later when the proto is
# added. Until then, fail gracefully so the service can still start.
try:
    from proto import service_a_pb2, service_a_pb2_grpc  # type: ignore
    grpc_available = True
except ImportError:
    grpc_available = False

aio_channel_type = getattr(grpc, "aio", None)

@app.get("/call-a")
async def call_a():
    """Invoke service-a via gRPC and return the response."""
    if not grpc_available or aio_channel_type is None:
        raise HTTPException(status_code=500, detail="gRPC stubs not available yet.")

    async with aio_channel_type.insecure_channel(SERVICE_A_GRPC_ADDR) as channel:  # type: ignore[attr-defined]
        stub = service_a_pb2_grpc.ServiceAStub(channel)  # type: ignore[name-defined]
        try:
            response = await stub.Ping(service_a_pb2.PingRequest())  # type: ignore[name-defined]
            return {"message": response.message}
        except grpc.aio.AioRpcError as exc:  # type: ignore[attr-defined]
            raise HTTPException(status_code=500, detail=f"gRPC error: {exc}") from exc


if __name__ == "__main__":
    import uvicorn

    uvicorn.run("main:app", host="0.0.0.0", port=int(os.getenv("PORT", 8000)), reload=False)
