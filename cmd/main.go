package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"lambda-runtime-simulator/pkg/config"
	"lambda-runtime-simulator/pkg/controller"
	"lambda-runtime-simulator/pkg/event"
	"lambda-runtime-simulator/pkg/server"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type EnvironmentVariables struct {
	Port          int    `envconfig:"port" default:"8080"`
	LambdaTimeout int    `envconfig:"function_timeout" default:"120"`
	Arn           string `envconfig:"lambda_arn"`
	LogLevel      string `envconfig:"log_level" default:"info"`
}

func main() {
	var env EnvironmentVariables
	err := envconfig.Process("", &env)
	if err != nil {
		log.Panicf("Error while processing environment variables: %v", err.Error())
	}

	// Init Logger
	lvl, err := logrus.ParseLevel(env.LogLevel)
	if err != nil {
		log.Panicf("Invalid Error Level [%v]: %v", env.LogLevel, err.Error())
	}

	logrus.SetLevel(lvl)

	// Init Http Server
	e := echo.New()
	controllers := setupControllers(&env)

	err = server.SetupServer(e, controllers...)
	if err != nil {
		log.Panicf("Error while initializing HTTP Server", err.Error())
	}

	if err := http.ListenAndServe(fmt.Sprint("0.0.0.0:", env.Port), e); err != http.ErrServerClosed {
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
