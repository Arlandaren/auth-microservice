package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
)

var (
	requestCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_requests_total",
			Help: "Total number of requests received",
		},
		[]string{"method", "endpoint"},
	)

	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "app_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	//grpcRequestCounter = prometheus.NewCounterVec(
	//	prometheus.CounterOpts{
	//		Name: "grpc_requests_total",
	//		Help: "Total number of gRPC requests received",
	//	},
	//	[]string{"method"},
	//)
	//
	//grpcRequestDuration = prometheus.NewHistogramVec(
	//	prometheus.HistogramOpts{
	//		Name:    "grpc_request_duration_seconds",
	//		Help:    "Duration of gRPC requests",
	//		Buckets: prometheus.DefBuckets,
	//	},
	//	[]string{"method"},
	//)		решили убрать до лучших времен
)

type Prometheus struct {
	metrics []prometheus.Collector
}

func NewPrometheus() *Prometheus {
	return &Prometheus{
		metrics: []prometheus.Collector{
			requestCounter,
			requestDuration,
			//grpcRequestCounter, решено убрать
			//grpcRequestDuration, решено убрать
		},
	}
}

func (p *Prometheus) RegisterMetrics() {
	prometheus.MustRegister(p.metrics...)
}

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timer := prometheus.NewTimer(requestDuration.WithLabelValues(r.Method, r.URL.Path))
		defer timer.ObserveDuration()

		requestCounter.WithLabelValues(r.Method, r.URL.Path).Inc()

		next.ServeHTTP(w, r)
	})
}
