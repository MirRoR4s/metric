// Package metric provides a simple and efficient way to create and manage various types of metrics, such as counters and gauges, for monitoring applications.
package metric

import (
	"fmt"
	"sync"

	"github.com/MirRoR4s/metric/internal/metric"
)

// Counter is a metric that represents a single numerical value that only ever goes up.
type Counter struct {
	*metric.Metric
	mu    sync.RWMutex
	count float64
}

// NewCounter creates a new Counter metric with the given name and help description.
//
// It panics if the name is invalid or the help description is empty.
func NewCounter(name, help string) (*Counter, error) {
	metric, err := metric.New(name, help, metric.TypeCounter)
	if err != nil {
		return nil, err
	}
	return &Counter{Metric: metric}, nil
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
func (c *Counter) Add(v float64) error {
	if v < 0 {
		return fmt.Errorf("counter cannot be incremented by a negative value")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.count += v
	return nil
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
	return c.Metric.WritePrometheus() + c.Name + " " + fmt.Sprintf("%d", int64(c.Value())) + "\n"
}

// Gauge is a metric that represents a single numerical value that can arbitrarily go up and down.
type Gauge struct {
	*metric.Metric
	mu    sync.RWMutex
	value float64
}

func NewGauge(name, help string) (*Gauge, error) {
	metric, err := metric.New(name, help, metric.TypeGauge)
	if err != nil {
		return nil, err
	}
	return &Gauge{Metric: metric}, nil
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
	return g.Metric.WritePrometheus() + g.Name + " " + fmt.Sprintf("%f", g.Value()) + "\n"
}
