package telemetry

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var MetricsIngestedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "metrics_ingested_total",
	Help: "Total number of metrics successfully ingested.",
})

var MetricsIngestErrorsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "metrics_ingest_errors_total",
	Help: "Total number of metrics that failed to be ingested.",
})

var IngestDurationSeconds = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "ingest_duration_seconds",
	Help:    "Duration of metric ingestion in seconds.",
	Buckets: prometheus.DefBuckets,
}, []string{"outcome"})