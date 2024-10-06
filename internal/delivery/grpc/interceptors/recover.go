package interceptors

import (
	"context"
	"runtime/debug"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// Recover returns a new unary server interceptor for panic recovery.
func Recover(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ any, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				event := logger.WithLevel(zerolog.PanicLevel).
					Interface("recover", rec).
					Bytes("stack", debug.Stack())

				span := trace.SpanFromContext(ctx)
				traceID := span.SpanContext().TraceID()
				if traceID.IsValid() {
					event.Str("trace_id", traceID.String())
				}

				event.Msg("recovered panic")
			}
		}()

		return handler(ctx, req)
	}
}
