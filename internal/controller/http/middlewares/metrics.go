package middlewares

import (
	"github.com/d1mitrii/authentication-service/internal/metrics"
	"net/http"
	"time"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		srw := NewStatusResponseWriter(w)

		next.ServeHTTP(srw, r)

		path := r.RequestURI
		method := r.Method
		status := srw.StatusString()

		metrics.HttpCounterRequestTotal(status, method, path)
		metrics.HttpHistogramResponseTimeObserve(status, method, path, time.Since(start).Seconds())
	})
}
