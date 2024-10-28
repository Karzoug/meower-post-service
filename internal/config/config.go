package config

import (
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel zerolog.Level `env:"LOG_LEVEL" envDefault:"info"`
}
