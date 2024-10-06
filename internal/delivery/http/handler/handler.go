package handler

import (
	"github.com/rs/zerolog"

	"github.com/Karzoug/meower-post-service/internal/delivery/http/gen"
)

var _ gen.StrictServerInterface = PostHandlers{}

type PostHandlers struct {
	logger zerolog.Logger
}

func New(logger zerolog.Logger) PostHandlers {
	logger = logger.With().Str("component", "http server handler").Logger()

	return PostHandlers{
		logger: logger,
	}
}
