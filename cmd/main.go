package main

import (
	"local.com/go-clean-lambda/internal/app"
	"local.com/go-clean-lambda/internal/controller"
	"local.com/go-clean-lambda/internal/logger"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gorillamux"
)

func main() {
	controllers, err := app.InitDummyControllers()
	if err != nil {
		logger.Error("execution end. failed to init lambda.", err)
		return
	}
	r := controller.NewRouter(controllers)
	logger.Info("muxproxy initialization done")
	lambda.Start(gorillamux.New(r).ProxyWithContext)
}
