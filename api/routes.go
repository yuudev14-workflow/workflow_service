package api

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/yuudev14-workflow/workflow-service/docs"
)

func StartApi(app *gin.RouterGroup) {

}

func InitRouter() *gin.Engine {

	app := gin.Default()

	docs.SwaggerInfo.BasePath = "./"

	api_group := app.Group("/api")

	StartApi(api_group)

	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return app

}
