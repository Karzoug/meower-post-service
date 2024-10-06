package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/Karzoug/meower-post-service/internal/delivery/http/config"
	"github.com/Karzoug/meower-post-service/internal/delivery/http/gen"
	"github.com/Karzoug/meower-post-service/internal/delivery/http/middleware"
	"github.com/Karzoug/meower-post-service/pkg/healthcheck"
)

const baseURL = "/api/v1"

type server struct {
	httpServer http.Server
	cfg        config.ServerConfig
	logger     zerolog.Logger
}

func New(cfg config.ServerConfig, handlers gen.StrictServerInterface, logger zerolog.Logger) *server {
	logger = logger.With().Str("component", "http server").Logger()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthcheck", healthcheck.HealthCheck(logger))
	mux.HandleFunc("GET /buildinfo", Buildinfo(logger))

	h := gen.HandlerWithOptions(
		gen.NewStrictHandler(handlers, nil),
		gen.StdHTTPServerOptions{
			Middlewares: []gen.MiddlewareFunc{
				middleware.Logger(logger),
				middleware.Recoverer(logger),
				// this middleware adds the following metrics:
				//   http_server_duration_milliseconds
				//   http_server_request_size_bytes_total
				//   http_server_response_size_bytes_total
				otelhttp.NewMiddleware("http server",
					otelhttp.WithMessageEvents(otelhttp.ReadEvents, otelhttp.WriteEvents),
					otelhttp.WithFilter(func(r *http.Request) bool {
						switch r.URL.Path {
						case "/buildinfo", "/healthcheck":
							return false
						default:
							return true
						}
					},
					),
				),
				middleware.Metrics(),
				middleware.Auth(),
				middleware.CircuitBreaker(cfg.CircuitBreaker),
			},
			BaseURL:    baseURL,
			BaseRouter: mux,
		})
	return &server{
		httpServer: http.Server{
			Addr:         cfg.Address(),
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			Handler:      h,
		},
		cfg:    cfg,
		logger: logger,
	}
}

func (s *server) Run(ctx context.Context) error {
	s.logger.Info().
		Str("address", s.cfg.Address()).
		Msg("listening")

	go func() {
		<-ctx.Done()

		closeCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := s.httpServer.Shutdown(closeCtx); err != nil {
			s.logger.Error().
				Err(err).
				Msg("shutdown error")
		}
	}()

	if err := s.httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}
