package config

import (
	grpcConfig "github.com/Karzoug/meower-post-service/internal/delivery/grpc/config"
	httpConfig "github.com/Karzoug/meower-post-service/internal/delivery/http/config"
	"github.com/Karzoug/meower-post-service/pkg/metric/prom"

	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel zerolog.Level           `env:"LOG_LEVEL"        envDefault:"info"`
	HTTP     httpConfig.ServerConfig `envPrefix:"HTTP_"`
	GRPC     grpcConfig.Config       `envPrefix:"GRPC_"`
	PromHTTP prom.ServerConfig       `envPrefix:"PROM_HTTP_"`
}
