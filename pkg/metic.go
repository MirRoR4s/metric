package metric

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v4/mem"
)

var (
	// rEMetricName is a regular expression that matches valid Prometheus metric names.
	rEMetricName = regexp.MustCompile("[a-zA-Z_:][a-zA-Z0-9_:]*")
)

type Type string

const (
	TypeCounter Type = "counter"
	TypeGauge   Type = "gauge"
	TypeUntyped Type = "untyped"
)

// Metric is the base struct for all metrics, containing common fields like name and help description.
type Metric struct {
	name string
	help string
	typ  Type
}

func (m *Metric) WritePrometheus() string {
	helpInfo := "# HELP " + m.name + " " + m.help + "\n"
	typeInfo := "# TYPE " + m.name + " " + string(m.typ) + "\n"
	return helpInfo + typeInfo
}

func NewMetric(name, help string, typ Type) *Metric {
	if !rEMetricName.MatchString(name) {
		panic(ErrInvalidMetricName)
	}
	if help == "" {
		panic(ErrEmptyMetricDesc)
	}
	if help[len(help)-1] != '.' {
		help += "."
	}
	return &Metric{name: name, help: help, typ: typ}
}

// Counter is a metric that represents a single numerical value that only ever goes up.
type Counter struct {
	*Metric
	mu    sync.RWMutex
	count float64
}

// NewCounter creates a new Counter metric with the given name and help description.
//
// It panics if the name is invalid or the help description is empty.
func NewCounter(name, help string) *Counter {
	return &Counter{Metric: NewMetric(name, help, TypeCounter)}
}

// Inc increments the counter by 1.
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count++
}

// Add increments the counter by the given value.
//
// If the value is negative, it panics because counters cannot be decremented.
func (c *Counter) Add(v float64) {
	if v < 0 {
		panic("counter cannot be incremented by a negative value")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count += v
}

// Value returns the current value of the counter.
func (c *Counter) Value() float64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.count
}

func (c *Counter) WritePrometheus() string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.Metric.WritePrometheus() + c.name + " " + fmt.Sprintf("%d", int64(c.Value())) + "\n"
}

// HttpRequestsTotal is a helper function that creates a Counter metric for tracking total HTTP requests and returns it along with a middleware function that increments the counter for each incoming HTTP request.
//
// For example, you can use it like this:
//
//	counter, middleware := metric.HttpRequestsTotal()
//	registry.Register(counter)
//	http.Handle("/metrics", registry.Handler())
//	http.Handle("/hello", middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("Hello, World!"))
//	})))
func HttpRequestsTotal() (*Counter, func(http.Handler) http.Handler) {
	counter := NewCounter("http_requests_total", "Total number of HTTP requests.")
	middleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			counter.Inc()
			next.ServeHTTP(w, r)
		})
	}
	return counter, middleware
}

// Gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
type Gauge struct {
	*Metric
	mu     sync.RWMutex
	value  float64
	labels map[string]string
}

func NewGauge(name, help string) *Gauge {
	if !rEMetricName.MatchString(name) {
		panic(ErrInvalidMetricName)
	}
	if help == "" {
		panic(ErrEmptyMetricDesc)
	}
	if help[len(help)-1] != '.' {
		help += "."
	}
	return &Gauge{Metric: NewMetric(name, help, TypeGauge), labels: make(map[string]string)}
}

// Set sets the gauge to the given value.
func (g *Gauge) Set(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = v
}

// Inc increments the gauge by 1.
func (g *Gauge) Inc() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value++
}

// Dec decrements the gauge by 1.
func (g *Gauge) Dec() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value--
}

// Add increments the gauge by the given value.
func (g *Gauge) Add(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += v
}

// Sub decrements the gauge by the given value.
func (g *Gauge) Value() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

func (g *Gauge) WritePrometheus() string {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.Metric.WritePrometheus() + g.name + " " + fmt.Sprintf("%f", g.Value()) + "\n"
}

// memory holds memory-related gauges and updates them periodically.
type memory struct {
	VirtualMemoryBytes    *Gauge
	VirtualMemoryMaxBytes *Gauge
	ticker                *time.Ticker
	cancel                context.CancelFunc
}

// Memory creates and starts memory metrics collection.
func Memory(ctx context.Context) *memory {
	ctx, cancel := context.WithCancel(ctx)
	mm := &memory{
		VirtualMemoryBytes:    NewGauge("process_virtual_memory_bytes", "Virtual memory size in bytes."),
		VirtualMemoryMaxBytes: NewGauge("process_virtual_memory_max_bytes", "Maximum amount of virtual memory available in bytes."),
		ticker:                time.NewTicker(500 * time.Millisecond),
		cancel:                cancel,
	}

	go mm.update(ctx)
	return mm
}

// update periodically updates all memory metrics.
func (mm *memory) update(ctx context.Context) {
	defer mm.ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Println("Stopping memory metrics updater")
			return
		case <-mm.ticker.C:
			v, err := mem.VirtualMemory()
			if err != nil {
				log.Printf("Error getting virtual memory: %v", err)
				continue
			}
			mm.VirtualMemoryBytes.Set(float64(v.Total))
			mm.VirtualMemoryMaxBytes.Set(float64(v.Available))
		}
	}
}

// Stop stops the memory metrics collection.
func (mm *memory) Stop() {
	mm.cancel()
}

func (mm *memory) WritePrometheus() string {
	return mm.VirtualMemoryBytes.WritePrometheus() + mm.VirtualMemoryMaxBytes.WritePrometheus()
}
