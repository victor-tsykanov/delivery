package commands

import (
	"context"
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/common/persistence"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	outPorts "github.com/victor-tsykanov/delivery/internal/core/ports/out"
)

type CreateOrderCommandHandler struct {
	transactionManager persistence.ITransactionManager
	orderRepository    outPorts.IOrderRepository
	geoClient          outPorts.IGeoClient
}

func NewCreateOrderCommandHandler(
	transactionManager persistence.ITransactionManager,
	orderRepository outPorts.IOrderRepository,
	geoClient outPorts.IGeoClient,
) (*CreateOrderCommandHandler, error) {
	if transactionManager == nil {
		return nil, errors.NewValueIsRequiredError("transactionManager")
	}

	if orderRepository == nil {
		return nil, errors.NewValueIsRequiredError("orderRepository")
	}

	if geoClient == nil {
		return nil, errors.NewValueIsRequiredError("geoClient")
	}

	return &CreateOrderCommandHandler{
		transactionManager: transactionManager,
		orderRepository:    orderRepository,
		geoClient:          geoClient,
	}, nil
}

func (h *CreateOrderCommandHandler) Handle(ctx context.Context, command inPorts.CreateOrderCommand) error {
	location, err := h.geoClient.GetLocation(ctx, command.Street())
	if err != nil {
		return fmt.Errorf("falied to get location for street %s: %w", command.Street(), err)
	}

	order, err := order.NewOrder(order.ID(command.BasketID()), *location)
	if err != nil {
		return fmt.Errorf("falied to create order %s: %w", command.BasketID(), err)
	}

	err = h.transactionManager.Execute(ctx, func(ctx context.Context) error {
		return h.orderRepository.Create(ctx, order)
	})
	if err != nil {
		return fmt.Errorf("falied to save order %s: %w", order.ID(), err)
	}

	return nil
}
