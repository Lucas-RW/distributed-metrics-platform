package utils

import (
	"math"
	"testing"
	"time"
)

func TestNormalizeMetricName_LowercasesInput(t *testing.T) {
	result := NormalizeMetricName("CPU_Usage")
	if result != "cpu_usage" {
		t.Errorf("Expected 'cpu_usage', got '%s'", result)
	}
}

func TestNormalizeMetricName_TrimsWhitespace(t *testing.T) {
	result := NormalizeMetricName("  cpu_usage  ")
	if result != "cpu_usage" {
		t.Errorf("Expected 'cpu_usage', got '%s'", result)
	}
}

func TestNormalizeMetricName_TrimsAndLowercases(t *testing.T) {
	result := NormalizeMetricName("  MEMORY_USAGE  ")
	if result != "memory_usage" {
		t.Errorf("Expected 'memory_usage', got '%s'", result)
	}
}

func TestNormalizeMetricName_AlreadyNormalized(t *testing.T) {
	result := NormalizeMetricName("cpu_usage")
	if result != "cpu_usage" {
		t.Errorf("Expected 'cpu_usage', got '%s'", result)
	}
}

func TestNormalizeMetricName_EmptyString(t *testing.T) {
	result := NormalizeMetricName("")
	if result != "" {
		t.Errorf("Expected empty string, got '%s'", result)
	}
}

func TestValidateMetricValue_ValidPositiveNumber(t *testing.T) {
	if err := ValidateMetricValue(42.5); err != nil {
		t.Errorf("Expected no error for valid value, got: %v", err)
	}
}

func TestValidateMetricValue_ValidZero(t *testing.T) {
	if err := ValidateMetricValue(0); err != nil {
		t.Errorf("Expected no error for zero, got: %v", err)
	}
}

func TestValidateMetricValue_ValidNegativeNumber(t *testing.T) {
	if err := ValidateMetricValue(-10.0); err != nil {
		t.Errorf("Expected no error for negative value, got: %v", err)
	}
}

func TestValidateMetricValue_NaN(t *testing.T) {
	if err := ValidateMetricValue(math.NaN()); err == nil {
		t.Error("Expected error for NaN, got nil")
	}
}

func TestValidateMetricValue_PositiveInfinity(t *testing.T) {
	if err := ValidateMetricValue(math.Inf(1)); err == nil {
		t.Error("Expected error for +Inf, got nil")
	}
}

func TestValidateMetricValue_NegativeInfinity(t *testing.T) {
	if err := ValidateMetricValue(math.Inf(-1)); err == nil {
		t.Error("Expected error for -Inf, got nil")
	}
}

func TestNormalizeLabels_LowercasesKeys(t *testing.T) {
	input := map[string]string{"Region": "us-west"}
	result := NormalizeLabels(input)
	if _, ok := result["region"]; !ok {
		t.Errorf("Expected key 'region', not found in result: %v", result)
	}
}

func TestNormalizeLabels_TrimsKeyWhitespace(t *testing.T) {
	input := map[string]string{"  region  ": "us-west"}
	result := NormalizeLabels(input)
	if _, ok := result["region"]; !ok {
		t.Errorf("Expected key 'region', not found in result: %v", result)
	}
}

func TestNormalizeLabels_TrimsValueWhitespace(t *testing.T) {
	input := map[string]string{"region": "  us-west  "}
	result := NormalizeLabels(input)
	if result["region"] != "us-west" {
		t.Errorf("Expected value 'us-west', got '%s'", result["region"])
	}
}

func TestNormalizeLabels_PreservesValueCase(t *testing.T) {
	input := map[string]string{"region": "US-West"}
	result := NormalizeLabels(input)
	if result["region"] != "US-West" {
		t.Errorf("Expected value 'US-West' to be unchanged, got '%s'", result["region"])
	}
}

func TestNormalizeLabels_EmptyMap(t *testing.T) {
	result := NormalizeLabels(map[string]string{})
	if len(result) != 0 {
		t.Errorf("Expected empty map, got %v", result)
	}
}

func TestNormalizeLabels_MultipleEntries(t *testing.T) {
	input := map[string]string{
		"  HOST  ": "  server-1  ",
		"ENV":      "  production  ",
	}
	result := NormalizeLabels(input)

	if result["host"] != "server-1" {
		t.Errorf("Expected 'server-1', got '%s'", result["host"])
	}
	if result["env"] != "production" {
		t.Errorf("Expected 'production', got '%s'", result["env"])
	}
}

func TestSetTimestampIfMissing_ReturnsExistingTimestamp(t *testing.T) {
	ts := int64(1700000000)
	result := SetTimestampIfMissing(ts)
	if result != ts {
		t.Errorf("Expected %d, got %d", ts, result)
	}
}

func TestSetTimestampIfMissing_SetsCurrentTimeWhenZero(t *testing.T) {
	before := time.Now().Unix()
	result := SetTimestampIfMissing(0)
	after := time.Now().Unix()

	if result < before || result > after {
		t.Errorf("Expected timestamp between %d and %d, got %d", before, result, after)
	}
}

func TestSetTimestampIfMissing_DoesNotModifyNonZero(t *testing.T) {
	ts := int64(9999999999)
	result := SetTimestampIfMissing(ts)
	if result != ts {
		t.Errorf("Expected %d to be unchanged, got %d", ts, result)
	}
}