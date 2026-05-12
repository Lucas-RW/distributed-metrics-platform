package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/models"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/services"
)

func MetricsHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var metric models.Metric

	err := json.NewDecoder(r.Body).Decode(&metric)

	if err != nil {
		http.Error(w, "Invalid metric data", http.StatusBadRequest)
		return
	}

	if metric.Name == "" {
		http.Error(w, "Metric name is required", http.StatusBadRequest)
		return
	}

	err = services.Ingest(metric)

	if err != nil {
		http.Error(w, "Failed to ingest metric", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metric ingested successfully"))

}