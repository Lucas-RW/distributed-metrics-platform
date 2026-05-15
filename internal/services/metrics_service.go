package services

import (
	"errors"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/models"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/storage"
	"github.com/Lucas-RW/distributed-metrics-platform/internal/utils"
)

var allowedTypes = map[string]bool{
	"counter":   true,
	"gauge":     true,
	"histogram": true,
}

func Ingest(data models.Metric) error {
	if data.Name == "" {
		return errors.New("Metric name is required")
	}

	if !allowedTypes[data.Type] {
		return errors.New("Metric type is not valid")
	}

	if err := utils.ValidateMetricValue(data.Value); err != nil {
		return err
	}

	data.Name = utils.NormalizeMetricName(data.Name)
	data.Labels = utils.NormalizeLabels(data.Labels)
	data.Timestamp = utils.SetTimestampIfMissing(data.Timestamp)

	return storage.Save(data)
}
