package config

import (
	"fmt"
	"github.com/altiby/son/internal/goconf"
	"github.com/altiby/son/internal/logger"
	"github.com/altiby/son/internal/metrics"
	"github.com/altiby/son/internal/storage"
)

// Config is the main config for the application
type Config struct {
	Server     ServerConfig
	Postgresql storage.Config
}

// Server configuration.
type ServerConfig struct {
	Port    int            `mapstructure:"port" valid:"required"`
	Logging logger.Config  `mapstructure:"logging" valid:"-"`
	Metrics metrics.Config `mapstructure:"metrics" valid:"-"`
}

// Get loads config from requested path for specified env.
func Get(path, envKey string) (*Config, error) {
	var cfg Config

	if err := goconf.Load(path, envKey, &cfg); err != nil {
		return nil, fmt.Errorf("could not load config: %w", err)
	}

	return &cfg, nil
}
