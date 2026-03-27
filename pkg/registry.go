package metric

import "net/http"

type PrometheusFormatter interface {
	WritePrometheus() string
}

type Registry struct {

}

func (r *Registry) Register(metric PrometheusFormatter) {

}

func (r *Registry) Handler() http.Handler {
	panic(1)
}
