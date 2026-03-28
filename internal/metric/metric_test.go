package metric

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	metric, err := New("test_metric", "A test metric.", TypeCounter)
	assert.NoError(t, err, "expected no error when creating a new metric")
	assert.Equal(t, "test_metric", metric.Name, "expected metric name to be 'test_metric'")
	assert.Equal(t, "A test metric.", metric.Help, "expected metric help to be 'A test metric.'")
	assert.Equal(t, TypeCounter, metric.typ, "expected metric type to be 'counter'")
}

func TestNewInvalidName(t *testing.T) {
	// Table driven test for invalid metric names
	invalidNames := []string{
		"?aaaa", // starts with a digit
		"",      // empty string
	}
	for _, name := range invalidNames {
		t.Run(name, func(t *testing.T) {
			_, err := New(name, "A test metric.", TypeCounter)
			assert.Error(t, err, "expected an error when creating a metric with an invalid name: "+name)
		})
	}
}

func TestNewEmptyHelp(t *testing.T) {
	_, err := New("test_metric", "", TypeCounter)
	assert.Error(t, err, "expected an error when creating a metric with an empty help description")
}

func TestWritePrometheus(t *testing.T) {
	metric, err := New("test_metric", "A test metric.", TypeCounter)
	assert.NoError(t, err, "expected no error when creating a new metric")
	expectedOutput := "# HELP test_metric A test metric.\n# TYPE test_metric counter\n"
	assert.Equal(t, expectedOutput, metric.WritePrometheus(), "expected Prometheus output to match the expected format")
}
