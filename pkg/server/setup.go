package server

import (
	"lambda-runtime-simulator/pkg/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func SetupServer(e *echo.Echo, controllers ...controller.Controller) error {
	e.Pre(middleware.RemoveTrailingSlash())

	for _, c := range controllers {
		err := c.RegisterRoutes(e)
		if err != nil {
			return err
		}
	}

	return nil
}
