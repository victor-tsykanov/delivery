package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/victor-tsykanov/delivery/pkg/servers"
)

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
