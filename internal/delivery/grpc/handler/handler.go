package handler

import (
	"github.com/rs/zerolog"

	pb "github.com/Karzoug/meower-post-service/internal/delivery/grpc/gen"
)

var _ pb.PostServer = PostHandlers{}

type PostHandlers struct {
	logger zerolog.Logger
	pb.UnimplementedPostServer
}

func New(logger zerolog.Logger) PostHandlers {
	logger = logger.With().Str("component", "grpc server handler").Logger()

	return PostHandlers{
		logger: logger,
	}
}
