package middleware

import (
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
)

func Metrics() func(next http.Handler) http.Handler {
	totalRequests := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name:      "requests_total",
			Subsystem: "http_server",
			Help:      "The number of the HTTP requests.",
		},
		[]string{"pattern", "code"},
	)

	mustRegisterMetrics(totalRequests)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cw := newWrapResponseWriter(w)

			defer func() {
				statusCode := strconv.Itoa(cw.Status())
				totalRequests.WithLabelValues(r.Pattern, statusCode).Inc()
			}()

			next.ServeHTTP(cw, r)
		})
	}
}

func mustRegisterMetrics(collectors ...prometheus.Collector) {
	for i := range collectors {
		if err := prometheus.Register(collectors[i]); err != nil {
			//nolint:errorlint
			if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
				panic(err)
			}
		}
	}
}
