package metric

import (
	"errors"
	"regexp"
)

type Type string

const (
	TypeCounter Type = "counter"
	TypeGauge   Type = "gauge"
	TypeUntyped Type = "untyped"
)

var (
	// rEMetricName is a regular expression that matches valid Prometheus metric names.
	rEMetricName = regexp.MustCompile("^[a-zA-Z_:][a-zA-Z0-9_:]*$")
)

// Metric is the base struct for all metrics, containing common fields like name and help description.
type Metric struct {
	Name string
	Help string
	typ  Type
}

func (m *Metric) WritePrometheus() string {
	helpInfo := "# HELP " + m.Name + " " + m.Help + "\n"
	typeInfo := "# TYPE " + m.Name + " " + string(m.typ) + "\n"
	return helpInfo + typeInfo
}

func New(name, help string, typ Type) (*Metric, error) {
	if !rEMetricName.MatchString(name) {
		return nil, errors.New("metric name must match the regex [a-zA-Z_:][a-zA-Z0-9_:]*")
	}
	if help == "" {
		return nil, errors.New("metric description cannot be empty")
	}
	if help[len(help)-1] != '.' {
		help += "."
	}
	return &Metric{Name: name, Help: help, typ: typ}, nil
}
