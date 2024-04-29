package main

import (
	"github.com/yuudev14-workflow/workflow-service/api"
	"github.com/yuudev14-workflow/workflow-service/environment"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/mq"
)

// @title 	Workflow Service API
// @version	1.0
// @description A Workflow Service in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func initApp() {
	environment.Setup()
	logging.Setup(environment.Settings.LOGGER_MODE)
}
func main() {
	initApp()
	mq.ConnectToMQ()
	app := api.InitRouter()
	app.Run()
}
