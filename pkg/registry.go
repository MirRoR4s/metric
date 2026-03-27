package metric

import "net/http"

// PrometheusFormatter is an interface that defines a method for writing metrics in Prometheus text format.
type PrometheusFormatter interface {
	// WritePrometheus returns the metric in Prometheus text format.
	WritePrometheus() string
}

type Registry struct {
	metircs []PrometheusFormatter
}

func NewRegistry() *Registry {
	return &Registry{
		metircs: make([]PrometheusFormatter, 0),
	}
}

func (r *Registry) Register(metrics ...PrometheusFormatter) {
	r.metircs = append(r.metircs, metrics...)
}

// Handler returns an HTTP handler that serves the metrics in Prometheus based-text format.
func (r *Registry) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		for _, m := range r.metircs {
			w.Write([]byte(m.WritePrometheus()))
		}
	})
}
