package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/google/uuid"
	"github.com/victor-tsykanov/delivery/internal/common/config"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	"github.com/victor-tsykanov/delivery/pkg/queues/basketconfirmedpb"
)

type BasketConfirmedConsumer struct {
	commandHandler inPorts.ICreateOrderCommandHandler
	kafkaConsumer  *kafka.Consumer
	topic          string
}

func NewBasketConfirmedConsumer(
	commandHandler inPorts.ICreateOrderCommandHandler,
	kafkaConfig *config.KafkaConfig,
) (*BasketConfirmedConsumer, error) {
	kafkaConsumer, err := kafka.NewConsumer(
		&kafka.ConfigMap{
			"bootstrap.servers":  kafkaConfig.Address,
			"group.id":           kafkaConfig.ConsumerGroup,
			"enable.auto.commit": false,
			"auto.offset.reset":  "earliest",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka client: %w", err)
	}

	return &BasketConfirmedConsumer{
		commandHandler: commandHandler,
		kafkaConsumer:  kafkaConsumer,
		topic:          kafkaConfig.BasketConfirmedTopic,
	}, nil
}

func (c *BasketConfirmedConsumer) Consume(ctx context.Context) {
	err := c.kafkaConsumer.Subscribe(c.topic, nil)
	if err != nil {
		log.Fatalf("failed to subscribe to topic: %s", err)
	}

	for {
		message, err := c.kafkaConsumer.ReadMessage(-1)
		if err != nil {
			log.Printf("failed to read message from Kafka: %v", err)
		} else {
			err = c.processMessage(ctx, message)
			if err != nil {
				log.Printf("failed to process event: %v", err)
			}
		}

		_, err = c.kafkaConsumer.CommitMessage(message)
		if err != nil {
			log.Printf("failed to commit message: %v", err)
		}
	}
}

func (c *BasketConfirmedConsumer) Close() error {
	err := c.kafkaConsumer.Close()
	if err != nil {
		return fmt.Errorf("failed to close Kafka consumer: %w", err)
	}

	return nil
}

func (c *BasketConfirmedConsumer) processMessage(ctx context.Context, message *kafka.Message) error {
	log.Printf("processing message %s: %s\n", c.topic, message.Value)

	var event basketconfirmedpb.BasketConfirmedIntegrationEvent
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	basketID, err := uuid.Parse(event.GetBasketId())
	if err != nil {
		return fmt.Errorf("failed to parse basket ID: %w", err)
	}

	command, err := inPorts.NewCreateOrderCommand(basketID, event.GetAddress().GetStreet())
	if err != nil {
		return fmt.Errorf("failed to create command: %w", err)
	}

	err = c.commandHandler.Handle(ctx, *command)
	if err != nil {
		return fmt.Errorf("failed to execute command: %w", err)
	}

	return nil
}
