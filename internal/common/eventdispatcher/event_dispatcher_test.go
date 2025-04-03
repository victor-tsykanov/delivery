package eventdispatcher_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/victor-tsykanov/delivery/internal/common/eventdispatcher"
)

type Foo struct {
}

type FooHandler struct {
	called bool
}

func (h *FooHandler) Handle(_ context.Context, _ *Foo) error {
	h.called = true

	return nil
}

func TestEventDispatcher_Dispatch(t *testing.T) {
	dispatcher := eventdispatcher.New()
	fooHandler := &FooHandler{}
	eventdispatcher.Register(dispatcher, fooHandler)

	err := dispatcher.Dispatch(context.Background(), &Foo{})

	assert.NoError(t, err)
	assert.True(t, fooHandler.called)
}

func TestEventDispatcher_Dispatch_UnknownEvent(t *testing.T) {
	dispatcher := eventdispatcher.New()

	err := dispatcher.Dispatch(context.Background(), &Foo{})

	assert.Equal(t, eventdispatcher.NewUnknownEventError("*eventdispatcher_test.Foo"), err)
}
