package controller_test

import (
	"encoding/json"
	"io"
	"lambda-runtime-simulator/pkg/config"
	"lambda-runtime-simulator/pkg/controller"
	"lambda-runtime-simulator/pkg/event"
	"lambda-runtime-simulator/pkg/server"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestInvocationCreation(t *testing.T) {
	e := echo.New()
	s := event.NewService(&config.Runtime{})
	c := controller.NewAdminController(s)

	err := server.SetupServer(e, c)
	assert.NoError(t, err)

	body := "{}"

	req := httptest.NewRequest(http.MethodPost, "/push", strings.NewReader(body))
	rec := httptest.NewRecorder()
	e.Server.Handler.ServeHTTP(rec, req)

	assert.Equal(t, rec.Result().StatusCode, http.StatusOK)
	b, err := io.ReadAll(rec.Result().Body)
	assert.NoError(t, err)

	var result controller.NewInvocationResponseDto
	err = json.Unmarshal(b, &result)
	assert.NoError(t, err)
	assert.NotNil(t, &result)
	assert.NotEmpty(t, result.Id)
}
