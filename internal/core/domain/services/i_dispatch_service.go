package services

import (
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/courier"
	"github.com/victor-tsykanov/delivery/internal/core/domain/model/order"
)

type IDispatchService interface {
	Dispatch(order *order.Order, couriers []*courier.Courier) (*courier.Courier, error)
}
