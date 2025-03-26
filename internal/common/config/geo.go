package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type GeoServiceConfig struct {
	Address string `env:"GEO_SERVICE_GRPC_ADDRESS,notEmpty"`
}

func LoadGeoServiceConfig() (*GeoServiceConfig, error) {
	config := &GeoServiceConfig{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}

func MustLoadGeoServiceConfig() *GeoServiceConfig {
	config, err := LoadGeoServiceConfig()
	if err != nil {
		log.Fatal(err)
	}

	return config
}
