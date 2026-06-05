package services

import (
	"context"
	"log"
	"os"
    "math"
    "sync"
    "testing"
    "github.com/Lucas-RW/distributed-metrics-platform/internal/models"
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

func TestIngest_EmptyName_ReturnsError(t *testing.T) {
    metric := models.Metric{Name: "", Type: "counter", Value: 1.0}
    if _, err := Ingest(metric); err == nil {
        t.Error("Expected error for empty name, got nil")
    }
}

func TestIngest_ValidName_NoError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "counter", Value: 1.0}
    if _, err := Ingest(metric); err != nil {
        t.Errorf("Expected no error for valid name, got: %v", err)
    }
}

func TestIngest_InvalidType_ReturnsError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "summary", Value: 1.0}
    if _, err := Ingest(metric); err == nil {
        t.Error("Expected error for invalid type, got nil")
    }
}

func TestIngest_EmptyType_ReturnsError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "", Value: 1.0}
    if _, err := Ingest(metric); err == nil {
        t.Error("Expected error for empty type, got nil")
    }
}

func TestIngest_CounterType_NoError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "counter", Value: 1.0}
    if _, err := Ingest(metric); err != nil {
        t.Errorf("Expected no error for 'counter' type, got: %v", err)
    }
}

func TestIngest_GaugeType_NoError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "gauge", Value: 1.0}
    if _, err := Ingest(metric); err != nil {
        t.Errorf("Expected no error for 'gauge' type, got: %v", err)
    }
}

func TestIngest_HistogramType_NoError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "histogram", Value: 1.0}
    if _, err := Ingest(metric); err != nil {
        t.Errorf("Expected no error for 'histogram' type, got: %v", err)
    }
}

func TestIngest_NaNValue_ReturnsError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "gauge", Value: math.NaN()}
    if _, err := Ingest(metric); err == nil {
        t.Error("Expected error for NaN value, got nil")
    }
}

func TestIngest_PositiveInfValue_ReturnsError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "gauge", Value: math.Inf(1)}
    if _, err := Ingest(metric); err == nil {
        t.Error("Expected error for +Inf value, got nil")
    }
}

func TestIngest_NegativeInfValue_ReturnsError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "gauge", Value: math.Inf(-1)}
    if _, err := Ingest(metric); err == nil {
        t.Error("Expected error for -Inf value, got nil")
    }
}

func TestIngest_ZeroValue_NoError(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "gauge", Value: 0}
    if _, err := Ingest(metric); err != nil {
        t.Errorf("Expected no error for zero value, got: %v", err)
    }
}

func TestIngest_NegativeValue_NoError(t *testing.T) {
    metric := models.Metric{Name: "temperature", Type: "gauge", Value: -12.5}
    if _, err := Ingest(metric); err != nil {
        t.Errorf("Expected no error for negative value, got: %v", err)
    }
}

func TestIngest_NormalisesName(t *testing.T) {
    metric := models.Metric{Name: "  CPU_Usage  ", Type: "gauge", Value: 1.0}
    result, err := Ingest(metric)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Name != "cpu_usage" {
        t.Errorf("Expected normalised name 'cpu_usage', got '%s'", result.Name)
    }
}

func TestIngest_NormalisesLabels(t *testing.T) {
    metric := models.Metric{
        Name:   "cpu_usage",
        Type:   "gauge",
        Value:  1.0,
        Labels: map[string]string{"  Region  ": "  us-west  "},
    }
    result, err := Ingest(metric)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Labels["region"] != "us-west" {
        t.Errorf("Expected normalised label region='us-west', got: %v", result.Labels)
    }
}

func TestIngest_SetsTimestampWhenZero(t *testing.T) {
    metric := models.Metric{Name: "cpu_usage", Type: "gauge", Value: 1.0, Timestamp: 0}
    result, err := Ingest(metric)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Timestamp == 0 {
        t.Error("Expected timestamp to be set, got 0")
    }
}

func TestIngest_PreservesExistingTimestamp(t *testing.T) {
    ts := int64(1700000000)
    metric := models.Metric{Name: "cpu_usage", Type: "gauge", Value: 1.0, Timestamp: ts}
    result, err := Ingest(metric)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if result.Timestamp != ts {
        t.Errorf("Expected timestamp %d, got %d", ts, result.Timestamp)
    }
}

func TestIngest_ConcurrentIngests_DoNotPanic(t *testing.T) {
    var wg sync.WaitGroup
    for i := 0; i < 50; i++ {
        wg.Add(1)
        go func(i int) {
            defer wg.Done()
            _, _ = Ingest(models.Metric{Name: "cpu_usage", Type: "gauge", Value: float64(i)})
        }(i)
    }
    wg.Wait()
}