package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Recoverer is a middleware that prevents panics and logs them.
func Recoverer(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					if rec == http.ErrAbortHandler { //nolint:errorlint
						// we don't recover http.ErrAbortHandler so the response
						// to the client is aborted, this should not be logged
						logger.Panic().
							Interface("recover", rec).
							Msg("got http.ErrAbortHandler")
					}

					event := logger.WithLevel(zerolog.PanicLevel).
						Interface("recover", rec).
						Bytes("stack", debug.Stack())

					span := trace.SpanFromContext(r.Context())
					traceID := span.SpanContext().TraceID()
					if traceID.IsValid() {
						event.Str("trace_id", traceID.String())
					}

					event.Msg("recovered panic")

					if r.Header.Get("Connection") != "Upgrade" {
						w.WriteHeader(http.StatusInternalServerError)
					}
				}
			}()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
