package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	requestsTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "Total number of /process requests",
		})
	workerURL string
)

func main() {
	workerURL = os.Getenv("WORKER_URL")
	if workerURL == "" {
		log.Fatal("WORKER_URL is not set")
	}

	prometheus.MustRegister(requestsTotal)

	http.HandleFunc("/process", handleProcess)
	http.Handle("/metrics", promhttp.Handler())

	log.Println("API Service listening on :8080")
	log.Printf("Worker URL: %s\n", workerURL)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	requestsTotal.Inc()
	log.Println("Processing /process request...")
	// Call worker-service using WORKER_URL env var
	resp, err := http.Post(fmt.Sprintf("%s/do-work", workerURL), "application/json", bytes.NewBuffer([]byte(`{}`)))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error calling worker: %v", err), http.StatusInternalServerError)
		return
	}
	body, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	w.Write(body)
}
