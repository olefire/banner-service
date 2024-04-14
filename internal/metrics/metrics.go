package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	HTTPRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "HTTP_request_total",
		Help: "The total number of processed requests",
	}, []string{"method", "result"})
	HTTPRequestTotalFail    = HTTPRequestTotal.MustCurryWith(prometheus.Labels{"result": "fail"})
	HTTPRequestTotalSuccess = HTTPRequestTotal.MustCurryWith(prometheus.Labels{"result": "success"})

	HTTPErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "HTTP_errors_total",
		Help: "The total number of HTTP errors",
	}, []string{"method"})

	HTTPLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "HTTP_duration_seconds",
		Help:    "Latency of HTTP requests in seconds",
		Buckets: prometheus.ExponentialBuckets(0.001, 2, 17),
	}, []string{"method"})
)
