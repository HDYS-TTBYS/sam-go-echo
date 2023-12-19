package main

import (
	"context"
	"log"
	"net/http"
	"sam-go-echo/api/healthcheck"
	"sam-go-echo/handler/auth"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	echoadapter "github.com/awslabs/aws-lambda-go-api-proxy/echo"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var echoLambda *echoadapter.EchoLambda

func init() {
	log.Printf("echo cold start")

	e := echo.New()
	e.Use(middleware.Recover())
	e.GET("/api/healthcheck", healthcheck.Healthcheck)
	r := e.Group("/restricted")
	{
		r.Use(auth.BasicAuth())
		r.GET("/", func(c echo.Context) error {
			return c.String(http.StatusOK, "Hello, World!")
		})
	}

	echoLambda = echoadapter.New(e)
}

func main() {
	lambda.Start(Handler)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return echoLambda.ProxyWithContext(ctx, req)
}
