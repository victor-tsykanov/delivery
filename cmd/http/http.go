package http

import (
	"context"
	"log"
	"net/http"

	"github.com/flowchartsman/swaggerui"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/victor-tsykanov/delivery/cmd/app"
	httpAdapters "github.com/victor-tsykanov/delivery/internal/adapters/in/http"
	"github.com/victor-tsykanov/delivery/internal/common/config"
	"github.com/victor-tsykanov/delivery/pkg/servers"
)

func Serve(_ context.Context, root *app.CompositionRoot, httpConfig *config.HTTPConfig) {
	server, err := httpAdapters.NewServer(
		root.CommandHandlers.CreateOrderCommandHandler,
		root.QueryHandlers.GetAllCouriersQueryHandler,
		root.QueryHandlers.GetPendingOrdersQueryHandler,
	)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())

	setUpCORS(e)
	setUpSwagger(e)

	handler := servers.NewStrictHandler(server, nil)
	servers.RegisterHandlers(e, handler)

	err = e.Start(httpConfig.Address())
	if err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}

func setUpCORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{
				echo.GET,
				echo.POST,
				echo.PUT,
				echo.DELETE,
				echo.OPTIONS,
			},
		}),
	)
}

func setUpSwagger(e *echo.Echo) {
	swagger, err := servers.GetSwagger()
	if err != nil {
		log.Fatalf("failed to load OpenAPI spec: %v", err)
	}

	spec, err := swagger.MarshalJSON()
	if err != nil {
		log.Fatalf("failed marshal OpenAPI spec to JSON: %v", err)
	}

	e.GET(
		"/swagger/*",
		echo.WrapHandler(
			http.StripPrefix(
				"/swagger",
				swaggerui.Handler(spec),
			),
		),
	)
}
