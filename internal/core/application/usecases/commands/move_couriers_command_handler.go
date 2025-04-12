package commands

import (
	"context"
	"errors"
	"fmt"

	commonErrors "github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/common/persistence"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	outPorts "github.com/victor-tsykanov/delivery/internal/core/ports/out"
)

type MoveCouriersCommandHandler struct {
	transactionManager persistence.ITransactionManager
	courierRepository  outPorts.ICourierRepository
	orderRepository    outPorts.IOrderRepository
}

func NewMoveCouriersCommandHandler(
	transactionManager persistence.ITransactionManager,
	courierRepository outPorts.ICourierRepository,
	orderRepository outPorts.IOrderRepository,
) (*MoveCouriersCommandHandler, error) {
	if transactionManager == nil {
		return nil, commonErrors.NewValueIsRequiredError("transactionManager")
	}

	if courierRepository == nil {
		return nil, commonErrors.NewValueIsRequiredError("courierRepository")
	}

	if orderRepository == nil {
		return nil, commonErrors.NewValueIsRequiredError("orderRepository")
	}

	return &MoveCouriersCommandHandler{
		transactionManager: transactionManager,
		courierRepository:  courierRepository,
		orderRepository:    orderRepository,
	}, nil
}

func (h *MoveCouriersCommandHandler) Handle(ctx context.Context) error {
	orders, err := h.orderRepository.FindAssigned(ctx)
	if err != nil {
		return fmt.Errorf("failed to find assigned orders: %w", err)
	}

	var errs []error
	for _, order := range orders {
		err = h.transactionManager.Execute(ctx, func(ctx context.Context) error {
			return h.processOrder(ctx, order)
		})

		if err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}

func (h *MoveCouriersCommandHandler) processOrder(ctx context.Context, order *order.Order) error {
	courier, err := h.courierRepository.Get(ctx, *order.CourierID())
	if err != nil {
		return fmt.Errorf("failed to find courier to which order %s is assigned: %w", order.ID(), err)
	}

	targetLocation := order.Location()
	err = courier.Move(&targetLocation)
	if err != nil {
		return fmt.Errorf("failed to move courier %s: %w", courier.ID(), err)
	}

	if courier.Location().Equals(targetLocation) {
		err = courier.SetFree()
		if err != nil {
			return fmt.Errorf("failed to set courier %s free: %w", courier.ID(), err)
		}

		err = order.Complete()
		if err != nil {
			return fmt.Errorf("failed to complete order %s: %w", order.ID(), err)
		}

		err = h.orderRepository.Update(ctx, order)
		if err != nil {
			return fmt.Errorf("failed to update order %s: %w", order.ID(), err)
		}
	}

	err = h.courierRepository.Update(ctx, courier)
	if err != nil {
		return fmt.Errorf("failed to update courier %s: %w", courier.ID(), err)
	}

	return nil
}
