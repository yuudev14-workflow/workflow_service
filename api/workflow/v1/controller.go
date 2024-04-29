package workflow_controller_v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/pkg/mq"
	rest "github.com/yuudev14-workflow/workflow-service/pkg/rests"
)

type WorkflowController struct {
}

func NewWorkflowController() *WorkflowController {
	return &WorkflowController{}
}

func (w *WorkflowController) Trigger(c *gin.Context) {
	response := rest.Response{C: c}

	// Publish a message to the queue
	body := "hello world"
	err := mq.MQChannel.Publish(
		"",                  // exchange
		mq.SenderQueue.Name, // routing key
		false,               // mandatory
		false,               // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		response.ResponseError(http.StatusBadGateway, err.Error())
	}

	response.ResponseSuccess(gin.H{})
}
