package services

import (
	"fmt"
	"slices"

	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

type DispatchService struct{}

func NewDispatchService() (*DispatchService, error) {
	return &DispatchService{}, nil
}

func (s *DispatchService) Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error) {
	if len(couriers) == 0 {
		return nil, errors.NewValueIsRequiredError("couriers")
	}

	orderLocation := order.Location()

	stepsByCourierIdx := make([]int, len(couriers))
	for i, c := range couriers {
		steps, err := c.CalculateStepsToLocation(&orderLocation)
		if err != nil {
			return nil, fmt.Errorf("failed to calculate steps to order location for courier %s: %w", c.ID(), err)
		}

		stepsByCourierIdx[i] = steps
	}

	closetCourierIdx := slices.Index(
		stepsByCourierIdx,
		slices.Min(stepsByCourierIdx),
	)

	if closetCourierIdx == -1 {
		return nil, fmt.Errorf("failed to find best courier")
	}

	return couriers[closetCourierIdx], nil
}
