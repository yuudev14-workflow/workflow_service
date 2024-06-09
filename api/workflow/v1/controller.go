package workflow_controller_v1

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/pkg/mq"
	rest "github.com/yuudev14-workflow/workflow-service/pkg/rests"
	"github.com/yuudev14-workflow/workflow-service/pkg/utils"
)

type WorkflowController struct {
}

func NewWorkflowController() *WorkflowController {
	return &WorkflowController{}
}

func (w *WorkflowController) UpdateWorkflow(c *gin.Context) {
	var payload dto.WorkflowPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"name": payload.Name.Value,
		"set":  payload.Name.Set,
	})

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
