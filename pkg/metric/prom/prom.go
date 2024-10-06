package prom

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func RegisterGlobal(
	ctx context.Context,
	serviceName,
	serviceVersion,
	namespace string,
) (shutdown func(context.Context) error, err error) {
	r, err := resource.Merge(resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion),
		))
	if err != nil {
		return nil, err
	}

	exporter, err := prometheus.New(
		prometheus.WithNamespace(namespace),
	)
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(exporter),
		metric.WithResource(r),
	)

	otel.SetMeterProvider(provider)

	return func(ctx context.Context) error {
		return provider.Shutdown(ctx)
	}, nil
}

func Serve(ctx context.Context, cfg ServerConfig, logger zerolog.Logger) error {
	logger = logger.With().Str("component", "prom http server").Logger()

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:         cfg.Address(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		Handler:      mux,
	}

	logger.Info().
		Str("address", cfg.Address()).
		Msg("listening")

	go func() {
		<-ctx.Done()

		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := srv.Shutdown(closeCtx); err != nil {
			logger.Error().
				Err(err).
				Msg("shutdown error")
		}
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
