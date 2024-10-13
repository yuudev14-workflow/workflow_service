package auth_api

import (
	"github.com/gin-gonic/gin"
	auth_api_v1 "github.com/yuudev14-workflow/workflow-service/api/auth/v1"
	"github.com/yuudev14-workflow/workflow-service/db"
	"github.com/yuudev14-workflow/workflow-service/service"
)

func SetupAuthController(route *gin.RouterGroup) {
	authService := service.NewUserService(db.DB)
	authController := auth_api_v1.NewAuthController(authService)
	r := route.Group("auth")
	{
		r.POST("/v1/sign-up", authController.SignUp)
		r.POST("/v1/login", authController.Login)
	}
}
