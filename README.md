![Workflow](https://github.com/driverpt/aws-lambda-runtime-simulator/actions/workflows/pr.yml/badge.svg)

# aws-lambda-runtime-simulator

## ⚠️ This Software is meant only for local testing purposes ⚠️
1. Use it at your own risk
2. Contributions Accepted

## Why was this project created?
I develop a lot of projects using AWS Lambda, and one of my biggest frustrations is the fact
that it's not possible to Test End-To-End the Lambda call (e.g.: Simulate the real Lambda Call).
Currently, AWS provides a [Lambda Runtime Interface Emulator](https://github.com/aws/aws-lambda-runtime-interface-emulator)
for testing purposes, but you're required to create a Docker Image and wrap your "executable" into the Container Entrypoint.
The goal of this is to use [TestContainers](https://www.testcontainers.org/) and just start your lambda Handler with `AWS_LAMBDA_RUNTIME_API`
pointing to this container.

## How to use
This package is also available as Docker Container.

### TL;DR
```shell
docker run -p 8080:8080 -d driverpt/lambda-api-simulator
```

Start your handler by setting `AWS_LAMBDA_RUNTIME_API` Env Var to `localhost:8080`

Send invocations
```shell
curl -X POST localhost:8080/invocation -d <Invocation Payload>
```

Will return something like
```json
{
  "InvocationId": "<UUID>"
}
```

Access invocation result with
```shell
curl localhost:8080/invocation/<InvocationId>
```

### Additional options
If you do not want to use Docker Container, go to the Releases Page and download the proper package for your System Arch.

### Configurations
Everything is configured via environment variables

```shell
SERVER_PORT=<NewServerPort> # Defaults to 8080
FUNCTION_TIMEOUT=<FunctionTimeoutInSeconds> # Defaults to 120, simulates Timeouts
ARN=<LambdaArn> # AWS Lambda ARN that you want in the HTTP Headers
LOG_LEVEL=<LogLevel> # Log verbosity, defaults to INFO
```

## Limitations
This only simulates AWS Lambda Runtime API Calls, does not simulate Lambda Interruptions.
If you need to simulate interruptions, I strongly recommend using official [Lambda Runtime Interface Emulator](https://github.com/aws/aws-lambda-runtime-interface-emulator)
