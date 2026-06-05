package handlers

import (
    "bytes"
	"context"
	"log"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"os"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/storage"
)

func TestMain(m *testing.M) {
    ctx := context.Background()
    if err := storage.InitDB(ctx); err != nil {
        log.Fatalf("could not connect to database: %v", err)
    }
    defer storage.CloseDB()
    os.Exit(m.Run())
}

func newPostRequest(t *testing.T, body any) *http.Request {
    t.Helper()
    b, err := json.Marshal(body)
    if err != nil {
        t.Fatalf("Failed to marshal request body: %v", err)
    }
    req, err := http.NewRequest(http.MethodPost, "/metrics", bytes.NewBuffer(b))
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Content-Type", "application/json")
    return req
}

func TestMetricsHandler_GetRequest_Returns405(t *testing.T) {
    req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
    rr := httptest.NewRecorder()
    MetricsHandler(rr, req)
    if rr.Code != http.StatusMethodNotAllowed {
        t.Errorf("Expected 405, got %d", rr.Code)
    }
}

func TestMetricsHandler_PutRequest_Returns405(t *testing.T) {
    req, _ := http.NewRequest(http.MethodPut, "/metrics", nil)
    rr := httptest.NewRecorder()
    MetricsHandler(rr, req)
    if rr.Code != http.StatusMethodNotAllowed {
        t.Errorf("Expected 405, got %d", rr.Code)
    }
}

func TestMetricsHandler_DeleteRequest_Returns405(t *testing.T) {
    req, _ := http.NewRequest(http.MethodDelete, "/metrics", nil)
    rr := httptest.NewRecorder()
    MetricsHandler(rr, req)
    if rr.Code != http.StatusMethodNotAllowed {
        t.Errorf("Expected 405, got %d", rr.Code)
    }
}

func TestMetricsHandler_MalformedJSON_Returns400(t *testing.T) {
    req, _ := http.NewRequest(http.MethodPost, "/metrics", bytes.NewBufferString("{not valid json}"))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    MetricsHandler(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Errorf("Expected 400 for malformed JSON, got %d", rr.Code)
    }
}

func TestMetricsHandler_EmptyBody_Returns400(t *testing.T) {
    req, _ := http.NewRequest(http.MethodPost, "/metrics", bytes.NewBufferString(""))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()
    MetricsHandler(rr, req)
    if rr.Code != http.StatusBadRequest {
        t.Errorf("Expected 400 for empty body, got %d", rr.Code)
    }
}

func TestMetricsHandler_MissingName_Returns400(t *testing.T) {
    body := map[string]any{"Type": "gauge", "Value": 1.0}
    rr := httptest.NewRecorder()
    MetricsHandler(rr, newPostRequest(t, body))
    if rr.Code != http.StatusBadRequest {
        t.Errorf("Expected 400 for missing name, got %d", rr.Code)
    }
}

func TestMetricsHandler_EmptyName_Returns400(t *testing.T) {
    body := map[string]any{"Name": "", "Type": "gauge", "Value": 1.0}
    rr := httptest.NewRecorder()
    MetricsHandler(rr, newPostRequest(t, body))
    if rr.Code != http.StatusBadRequest {
        t.Errorf("Expected 400 for empty name, got %d", rr.Code)
    }
}

func TestMetricsHandler_InvalidType_Returns500(t *testing.T) {
    body := map[string]any{"Name": "cpu_usage", "Type": "summary", "Value": 1.0}
    rr := httptest.NewRecorder()
    MetricsHandler(rr, newPostRequest(t, body))
    if rr.Code != http.StatusInternalServerError {
        t.Errorf("Expected 500 for invalid metric type, got %d", rr.Code)
    }
}

func TestMetricsHandler_ValidMetric_Returns200(t *testing.T) {
    body := map[string]any{"Name": "cpu_usage", "Type": "gauge", "Value": 42.0}
    rr := httptest.NewRecorder()
    MetricsHandler(rr, newPostRequest(t, body))
    if rr.Code != http.StatusOK {
        t.Errorf("Expected 200 for valid metric, got %d", rr.Code)
    }
}

func TestMetricsHandler_ValidMetric_ReturnsSuccessBody(t *testing.T) {
    body := map[string]any{"Name": "cpu_usage", "Type": "gauge", "Value": 42.0}
    rr := httptest.NewRecorder()
    MetricsHandler(rr, newPostRequest(t, body))
    if rr.Body.String() != "Metric ingested successfully" {
        t.Errorf("Unexpected response body: %q", rr.Body.String())
    }
}

func TestMetricsHandler_ValidMetric_WithLabels_Returns200(t *testing.T) {
    body := map[string]any{
        "Name":   "cpu_usage",
        "Type":   "gauge",
        "Value":  42.0,
        "Labels": map[string]string{"region": "us-west"},
    }
    rr := httptest.NewRecorder()
    MetricsHandler(rr, newPostRequest(t, body))
    if rr.Code != http.StatusOK {
        t.Errorf("Expected 200 for metric with labels, got %d", rr.Code)
    }
}

func TestMetricsHandler_ValidMetric_WithTimestamp_Returns200(t *testing.T) {
    body := map[string]any{
        "Name":      "cpu_usage",
        "Type":      "counter",
        "Value":     1.0,
        "Timestamp": 1700000000,
    }
    rr := httptest.NewRecorder()
    MetricsHandler(rr, newPostRequest(t, body))
    if rr.Code != http.StatusOK {
        t.Errorf("Expected 200 for metric with timestamp, got %d", rr.Code)
    }
}

func TestMetricsHandler_ValidMetric_AllThreeTypes(t *testing.T) {
    types := []string{"counter", "gauge", "histogram"}
    for _, metricType := range types {
        body := map[string]any{"Name": "cpu_usage", "Type": metricType, "Value": 1.0}
        rr := httptest.NewRecorder()
        MetricsHandler(rr, newPostRequest(t, body))
        if rr.Code != http.StatusOK {
            t.Errorf("Expected 200 for type '%s', got %d", metricType, rr.Code)
        }
    }
}