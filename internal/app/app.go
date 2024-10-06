package app

import (
	"context"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"golang.org/x/sync/errgroup"

	"github.com/Karzoug/meower-post-service/internal/config"
	"github.com/Karzoug/meower-post-service/pkg/buildinfo"
	"github.com/Karzoug/meower-post-service/pkg/metric/prom"
	"github.com/Karzoug/meower-post-service/pkg/trace/otlp"
)

const (
	serviceName     = "PostService"
	metricNamespace = "post_service"
	pkgName         = "github.com/Karzoug/meower-post-service"
	initTimeout     = 10 * time.Second
)

var serviceVersion = buildinfo.Get().Version

func Run(ctx context.Context, logger zerolog.Logger) error {
	cfg, err := env.ParseAs[config.Config]()
	if err != nil {
		return err
	}
	zerolog.SetGlobalLevel(cfg.LogLevel)

	// set timeout for initialization
	ctxInit, closeCtx := context.WithTimeout(ctx, initTimeout)
	defer closeCtx()

	// set up tracer
	shutdownTracer, err := otlp.RegisterGlobal(ctxInit, serviceName, serviceVersion)
	if err != nil {
		return err
	}
	defer close(shutdownTracer, logger)

	_ = otel.GetTracerProvider().Tracer(pkgName)

	// set up meter
	shutdownMeter, err := prom.RegisterGlobal(ctxInit, serviceName, serviceVersion, metricNamespace)
	if err != nil {
		return err
	}
	defer close(shutdownMeter, logger)

	eg, ctx := errgroup.WithContext(ctx)

	// run prometheus metrics http server
	eg.Go(func() error {
		return prom.Serve(ctx, cfg.PromHTTP, logger)
	})

	return eg.Wait()
}

func close(fn func(context.Context) error, logger zerolog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := fn(ctx); err != nil {
		logger.Error().
			Err(err).
			Msg("error closing")
	}
}
