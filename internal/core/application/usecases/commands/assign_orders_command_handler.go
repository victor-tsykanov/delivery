package commands

import (
	"context"
	"fmt"

	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
	"github.com/victor-tsykanov/delivery/internal/core/domain/services"
	outPorts "github.com/victor-tsykanov/delivery/internal/core/ports/out"
)

type AssignOrdersCommandHandler struct {
	transactionManager outPorts.ITransactionManager
	dispatchService    services.IDispatchService
	courierRepository  outPorts.ICourierRepository
	orderRepository    outPorts.IOrderRepository
}

func NewAssignOrdersCommandHandler(
	transactionManager outPorts.ITransactionManager,
	dispatchService services.IDispatchService,
	courierRepository outPorts.ICourierRepository,
	orderRepository outPorts.IOrderRepository,
) (*AssignOrdersCommandHandler, error) {
	if transactionManager == nil {
		return nil, errors.NewValueIsRequiredError("transactionManager")
	}

	if dispatchService == nil {
		return nil, errors.NewValueIsRequiredError("dispatchService")
	}

	if courierRepository == nil {
		return nil, errors.NewValueIsRequiredError("courierRepository")
	}

	if orderRepository == nil {
		return nil, errors.NewValueIsRequiredError("orderRepository")
	}

	return &AssignOrdersCommandHandler{
		transactionManager: transactionManager,
		dispatchService:    dispatchService,
		courierRepository:  courierRepository,
		orderRepository:    orderRepository,
	}, nil
}

func (h *AssignOrdersCommandHandler) Handle(ctx context.Context) error {
	newOrders, err := h.orderRepository.FindNew(ctx)
	if err != nil {
		return err
	}

	for _, order := range newOrders {
		err = h.transactionManager.Execute(
			ctx,
			func(ctx context.Context) error {
				return h.assignOrder(ctx, order)
			},
		)

		if err != nil {
			return fmt.Errorf("failed to assign order %s: %w", order.ID(), err)
		}
	}

	return nil
}

func (h *AssignOrdersCommandHandler) assignOrder(ctx context.Context, order *order.Order) error {
	freeCouriers, err := h.courierRepository.FindFree(ctx)
	if err != nil {
		return err
	}

	if len(freeCouriers) == 0 {
		return errors.NoAvailableCouriersError
	}

	courier, err := h.dispatchService.Dispatch(order, freeCouriers)
	if err != nil {
		return err
	}

	err = order.Assign(courier)
	if err != nil {
		return err
	}

	err = courier.SetBusy()
	if err != nil {
		return err
	}

	err = h.orderRepository.Update(ctx, order)
	if err != nil {
		return err
	}

	err = h.courierRepository.Update(ctx, courier)
	if err != nil {
		return err
	}

	return nil
}
