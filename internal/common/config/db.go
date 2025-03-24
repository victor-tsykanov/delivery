package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type DBConfig struct {
	Host     string `env:"DB_HOST,notEmpty"`
	Port     int    `env:"DB_PORT,notEmpty"      envDefault:"5432"`
	User     string `env:"DB_USER,notEmpty"`
	Password string `env:"DB_PASSWORD,notEmpty"`
	Database string `env:"DB_DATABASE,notEmpty"`
	SSLMode  string `env:"DB_SSL_MODE,notEmpty"  envDefault:"disable"`
}

func LoadDBConfig() (*DBConfig, error) {
	config := &DBConfig{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}

func MustLoadDBConfig() *DBConfig {
	config, err := LoadDBConfig()
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func (c *DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		c.Host,
		c.User,
		c.Password,
		c.Database,
		c.Port,
		c.SSLMode,
	)
}
