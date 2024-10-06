package middleware

import (
	"net/http"

	"github.com/sony/gobreaker"

	"github.com/Karzoug/meower-post-service/internal/delivery/http/config"
)

func CircuitBreaker(cfg config.CircuitBreakerConfig) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &circuitHandler{
			cb: gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
				Name:        "circuit-breaker",
				Interval:    cfg.Interval,
				Timeout:     cfg.Timeout,
				MaxRequests: cfg.MaxRequests,
			}),
			next: next,
		}
	}
}

type circuitHandler struct {
	cb   *gobreaker.TwoStepCircuitBreaker
	next http.Handler
}

func (ch *circuitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	done, err := ch.cb.Allow()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	cw := newWrapResponseWriter(w)
	ch.next.ServeHTTP(cw, r)

	done(cw.code < http.StatusInternalServerError)
}
