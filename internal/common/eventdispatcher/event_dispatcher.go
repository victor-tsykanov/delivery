package eventdispatcher

import (
	"context"
	"fmt"
	"reflect"
)

type EventDispatcher struct {
	handlers map[reflect.Type]internalHandler
}

func New() *EventDispatcher {
	return &EventDispatcher{
		handlers: make(map[reflect.Type]internalHandler),
	}
}

func Register[T any](dispatcher *EventDispatcher, handler IEventHandler[T]) {
	eventType := reflect.TypeOf(*new(T))
	dispatcher.handlers[eventType] = &handlerAdapter[T]{handler: handler}
}

func (d *EventDispatcher) Dispatch(ctx context.Context, event interface{}) error {
	if event == nil {
		return fmt.Errorf("event cannot be nil")
	}

	eventType := reflect.TypeOf(event)
	handler, exists := d.handlers[eventType]

	if !exists {
		return NewUnknownEventError(eventType.String())
	}

	return handler.handle(ctx, event)
}

// internalHandler is the type-erased interface for actual dispatch
type internalHandler interface {
	handle(ctx context.Context, message interface{}) error
	messageType() reflect.Type
}

// handlerAdapter bridges generic IEventHandler and internalHandler
type handlerAdapter[T any] struct {
	handler IEventHandler[T]
}

func (a *handlerAdapter[T]) handle(ctx context.Context, event interface{}) error {
	typedMessage, ok := event.(T)
	if !ok {
		return fmt.Errorf("type mismatch: expected %T, got %T", *new(T), event)
	}

	return a.handler.Handle(ctx, typedMessage)
}

func (a *handlerAdapter[T]) messageType() reflect.Type {
	return reflect.TypeOf(*new(T))
}
