package interceptors

import (
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func Logger(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	l := logger.With().Logger()

	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		//nolint:zerologlint
		event := l.Info().
			Str("method", info.FullMethod)

		h, err := handler(ctx, req)

		span := trace.SpanFromContext(ctx)
		traceID := span.SpanContext().TraceID()
		if traceID.IsValid() {
			event.Str("trace_id", traceID.String())
		}

		defer func() {
			event.
				Err(err).
				Msg("incoming_request")
		}()

		return h, err
	}
}
