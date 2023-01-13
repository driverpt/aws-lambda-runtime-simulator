package controller

import "github.com/labstack/echo/v4"

type Controller interface {
	RegisterRoutes(e *echo.Echo) error
}
