package http

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	inPorts "github.com/victor-tsykanov/delivery/internal/core/ports/in"
	"github.com/victor-tsykanov/delivery/pkg/servers"
)

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
