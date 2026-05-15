package utils

import "strings"

func NormalizeMetricName(name string) string { 
	return strings.ToLower(strings.TrimSpace(name))
}

func NormalizeMetricValue(value string) string {
	return strings.TrimSpace(value)
}

func NormalizeLabels(labels map[string]string) map[string]string {
	normalized := make(map[string]string)
	for key, value := range labels {
		normalized[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(value)
	}
	return normalized
}