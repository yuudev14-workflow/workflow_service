package workflow_controller_v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/mq"
	rest "github.com/yuudev14-workflow/workflow-service/pkg/rests"
	"github.com/yuudev14-workflow/workflow-service/pkg/utils"
	"github.com/yuudev14-workflow/workflow-service/service"
)

type WorkflowController struct {
	WorkflowService service.WorkflowService
}

func NewWorkflowController(WorkflowService service.WorkflowService) *WorkflowController {
	return &WorkflowController{
		WorkflowService: WorkflowService,
	}
}

func (w *WorkflowController) CreateWorkflow(c *gin.Context) {
	var body dto.WorkflowPayload
	response := rest.Response{C: c}

	check, code, validErr := rest.BindFormAndValidate(c, &body)

	if !check {
		logging.Logger.Errorf(fmt.Sprintf("%v", validErr))
		response.ResponseError(code, validErr)
		return
	}

	workflow, err := w.WorkflowService.CreateWorkflow(body)

	logging.Logger.Debug("added workflow...")

	if err != nil {
		response.ResponseError(http.StatusBadRequest, err.Error())
		return
	}

	response.ResponseSuccess(workflow)

}

func (w *WorkflowController) UpdateWorkflow(c *gin.Context) {
	var body dto.UpdateWorkflowData
	response := rest.Response{C: c}
	workflowId := c.Param("workflow_id")

	check, code, validErr := rest.BindFormAndValidate(c, &body)

	if !check {
		logging.Logger.Errorf(fmt.Sprintf("%v", validErr))
		response.ResponseError(code, validErr)
		return
	}

	// node to insert
	// node to update
	// node to delete

	// 1. check if node exists in the existing nodes else update
	// 2. if node not in new nodes to be updated, delete

	// get all edges
	// 1. check if node exists in the existing nodes else update
	// 2. if node not in new nodes to be updated, delete

	// save everything

	workflow, err := w.WorkflowService.UpdateWorkflow(workflowId, body)

	logging.Logger.Debug("added workflow...")

	if err != nil {
		response.ResponseError(http.StatusBadRequest, err.Error())
		return
	}

	response.ResponseSuccess(workflow)

}

func (w *WorkflowController) Trigger(c *gin.Context) {
	response := rest.Response{C: c}
	graph := map[string][]string{
		"A": {"B"},
		"B": {"C"},
		"C": {},
	}

	currentNode := "A"

	currentQueue := []string{
		"A",
	}

	// Publish a message to the queue
	body := utils.WorkflowData{
		Graph:        graph,
		CurrentNode:  currentNode,
		CurrentQueue: currentQueue,
		Visited:      currentQueue,
	}

	jsonData, jsonErr := json.Marshal(body)

	if jsonErr != nil {
		response.ResponseError(http.StatusBadGateway, jsonErr.Error())
	}
	err := mq.MQChannel.Publish(
		"",                  // exchange
		mq.SenderQueue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(jsonData),
		})
	if err != nil {
		response.ResponseError(http.StatusBadGateway, err.Error())
	}

	response.ResponseSuccess(gin.H{})
}
