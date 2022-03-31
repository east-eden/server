package utils

import (
	"context"
	"fmt"

	"github.com/asim/go-micro/v3/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	defaultMetricPrefix = "micro"
)

func NewPrometheusHandlerWrapper(opts ...server.Option) server.HandlerWrapper {
	md := make(map[string]string)
	sopts := server.Options{}

	for _, opt := range opts {
		opt(&sopts)
	}

	for k, v := range sopts.Metadata {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, k)] = v
	}
	if len(sopts.Name) > 0 {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, "name")] = sopts.Name
	}
	if len(sopts.Id) > 0 {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, "id")] = sopts.Id
	}
	if len(sopts.Version) > 0 {
		md[fmt.Sprintf("%s_%s", defaultMetricPrefix, "version")] = sopts.Version
	}

	opsCounter := promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "micro",
			Name:      "request_total",
			Help:      "How many go-micro requests processed, partitioned by method and status",
		},
		[]string{"method", "status"},
	)

	timeCounterSummary := promauto.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "micro",
			Name:      "upstream_latency_microseconds",
			Help:      "Service backend method request latencies in microseconds",
		},
		[]string{"method"},
	)

	timeCounterHistogram := promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "micro",
			Name:      "request_duration_seconds",
			Help:      "Service method request time in seconds",
		},
		[]string{"method"},
	)

	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp any) error {
			name := req.Endpoint()

			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				timeCounterSummary.WithLabelValues(name).Observe(v)
				timeCounterHistogram.WithLabelValues(name).Observe(v)
			}))
			defer timer.ObserveDuration()

			err := fn(ctx, req, rsp)
			if err == nil {
				opsCounter.WithLabelValues(name, "success").Inc()
			} else {
				opsCounter.WithLabelValues(name, "fail").Inc()
			}

			return err
		}
	}
}
