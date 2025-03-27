package http

import (
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
)

type Server struct {
	createOrderCommandHandler    inPorts.ICreateOrderCommandHandler
	getAllCouriersQueryHandler   inPorts.IGetAllCouriersQueryHandler
	getPendingOrdersQueryHandler inPorts.IGetPendingOrdersQueryHandler
}

func NewServer(
	createOrderCommandHandler inPorts.ICreateOrderCommandHandler,
	getAllCouriersQueryHandler inPorts.IGetAllCouriersQueryHandler,
	getPendingOrdersQueryHandler inPorts.IGetPendingOrdersQueryHandler,
) (*Server, error) {
	return &Server{
		createOrderCommandHandler:    createOrderCommandHandler,
		getAllCouriersQueryHandler:   getAllCouriersQueryHandler,
		getPendingOrdersQueryHandler: getPendingOrdersQueryHandler,
	}, nil
}
