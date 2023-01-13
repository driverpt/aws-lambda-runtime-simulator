package main

import (
	"fmt"
	"lambda-runtime-simulator/pkg/config"
	"lambda-runtime-simulator/pkg/controller"
	"lambda-runtime-simulator/pkg/event"
	"lambda-runtime-simulator/pkg/server"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
)

type EnvironmentVariables struct {
	Port          int    `envconfig:"port" default:"8080"`
	LambdaTimeout int    `envconfig:"function_timeout" default:"120"`
	Arn           string `envconfig:"lambda_arn"`
}

func main() {
	var env EnvironmentVariables
	err := envconfig.Process("", &env)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	controllers := setupControllers(&env)

	server.SetupServer(e, controllers...)

	if err := e.Start(fmt.Sprint("0.0.0.0:", env.Port)); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func setupControllers(env *EnvironmentVariables) []controller.Controller {
	var result []controller.Controller

	cfg := config.Runtime{
		Port:             env.Port,
		TimeoutInSeconds: env.LambdaTimeout,
		Arn:              env.Arn,
	}

	svc := event.NewService(&cfg)

	runtimeController := controller.NewRuntimeController(cfg, svc)
	result = append(result, runtimeController)

	adminController := controller.NewAdminController(svc)
	result = append(result, adminController)

	return result
}
