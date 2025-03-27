package http

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/victor-tsykanov/delivery/pkg/servers"
)

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
