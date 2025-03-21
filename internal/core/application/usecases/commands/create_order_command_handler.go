package commands

import (
	"context"
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/kernel"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	outPorts "github.com/victor-tsykanov/delivery/internal/core/ports/out"
)

type CreateOrderCommandHandler struct {
	transactionManager outPorts.ITransactionManager
	orderRepository    outPorts.IOrderRepository
}

func NewCreateOrderCommandHandler(
	transactionManager outPorts.ITransactionManager,
	orderRepository outPorts.IOrderRepository,
) (*CreateOrderCommandHandler, error) {
	if transactionManager == nil {
		return nil, errors.NewValueIsRequiredError("transactionManager")
	}

	if orderRepository == nil {
		return nil, errors.NewValueIsRequiredError("orderRepository")
	}

	return &CreateOrderCommandHandler{
		transactionManager: transactionManager,
		orderRepository:    orderRepository,
	}, nil
}

func (h *CreateOrderCommandHandler) Handle(ctx context.Context, command inPorts.CreateOrderCommand) error {
	location := kernel.RandomLocation()
	order, err := order.NewOrder(command.ID(), *location)
	if err != nil {
		return fmt.Errorf("falied to create order %s: %w", command.ID(), err)
	}

	err = h.transactionManager.Execute(ctx, func(ctx context.Context) error {
		return h.orderRepository.Create(ctx, order)
	})
	if err != nil {
		return fmt.Errorf("falied to save order %s: %w", order.ID(), err)
	}

	return nil
}
