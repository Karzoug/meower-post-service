package interceptors

import (
	"context"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	grpcCodes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Karzoug/meower-post-service/internal/delivery/http/config"
)

func CircuitBreaker(cfg config.CircuitBreakerConfig) grpc.UnaryServerInterceptor {
	cb := gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
		Name:        "circuit-breaker",
		Interval:    cfg.Interval,
		Timeout:     cfg.Timeout,
		MaxRequests: cfg.MaxRequests,
	})
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		done, err := cb.Allow()
		if err != nil {
			return nil, err
		}
		h, err := handler(ctx, req)

		s, ok := status.FromError(err)
		done(s.Code() != grpcCodes.Internal && ok)

		return h, err
	}
}
