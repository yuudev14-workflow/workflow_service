package main

import "github.com/yuudev14-workflow/workflow-service/api"

// @title 	Workflow Service API
// @version	1.0
// @description A Workflow Service in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	app := api.InitRouter()
	app.Run()
}
