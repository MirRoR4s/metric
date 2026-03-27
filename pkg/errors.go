package metric

import "errors"

var (
	ErrInvalidMetricName = errors.New("invalid metric name")
	ErrInvalidLabelName  = errors.New("invalid label name")
	ErrEmptyMetricDesc   = errors.New("metric description cannot be empty")
)
