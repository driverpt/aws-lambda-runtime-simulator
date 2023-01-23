package controller

import (
	"errors"
	"io"
	"lambda-runtime-simulator/pkg/event"
	"net/http"

	"github.com/labstack/echo/v4"
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
	e.GET("/invocation", a.GetLoggedInvocations)
	e.GET("/invocation/:invocationId", a.GetInvocationById)
	e.POST("/invocation/clear", a.ResetAllCaches)
	e.POST("/invocation", a.Push)
	return nil
}

func (a AdminController) GetLoggedInvocations(c echo.Context) error {
	invocations := a.service.GetCachedInvocations()
	return c.JSON(http.StatusOK, invocations)
}

func (a AdminController) GetInvocationById(c echo.Context) error {
	id := c.QueryParam("invocationId")
	if id == "" {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("invalid id"))
	}

	invocation := a.service.GetById(id)
	if invocation == nil {
		return c.NoContent(http.StatusNotFound)
	}

	return c.JSON(http.StatusOK, invocation)
}

func (a AdminController) ResetAllCaches(c echo.Context) error {
	err := a.service.ResetAll()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusOK)
}

func (a AdminController) Push(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	id, err := a.service.PushInvocation(string(body))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, &NewInvocationResponseDto{Id: id})
}
