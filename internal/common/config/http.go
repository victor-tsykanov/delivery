package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type HTTPConfig struct {
	Host string `env:"HTTP_HOST,notEmpty" envDefault:"localhost"`
	Port int    `env:"HTTP_PORT,notEmpty"`
}

func LoadHTTPConfig() (*HTTPConfig, error) {
	config := &HTTPConfig{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}

func MustLoadHTTPConfig() *HTTPConfig {
	config, err := LoadHTTPConfig()
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func (c *HTTPConfig) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
