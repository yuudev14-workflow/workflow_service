package workflow_api

import (
	"github.com/gin-gonic/gin"
	workflow_controller_v1 "github.com/yuudev14-workflow/workflow-service/api/workflow/v1"
	"github.com/yuudev14-workflow/workflow-service/db"
	"github.com/yuudev14-workflow/workflow-service/service"
)

func SetupWorkflowController(route *gin.RouterGroup) {
	workflowService := service.NewWorkflowService(db.DB)
	nodeService := service.NewEdgeServiceImpl(db.DB)
	taskService := service.NewTaskServiceImpl(db.DB)
	workflowController := workflow_controller_v1.NewWorkflowController(workflowService, taskService, nodeService)

	r := route.Group("v1/workflows")
	{
		r.POST("/trigger", workflowController.Trigger)
		r.POST("/", workflowController.CreateWorkflow)
		r.PUT("/:workflow_id", workflowController.UpdateWorkflow)
		r.PUT("/tasks/:workflow_id", workflowController.UpdateWorkflowTasks)
	}
}
