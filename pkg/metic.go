package metric

import (
	"regexp"
	"sync"
)

var (
	// rEMetricName is a regular expression that matches valid Prometheus metric names.
	rEMetricName = regexp.MustCompile("[a-zA-Z_:][a-zA-Z0-9_:]*")
	// rELabelName is a regular expression that matches valid Prometheus label names.
	rELabelName = regexp.MustCompile("[a-zA-Z][a-zA-Z0-9_]*")
)

type Counter struct {
	name   string
	help   string
	mu     sync.RWMutex
	count  float64
	labels map[string]string
}

// NewCounter creates a new Counter metric with the given name and help description.
//
// It panics if the name is invalid or the help description is empty.
func NewCounter(name, help string) *Counter {
	if !rEMetricName.MatchString(name) {
		panic(ErrInvalidMetricName)
	}
	if help == "" {
		panic(ErrEmptyMetricDesc)
	}
	if help[len(help)-1] != '.' {
		help += "."
	}
	return &Counter{name: name, help: help, labels: make(map[string]string)}
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
	info := "# HELP " + c.name + " " + c.help + "\n" +
		"# TYPE " + c.name + " counter\n" 
	for k, v := range c.labels {
		
	}
	
}

type Gauge struct {
	name string
	help string
	mu   sync.RWMutex
	value float64
	labels map[string]string
}

func NewGauge(name, help string) *Gauge {
	if !rEMetricName.MatchString(name) {
		panic(ErrInvalidMetricName)
	}
	if help == "" {
		panic(ErrEmptyMetricDesc)
	}
	return &Gauge{name: name, help: help}
}
