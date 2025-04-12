package outbox

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/outbox/eventserializers"
	"github.com/victor-tsykanov/delivery/internal/common/ddd"
	"github.com/victor-tsykanov/delivery/internal/common/eventdispatcher"
)

type EventsOutbox struct {
	repository               IRepository
	eventSerializersRegistry eventserializers.IRegistry
	eventDispatcher          eventdispatcher.IEventDispatcher
	processNewMessagesNowCh  chan struct{}
	cancelCh                 chan struct{}
}

func NewEventsOutbox(
	repository IRepository,
	eventSerializersRegistry eventserializers.IRegistry,
	eventDispatcher eventdispatcher.IEventDispatcher,
) (*EventsOutbox, error) {
	return &EventsOutbox{
		repository:               repository,
		eventSerializersRegistry: eventSerializersRegistry,
		eventDispatcher:          eventDispatcher,
		processNewMessagesNowCh:  make(chan struct{}),
		cancelCh:                 make(chan struct{}),
	}, nil
}

func (o *EventsOutbox) StoreAggregateEvents(ctx context.Context, aggregate ddd.IAggregateRoot) error {
	for _, event := range aggregate.DomainEvents() {
		err := o.storeEvent(ctx, event)
		if err != nil {
			return err
		}
	}

	aggregate.ClearDomainEvents()

	return nil
}

func (o *EventsOutbox) ProcessNewMessages(ctx context.Context) {
	const pollInterval = 2 * time.Second

	for {
		o.processNewMessages(ctx)

		select {
		case <-time.After(pollInterval):
			continue
		case <-o.processNewMessagesNowCh:
			continue
		case <-o.cancelCh:
			log.Println("shutting down")
			return
		}
	}
}

func (o *EventsOutbox) ShutDown() {
	o.cancelCh <- struct{}{}
}

func (o *EventsOutbox) processNewMessages(ctx context.Context) {
	const batchSize = 20

	messages, hasMore, err := o.repository.FindUnprocessed(batchSize)
	if err != nil {
		log.Printf("failed to get unprocessd messages: %v\n", err)
		return
	}

	for _, event := range messages {
		o.processMessage(ctx, event)
	}

	if hasMore {
		o.processNewMessagesNow()
	}
}

func (o *EventsOutbox) processMessage(ctx context.Context, message *Message) {
	log.Printf("processing message %s, id=%s", message.Type, message.ID)

	event, err := o.eventSerializersRegistry.Deserialize(message.Type, message.Payload)
	if err != nil {
		log.Printf("failed to deserialize event from message %s, id=%s: %v\n", message.Type, message.ID, err)
		return
	}

	err = o.eventDispatcher.Dispatch(ctx, event)
	if err != nil {
		log.Printf("failed to process event %s, id=%s: %v\n", event.Type(), event.ID(), err)
		return
	}

	err = o.repository.MarkProcessed(ctx, message.ID)
	if err != nil {
		log.Printf("failed to mark message %s, id=%s as processed: %v\n", message.ID, message.Type, err)
		return
	}
}

func (o *EventsOutbox) processNewMessagesNow() {
	go func() {
		o.processNewMessagesNowCh <- struct{}{}
	}()
}

func (o *EventsOutbox) storeEvent(ctx context.Context, event ddd.IDomainEvent) error {
	payload, err := o.eventSerializersRegistry.Serialize(event)
	if err != nil {
		return err
	}

	message := &Message{
		ID:      event.ID(),
		Type:    event.Type(),
		Payload: payload,
	}

	err = o.repository.Create(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to save message: %w", err)
	}

	return nil
}
