package eventdispatcher

import "context"

type IEventDispatcher interface {
	Dispatch(ctx context.Context, event interface{}) error
}

type IEventHandler[T any] interface {
	Handle(ctx context.Context, message T) error
}
