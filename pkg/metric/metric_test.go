package metric

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	counter, err := NewCounter("test_counter", "A test counter.")
	assert.NoError(t, err, "expected no error when creating a new counter")
	assert.Equal(t, "test_counter", counter.Name, "expected counter name to be 'test_counter'")
	assert.Equal(t, "A test counter.", counter.Help, "expected counter help to be 'A test counter.'")
}

func TestNewCounterEmptyName(t *testing.T) {
	// Testing that creating a counter with an empty name will error
	_, err := NewCounter("", "A test counter with empty name.")
	assert.Error(t, err, "expected error when creating a counter with an empty name")
}

func TestNewCounterEmptyHelp(t *testing.T) {
	// Testing that creating a counter with an empty help description will error
	_, err := NewCounter("test_counter_empty_help", "")
	assert.Error(t, err, "expected error when creating a counter with an empty help description")
}

func TestCounterInc(t *testing.T) {
	counter, err := NewCounter("test_counter", "A test counter.")
	assert.NoError(t, err, "expected no error when creating a new counter")
	counter.Inc()
	assert.Equal(t, 1.0, counter.Value(), "expected counter value 1 after one increment")

	var wg sync.WaitGroup
	// Testing thread safety by incrementing the counter from multiple goroutines.
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Inc()
		}()
	}
	wg.Wait()
	assert.Equal(t, 101.0, counter.Value(), "expected counter value 101 after 100 increments")
}

func TestCounterAdd(t *testing.T) {
	counter, err := NewCounter("test_counter_add", "A test counter for Add method.")
	assert.NoError(t, err, "expected no error when creating a new counter")
	assert.NoError(t, counter.Add(5))
	assert.Equal(t, 5.0, counter.Value(), "expected counter value 5 after adding 5")

	// Testing thread safety by adding to the counter from multiple goroutines.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, counter.Add(1))
		}()
	}
	wg.Wait()
	assert.Equal(t, 105.0, counter.Value(), "expected counter value 105 after 100 additions")

	// Testing that adding a negative value will return an error
	err = counter.Add(-1)
	assert.Error(t, err, "expected error when adding a negative value to the counter")
}

func TestCounterValue(t *testing.T) {
	counter, err := NewCounter("test_counter_value", "A test counter for Value method.")
	assert.NoError(t, err, "expected no error when creating a new counter")
	assert.Equal(t, 0.0, counter.Value(), "expected initial counter value to be 0")

	counter.Inc()
	assert.Equal(t, 1.0, counter.Value(), "expected counter value to be 1 after one increment")
}

func TestCounterWritePrometheus(t *testing.T) {
	counter, err := NewCounter("test_counter_prometheus", "A test counter for Prometheus output.")
	assert.NoError(t, err, "expected no error when creating a new counter")
	expectedOutput := "# HELP test_counter_prometheus A test counter for Prometheus output.\n# TYPE test_counter_prometheus counter\ntest_counter_prometheus 0.000000\n"
	assert.Equal(t, expectedOutput, counter.WritePrometheus(), "expected Prometheus output to match the expected format")
}

func TestNewGauge(t *testing.T) {
	gauge, err := NewGauge("test_gauge", "A test gauge.")
	assert.NoError(t, err, "expected no error when creating a new gauge")
	assert.Equal(t, "test_gauge", gauge.Name, "expected gauge name to be 'test_gauge'")
	assert.Equal(t, "A test gauge.", gauge.Help, "expected gauge help to be 'A test gauge.'")
}

func TestGaugeSet(t *testing.T) {
	gauge, err := NewGauge("test_gauge_set", "A test gauge for Set method.")
	assert.NoError(t, err, "expected no error when creating a new gauge")
	gauge.Set(3.14)
	assert.Equal(t, 3.14, gauge.Value(), "expected gauge value to be 3.14 after setting it")
}

func TestGaugeInc(t *testing.T) {
	gauge, err := NewGauge("test_gauge_inc", "A test gauge for Inc method.")
	assert.NoError(t, err, "expected no error when creating a new gauge")
	gauge.Inc()
	assert.Equal(t, 1.0, gauge.Value(), "expected gauge value to be 1 after one increment")

	// Testing thread safety by incrementing the gauge from multiple goroutines.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gauge.Inc()
		}()
	}
	wg.Wait()
	assert.Equal(t, 101.0, gauge.Value(), "expected gauge value 101 after 100 increments")

}

func TestGaugeDec(t *testing.T) {
	gauge, err := NewGauge("test_gauge_dec", "A test gauge for Dec method.")
	assert.NoError(t, err, "expected no error when creating a new gauge")
	gauge.Dec()
	assert.Equal(t, -1.0, gauge.Value(), "expected gauge value to be -1 after one decrement")

	// Testing thread safety by decrementing the gauge from multiple goroutines.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gauge.Dec()
		}()
	}
	wg.Wait()
	assert.Equal(t, -101.0, gauge.Value(), "expected gauge value -101 after 100 decrements")
}

func TestGaugeAdd(t *testing.T) {
	gauge, err := NewGauge("test_gauge_add", "A test gauge for Add method.")
	assert.NoError(t, err, "expected no error when creating a new gauge")
	gauge.Add(2.5)
	assert.Equal(t, 2.5, gauge.Value(), "expected gauge value to be 2.5 after adding 2.5")

	// Testing thread safety by adding to the gauge from multiple goroutines.
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			gauge.Add(1)
		}()
	}
	wg.Wait()
	assert.Equal(t, 102.5, gauge.Value(), "expected gauge value 102.5 after 100 additions")
}

func TestGaugeWritePrometheus(t *testing.T) {
	gauge, err := NewGauge("test_gauge_prometheus", "A test gauge for Prometheus output.")
	assert.NoError(t, err, "expected no error when creating a new gauge")
	expectedOutput := "# HELP test_gauge_prometheus A test gauge for Prometheus output.\n# TYPE test_gauge_prometheus gauge\ntest_gauge_prometheus 0.000000\n"
	assert.Equal(t, expectedOutput, gauge.WritePrometheus(), "expected Prometheus output to match the expected format")
}
