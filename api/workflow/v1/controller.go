package workflow_controller_v1

import (
	"github.com/gin-gonic/gin"
	rest "github.com/yuudev14-workflow/workflow-service/pkg/rests"
)

type WorkflowController struct {
}

func NewWorkflowController() *WorkflowController {
	return &WorkflowController{}
}

func (w *WorkflowController) Trigger(c *gin.Context) {
	response := rest.Response{C: c}

	response.ResponseSuccess(gin.H{})
}
