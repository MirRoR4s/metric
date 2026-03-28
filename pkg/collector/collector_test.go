package collector

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHttpRequestTotal(t *testing.T) {
	counter, middleware, err := NewHttpRequestsTotal()
	assert.NoError(t, err, "expected no error when creating HTTP requests total counter")
	// Simulate an HTTP request by calling the middleware function.
	mux := middleware(&http.ServeMux{})
	// Using httptest package to create a test HTTP request and response recorder.
	req := httptest.NewRequest("GET", "http://example.com/test", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	assert.Equal(t, 1.0, counter.Value(), "expected counter value to be 1 after one HTTP request")
}

func TestNewProcess(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	memory, err := NewProcess(ctx) // Example: 1GB max virtual memory
	assert.NoError(t, err, "expected no error when creating memory metrics")
	// Allow some time for the memory metrics to be updated.
	time.Sleep(1 * time.Second)
	assert.Greater(t, memory.virtualMemory.Value(), 0.0, "expected VirtualMemory to be greater than 0")
	assert.Greater(t, memory.virtualMemoryMax.Value(), 0.0, "expected VirtualMemoryMax to be greater than 0")
	cancel() // Stop the memory update goroutine
}

func TestNewProcessWithOptions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	memory, err := NewProcess(ctx, WithMaxVirtualMemory(1*1024*1024*1024)) // Example: 1GB max virtual memory
	assert.NoError(t, err, "expected no error when creating memory metrics with options")
	// Allow some time for the memory metrics to be updated.
	time.Sleep(1 * time.Second)
	assert.Greater(t, memory.virtualMemory.Value(), 0.0, "expected VirtualMemory to be greater than 0")
	assert.Equal(t, 1*1024*1024*1024, int(memory.virtualMemoryMax.Value()), "expected VirtualMemoryMax to be 1GB")
	cancel() // Stop the memory update goroutine
}
