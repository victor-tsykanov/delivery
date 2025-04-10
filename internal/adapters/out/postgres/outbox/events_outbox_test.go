package outbox_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/outbox"
	"github.com/victor-tsykanov/delivery/internal/common/ddd"
	outboxMocks "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/adapters/out/postgres/outbox"
	eventserializersMocks "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/adapters/out/postgres/outbox/eventserializers"
	dddMocks "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/common/ddd"
	eventdispatcherMocks "github.com/victor-tsykanov/delivery/mocks/github.com/victor-tsykanov/delivery/internal_/common/eventdispatcher"
)

type testEvent struct {
	id uuid.UUID
}

func newTestEvent() *testEvent {
	return &testEvent{id: uuid.New()}
}

func (e *testEvent) ID() uuid.UUID {
	return e.id
}

func (e *testEvent) Type() string {
	return "test"
}

func TestEventsOutbox_StoreAggregateEvents(t *testing.T) {
	// Arrange
	ctx := context.Background()
	repository := outboxMocks.NewMockIRepository(t)
	eventSerializersRegistry := eventserializersMocks.NewMockIRegistry(t)
	eventDispatcher := eventdispatcherMocks.NewMockIEventDispatcher(t)
	aggregateRoot := dddMocks.NewMockIAggregateRoot(t)
	events := []ddd.IDomainEvent{newTestEvent(), newTestEvent()}

	eventsOutbox, err := outbox.NewEventsOutbox(repository, eventSerializersRegistry, eventDispatcher)
	require.NoError(t, err)

	aggregateRoot.
		EXPECT().
		DomainEvents().
		Return(events).
		Once()

	aggregateRoot.
		EXPECT().
		ClearDomainEvents().
		Return().
		Once()

	for _, event := range events {
		eventSerializersRegistry.
			EXPECT().
			Serialize(event).
			Return([]byte("..."), nil).
			Once()

		repository.
			EXPECT().
			Create(ctx, &outbox.Message{
				ID:          event.ID(),
				Type:        event.Type(),
				Payload:     []byte("..."),
				ProcessedAt: sql.NullTime{},
			}).
			Return(nil).
			Once()
	}

	// Act
	err = eventsOutbox.StoreAggregateEvents(ctx, aggregateRoot)

	// Assert
	require.NoError(t, err)
}
