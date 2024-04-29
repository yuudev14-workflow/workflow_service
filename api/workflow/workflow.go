package workflow_api

import (
	"github.com/gin-gonic/gin"
	workflow_controller_v1 "github.com/yuudev14-workflow/workflow-service/api/workflow/v1"
)

func SetupWorkflowController(route *gin.RouterGroup) {
	workflowController := workflow_controller_v1.NewWorkflowController()

	r := route.Group("workflows")
	{
		r.POST("/v1/trigger", workflowController.Trigger)
	}
}
