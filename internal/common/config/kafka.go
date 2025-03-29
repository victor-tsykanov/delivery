package config

import (
	"fmt"
	"log"

	"github.com/caarlos0/env/v11"
)

type KafkaConfig struct {
	Address              string `env:"KAFKA_ADDRESS,notEmpty"`
	ConsumerGroup        string `env:"KAFKA_CONSUMER_GROUP,notEmpty"`
	BasketConfirmedTopic string `env:"KAFKA_BASKET_CONFIRMED_TOPIC,notEmpty"`
}

func LoadKafkaConfig() (*KafkaConfig, error) {
	config := &KafkaConfig{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	return config, nil
}

func MustLoadKafkaConfig() *KafkaConfig {
	config, err := LoadKafkaConfig()
	if err != nil {
		log.Fatal(err)
	}

	return config
}
