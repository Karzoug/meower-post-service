package config

import (
	"github.com/rs/zerolog"

	"github.com/Karzoug/meower-common-go/metric/prom"
	"github.com/Karzoug/meower-common-go/mongo"
	"github.com/Karzoug/meower-common-go/trace/otlp"

	grpc "github.com/Karzoug/meower-post-service/internal/delivery/grpc/server"
)

type Config struct {
	LogLevel zerolog.Level     `env:"LOG_LEVEL" envDefault:"info"`
	GRPC     grpc.Config       `envPrefix:"GRPC_"`
	PromHTTP prom.ServerConfig `envPrefix:"PROM_"`
	OTLP     otlp.Config       `envPrefix:"OTLP_"`
	Mongo    mongo.Config      `envPrefix:"MONGO_"`
}
