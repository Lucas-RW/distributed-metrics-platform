package services

import (
    "errors"
    "time"
    "github.com/Lucas-RW/distributed-metrics-platform/internal/models"
    "github.com/Lucas-RW/distributed-metrics-platform/internal/storage"
    "github.com/Lucas-RW/distributed-metrics-platform/internal/utils"
    "github.com/Lucas-RW/distributed-metrics-platform/internal/telemetry"
)

var allowedTypes = map[string]bool{
    "counter":   true,
    "gauge":     true,
    "histogram": true,
}

func Ingest(data models.Metric) (result models.Metric, err error) {

	start := time.Now()
	defer func() {
		outcome := "success"
		if err != nil {
			outcome = "error"
		}
		telemetry.IngestDurationSeconds.WithLabelValues(outcome).Observe(time.Since(start).Seconds())
	}()

	if data.Name == "" {
		return models.Metric{}, errors.New("Metric name is required")
	}

	if !allowedTypes[data.Type] {
		return models.Metric{}, errors.New("Metric type is not valid")
	}

	if err = utils.ValidateMetricValue(data.Value); err != nil {
		return models.Metric{}, err
	}

	data.Name = utils.NormalizeMetricName(data.Name)
	data.Labels = utils.NormalizeLabels(data.Labels)
	data.Timestamp = utils.SetTimestampIfMissing(data.Timestamp)

	err = storage.Save(data)
	return data, err
}