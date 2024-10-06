package grpc

import (
	"context"
	"fmt"
	"net"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/Karzoug/meower-post-service/internal/delivery/grpc/config"
	pb "github.com/Karzoug/meower-post-service/internal/delivery/grpc/gen"
	"github.com/Karzoug/meower-post-service/internal/delivery/grpc/interceptors"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type server struct {
	cfg    config.Config
	logger zerolog.Logger

	grpcServer *grpc.Server
}

func New(cfg config.Config, handlers pb.PostServer, logger zerolog.Logger) (*server, error) {
	logger = logger.With().Str("component", "grpc server").Logger()

	grpcServer := grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		// this interceptor adds the following metrics:
		//   rpc_server_duration_milliseconds
		//   rpc_server_request_size_bytes
		//   rpc_server_requests_per_rpc
		//   rpc_server_response_size_bytes
		//   rpc_server_responses_per_rpc
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithMessageEvents(otelgrpc.ReceivedEvents, otelgrpc.SentEvents),
		)),
		grpc.ChainUnaryInterceptor(
			interceptors.Logger(logger),
			interceptors.Recover(logger),
		),
	)

	ss := &server{
		cfg:        cfg,
		logger:     logger,
		grpcServer: grpcServer,
	}

	pb.RegisterPostServer(ss.grpcServer, handlers)
	reflection.Register(grpcServer)

	return ss, nil
}

func (s *server) Run(ctx context.Context) error {
	s.logger.Info().
		Str("address", s.cfg.Address()).
		Msg("listening")

	idleConnsClosed := make(chan struct{})

	var lc net.ListenConfig
	listen, err := lc.Listen(ctx, "tcp", s.cfg.Address())
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.cfg.Address(), err)
	}

	go func() {
		<-ctx.Done()
		s.grpcServer.GracefulStop()
		close(idleConnsClosed)
	}()

	if err := s.grpcServer.Serve(listen); err != nil {
		return fmt.Errorf("failed to serve grpc server: %w", err)
	}

	<-idleConnsClosed

	return nil
}
