package utils

import (
	"errors"
	"math"
	"strings"
	"time"
)

func NormalizeMetricName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func ValidateMetricValue(value float64) error {
	if math.IsNaN(value) || math.IsInf(value, 0) {
		return errors.New("Metric value must be a finite number")
	}
	return nil
}

func NormalizeLabels(labels map[string]string) map[string]string {
	normalized := make(map[string]string)
	for key, value := range labels {
		normalized[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
	}
	return normalized
}

func SetTimestampIfMissing(timestamp int64) int64 {
	if timestamp == 0 {
		return time.Now().Unix()
	}
	return timestamp
}
