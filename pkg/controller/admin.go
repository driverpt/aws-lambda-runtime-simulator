package controller

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io"
	"lambda-runtime-simulator/pkg/event"
	"net/http"
)

type AdminController struct {
	service *event.Service
}

var _ Controller = &AdminController{}

func NewAdminController(svc *event.Service) *AdminController {
	return &AdminController{
		service: svc,
	}
}

func (a AdminController) RegisterRoutes(e *echo.Echo) error {
	e.GET("/log", nil)
	e.GET("/log/stream", nil)
	e.GET("/response/:requestId", nil)
	e.POST("/log/clear", nil)
	e.POST("/push", a.Push)

	return nil
}

func (a AdminController) Push(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Json sanity check
	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = a.service.PushInvocation(string(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
}
