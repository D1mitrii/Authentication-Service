package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	namespace = "auth"
)

type Metrics struct {
	grpcRequestTotal *prometheus.CounterVec
	grpcDuration     *prometheus.HistogramVec
	httpRequestTotal *prometheus.CounterVec
	httpDuration     *prometheus.HistogramVec
}

var metrics *Metrics

func Init(reg prometheus.Registerer) error {

	metrics = &Metrics{
		grpcRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "grpc_request_total",
				Help:      "Toatl gRPC request counter",
			},
			[]string{"status", "method"},
		),
		grpcDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "grpc_request_duration_seconds",
				Help:      "Duration of gRPC request",
				Buckets: []float64{
					0.1,
					0.2,
					0.25,
					0.5,
					1,
				},
			},
			[]string{"status", "method"},
		),
		httpRequestTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      "http_request_total",
				Help:      "Total http request counter",
			},
			[]string{"status", "method", "path"},
		),
		httpDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      "http_request_duration_seconds",
				Help:      "Duration of HTTP request",
				Buckets: []float64{
					0.1,
					0.2,
					0.25,
					0.5,
					1,
				},
			},
			[]string{"status", "method", "path"},
		),
	}

	reg.MustRegister(
		metrics.grpcRequestTotal,
		metrics.grpcDuration,
		metrics.httpRequestTotal,
		metrics.httpDuration,
	)

	return nil
}

func GrpcCounterRequestTotal(status string, method string) {
	metrics.grpcRequestTotal.WithLabelValues(status, method).Inc()
}

func GrpcHistogramResponseTimeObserve(status string, method string, time float64) {
	metrics.grpcDuration.WithLabelValues(status, method).Observe(time)
}

func HttpCounterRequestTotal(status string, method string, path string) {
	metrics.httpRequestTotal.WithLabelValues(status, method, path).Inc()
}

func HttpHistogramResponseTimeObserve(status string, method string, path string, time float64) {
	metrics.httpDuration.WithLabelValues(status, method, path).Observe(time)
}
