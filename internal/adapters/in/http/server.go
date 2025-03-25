package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	"github.com/victor-tsykanov/delivery/pkg/servers"
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

func (s *Server) GetCouriers(
	ctx context.Context,
	_ servers.GetCouriersRequestObject,
) (servers.GetCouriersResponseObject, error) {
	couriers, err := s.getAllCouriersQueryHandler.Handle(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	courierDTOs := make([]servers.Courier, len(couriers))
	for i, courier := range couriers {
		courierDTOs[i] = servers.Courier{
			Id: courier.ID,
			Location: servers.Location{
				X: courier.Location.X,
				Y: courier.Location.Y,
			},
			Name: courier.Name,
		}
	}

	return servers.GetCouriers200JSONResponse(courierDTOs), nil
}

func (s *Server) CreateOrder(
	ctx context.Context,
	_ servers.CreateOrderRequestObject,
) (servers.CreateOrderResponseObject, error) {
	command, err := inPorts.NewCreateOrderCommand(uuid.New(), "...")
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	err = s.createOrderCommandHandler.Handle(ctx, *command)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusUnprocessableEntity, err)
	}

	return servers.CreateOrder201Response{}, nil
}

func (s *Server) GetOrders(
	ctx context.Context,
	_ servers.GetOrdersRequestObject,
) (servers.GetOrdersResponseObject, error) {
	orders, err := s.getPendingOrdersQueryHandler.Handle(ctx)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	orderDTOs := make([]servers.Order, len(orders))
	for i, order := range orders {
		orderDTOs[i] = servers.Order{
			Id: order.ID,
			Location: servers.Location{
				X: order.Location.X,
				Y: order.Location.Y,
			},
		}
	}

	return servers.GetOrders200JSONResponse(orderDTOs), nil
}
