package out

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/victor-tsykanov/delivery/internal/common/config"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	"github.com/victor-tsykanov/delivery/pkg/queues/orderstatuschangedpb"
)

type OrderCompletedEventProducer struct {
	kafkaProducer *kafka.Producer
	topic         string
}

func NewOrderCompletedEventProducer(kafkaConfig *config.KafkaConfig) (*OrderCompletedEventProducer, error) {
	kafkaProducer, err := kafka.NewProducer(
		&kafka.ConfigMap{
			"bootstrap.servers": kafkaConfig.Address,
			"acks":              "all",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	return &OrderCompletedEventProducer{
		kafkaProducer: kafkaProducer,
		topic:         kafkaConfig.OrderStatusChangedTopic,
	}, nil
}

func (p *OrderCompletedEventProducer) Publish(_ context.Context, event order.CompletedEvent) error {
	integrationEvent := p.createIntegrationEvent(&event)
	payload, err := json.Marshal(integrationEvent)
	if err != nil {
		return fmt.Errorf("failed to marshal intgration event to JSON: %w", err)
	}

	err = p.kafkaProducer.Produce(
		&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &(p.topic), Partition: kafka.PartitionAny},
			Value:          payload,
		},
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to publish %s event to Kafka: %w", event.Type(), err)
	}

	return nil
}

func (p *OrderCompletedEventProducer) createIntegrationEvent(
	domainEvent *order.CompletedEvent,
) *orderstatuschangedpb.OrderStatusChangedIntegrationEvent {
	order := domainEvent.Order()

	return &orderstatuschangedpb.OrderStatusChangedIntegrationEvent{
		OrderId:     order.ID().String(),
		OrderStatus: orderstatuschangedpb.OrderStatus_Completed,
	}
}

func (p *OrderCompletedEventProducer) Close() error {
	p.kafkaProducer.Close()

	return nil
}
