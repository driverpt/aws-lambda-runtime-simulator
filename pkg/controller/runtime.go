package controller

import (
	"fmt"
	"io"
	"lambda-runtime-simulator/pkg/config"
	"lambda-runtime-simulator/pkg/event"
	"lambda-runtime-simulator/pkg/lambda"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type RuntimeController struct {
	config  config.Runtime
	service *event.Service
}

var _ Controller = &RuntimeController{}

func NewRuntimeController(cfg config.Runtime, service *event.Service) *RuntimeController {
	return &RuntimeController{
		config:  cfg,
		service: service,
	}
}

func (r RuntimeController) RegisterRoutes(e *echo.Echo) error {
	e.GET("/2018-06-01/runtime/invocation/next", r.NextInvocation)
	e.POST("/2018-06-01/runtime/invocation/:requestId/response", r.SendResponse)
	e.POST("/2018-06-01/runtime/invocation/:requestId/error", r.SendError)
	e.POST("/2018-06-01/runtime/init/error", r.SendRuntimeInitError)

	return nil
}

func (r RuntimeController) NextInvocation(c echo.Context) error {
	next, err := r.service.GetNextInvocation()
	if err != nil {
		return err
	}

	if next == nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	c.Response().Header().Set(lambda.HeaderLambdaRuntimeAwsRequestId, next.Id)
	c.Response().Header().Set(lambda.HeaderLambdaRuntimeInvokedFunctionArn, r.config.Arn)
	c.Response().Header().Set(lambda.HeaderLambdaRuntimeDeadlineMs, fmt.Sprint(next.Timeout.Unix()))
	// TODO: Implement Tracing Part
	return c.JSON(http.StatusOK, next.Body)
}

func (r RuntimeController) SendResponse(c echo.Context) error {
	requestId := c.Param("requestId")

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	err = r.service.SendResponse(requestId, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (r RuntimeController) SendError(c echo.Context) error {
	errorType := c.Request().Header.Get(lambda.HeaderLambdaRuntimeFunctionErrorType)
	requestId := c.Param("requestId")

	var body event.RuntimeError
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid body")
	}

	err = r.service.SendError(requestId, &body, errorType)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return err
}

func (r RuntimeController) SendRuntimeInitError(c echo.Context) error {
	hErrorType := c.Request().Header.Get(lambda.HeaderLambdaRuntimeFunctionErrorType)
	log.Warnf("Runtime Init Received: %v", hErrorType)

	var body event.RuntimeError
	err := c.Bind(&body)
	if err != nil {
		return c.String(http.StatusBadRequest, "invalid body")
	}

	err = r.service.SendRuntimeInitError(body.ErrorMessage, body.ErrorType, body.StackTrace...)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return err
}
