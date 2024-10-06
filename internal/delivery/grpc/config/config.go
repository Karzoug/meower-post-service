package config

import (
	"fmt"
	"time"
)

type Config struct {
	Host           string               `env:"HOST"`
	Port           uint                 `env:"PORT,notEmpty"`
	CircuitBreaker CircuitBreakerConfig `envPrefix:"CB_"`
}

func (cfg Config) Address() string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

type CircuitBreakerConfig struct {
	// Interval is a cyclic period of the closed state for CircuitBreaker to clear the internal Counts
	Interval time.Duration `env:"INTERVAL" envDefault:"10s"`
	// Timeout is a period of the open state, after which the state of CircuitBreaker becomes half-open
	Timeout time.Duration `env:"TIMEOUT" envDefault:"30s"`
	// MaxRequests is a maximum number of requests allowed to pass through when the CircuitBreaker is half-open
	MaxRequests uint32 `env:"MAX_REQUESTS" envDefault:"10"`
}
