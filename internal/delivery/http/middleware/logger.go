package middleware

import (
	"net/http"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
)

// Logger is a middleware that logs incoming requests.
func Logger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	l := logger.With().Logger()

	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			cw := newWrapResponseWriter(w)

			//nolint:zerologlint
			event := l.Info().
				Str("remote_addr", r.RemoteAddr).
				Str("path", r.URL.Path).
				Str("proto", r.Proto).
				Str("method", r.Method).
				Str("user_agent", r.UserAgent()).
				Int64("bytes_in", r.ContentLength)

			span := trace.SpanFromContext(r.Context())
			traceID := span.SpanContext().TraceID()
			if traceID.IsValid() {
				w.Header().Add("X-Trace-Id", traceID.String())
				event.Str("trace_id", traceID.String())
			}

			defer func() {
				event.
					Int("status_code", cw.Status()).
					Msg("incoming_request")
			}()
			next.ServeHTTP(cw, r)
		}
		return http.HandlerFunc(fn)
	}
}
