package metric

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCounter(t *testing.T) {
	counter := NewCounter("test_counter", "A test counter.")
	assert.Equal(t, "test_counter", counter.name, "expected counter name to be 'test_counter'")
	assert.Equal(t, "A test counter.", counter.help, "expected counter help to be 'A test counter.'")
	assert.Equal(t, TypeCounter, counter.typ, "expected counter type to be 'counter'")
}

func TestNewCounterEmptyName(t *testing.T) {
	// Testing that creating a counter with an empty name will panic
	assert.Panics(t, func() {
		NewCounter("", "A counter with an empty name.")
	}, "expected panic when creating a counter with an empty name")
}

func TestNewCounterEmptyHelp(t *testing.T) {
	// Testing that creating a counter with an empty help description will panic
	assert.Panics(t, func() {
		NewCounter("test_counter_empty_help", "")
	}, "expected panic when creating a counter with an empty help description")
}

func TestCounterInc(t *testing.T) {
	counter := NewCounter("test_counter", "A test counter.")
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
	counter := NewCounter("test_counter_add", "A test counter for Add method.")
	counter.Add(5)
	assert.Equal(t, 5.0, counter.Value(), "expected counter value 5 after adding 5")

	// Testing thread safety by adding to the counter from multiple goroutines.
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add(1)
		}()
	}
	wg.Wait()
	assert.Equal(t, 105.0, counter.Value(), "expected counter value 105 after 100 additions")

	// Testing that adding a negative value will panic
	assert.Panics(t, func() {
		counter.Add(-1)
	}, "expected panic when adding a negative value to the counter")
}

func TestCounterValue(t *testing.T) {
	counter := NewCounter("test_counter_value", "A test counter for Value method.")
	assert.Equal(t, 0.0, counter.Value(), "expected initial counter value to be 0")

	counter.Inc()
	assert.Equal(t, 1.0, counter.Value(), "expected counter value to be 1 after one increment")
}

func TestCounterWritePrometheus(t *testing.T) {
	counter := NewCounter("test_counter_prometheus", "A test counter for Prometheus output.")
	expectedOutput := "# HELP test_counter_prometheus A test counter for Prometheus output.\n# TYPE test_counter_prometheus counter\ntest_counter_prometheus 0\n"
	assert.Equal(t, expectedOutput, counter.WritePrometheus(), "expected Prometheus output to match the expected format")
}

func TestNewGauge(t *testing.T) {
	gauge := NewGauge("test_gauge", "A test gauge.")
	assert.Equal(t, "test_gauge", gauge.name, "expected gauge name to be 'test_gauge'")
	assert.Equal(t, "A test gauge.", gauge.help, "expected gauge help to be 'A test gauge.'")
	assert.Equal(t, TypeGauge, gauge.typ, "expected gauge type to be 'gauge'")
}

func TestGaugeSet(t *testing.T) {
	gauge := NewGauge("test_gauge_set", "A test gauge for Set method.")
	gauge.Set(3.14)
	assert.Equal(t, 3.14, gauge.Value(), "expected gauge value to be 3.14 after setting it")
}

func TestGaugeInc(t *testing.T) {
	gauge := NewGauge("test_gauge_inc", "A test gauge for Inc method.")
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
	gauge := NewGauge("test_gauge_dec", "A test gauge for Dec method.")
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
	gauge := NewGauge("test_gauge_add", "A test gauge for Add method.")
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
	gauge := NewGauge("test_gauge_prometheus", "A test gauge for Prometheus output.")
	expectedOutput := "# HELP test_gauge_prometheus A test gauge for Prometheus output.\n# TYPE test_gauge_prometheus gauge\ntest_gauge_prometheus 0.000000\n"
	assert.Equal(t, expectedOutput, gauge.WritePrometheus(), "expected Prometheus output to match the expected format")
}

func TestHttpRequestsTotal(t *testing.T) {
	counter, middleware := HttpRequestsTotal()
	// Simulate an HTTP request by calling the middleware function.
	mux := middleware(&http.ServeMux{})
	// Using httptest package to create a test HTTP request and response recorder.
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, 1.0, counter.Value(), "expected counter value to be 1 after one HTTP request")
}

func TestMemory(t *testing.T) {
	memory := Memory(context.Background())
	assert.NotNil(t, memory, "expected memory to be initialized")

}

func TestMemoryUpdate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	memory := Memory(ctx)
	// Allow some time for the memory metrics to be updated.
	time.Sleep(1 * time.Second)
	assert.Greater(t, memory.VirtualMemoryBytes.Value(), 0.0, "expected VirtualMemoryBytes to be greater than 0")
	assert.Greater(t, memory.VirtualMemoryMaxBytes.Value(), 0.0, "expected VirtualMemoryMaxBytes to be greater than 0")
	cancel() // Stop the memory update goroutine
}
