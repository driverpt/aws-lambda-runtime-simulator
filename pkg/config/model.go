package config

const (
	PushApiPort            = "PUSH_API_PORT"
	RuntimeApiPort         = "RUNTIME_API_PORT"
	EnvVariableFunctionArn = "FUNCTION_ARN"
)

type Runtime struct {
	Port             int
	TimeoutInSeconds int
	Arn              string
}
