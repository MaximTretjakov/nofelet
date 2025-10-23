package config

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	WS WSConfig `env:",prefix=WS_"`
}

type WSConfig struct {
	Port              string        `env:"PORT,required"`
	ReadTimeout       time.Duration `env:"READ_TIMEOUT,default=30s"`
	WriteTimeout      time.Duration `env:"WRITE_TIMEOUT,default=30s"`
	ReadHeaderTimeout time.Duration `env:"READ_HEADER_TIMEOUT,default=500ms"`
	ShutdownTimeout   time.Duration `env:"SHUTDOWN_TIMEOUT,default=3s"`
}

func init() {
	_ = godotenv.Load()
}

// NewConfig load current configuration.
func newConfig() (*Config, error) {
	var config Config
	if err := envconfig.Process(context.Background(), &config); err != nil {
		return nil, err
	}
	return &config, nil
}
