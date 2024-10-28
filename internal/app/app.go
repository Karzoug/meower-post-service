package app

import (
	"context"
	"runtime"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"

	"github.com/Karzoug/meower-post-service/internal/config"
)

func Run(ctx context.Context, logger zerolog.Logger) error {
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(cfg.LogLevel)

	logger.Info().
		Int("GOMAXPROCS", runtime.GOMAXPROCS(0)).
		Str("log level", cfg.LogLevel.String()).
		Msg("starting up")

	return nil
}
