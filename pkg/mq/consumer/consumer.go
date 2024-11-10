package consumer

import (
	"encoding/json"
	"fmt"

	"github.com/yuudev14-workflow/workflow-service/db"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/mq"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
	"github.com/yuudev14-workflow/workflow-service/service"
)

type MessageBody struct {
	Action string                 `json:"action"`
	Params map[string]interface{} `json:"params"`
}

type TaskStatusPayload struct {
	WorkflowHistoryId string      `json:"workflow_history_id"`
	TaskId            string      `json:"task_id"`
	Status            string      `json:"status"`
	Result            interface{} `json:"result,omitempty"`
	Error             *string     `json:"error,omitempty"`
}

type WorkflowStatusPayload struct {
	WorkflowHistoryId string      `json:"workflow_history_id"`
	Status            string      `json:"status"`
	Result            interface{} `json:"result,omitempty"`
	Error             *string     `json:"error,omitempty"`
}

type ConsumeMessage struct {
	WorkflowService service.WorkflowService
	TaskService     service.TaskService
}

func NewConsumeMessage(
	WorkflowService service.WorkflowService,
	TaskService service.TaskService,
) *ConsumeMessage {
	return &ConsumeMessage{
		WorkflowService: WorkflowService,
		TaskService:     TaskService,
	}
}

// Example handler functions for different message types
func (c *ConsumeMessage) handleTask(params []byte) {
	var taskParams TaskStatusPayload
	if err := json.Unmarshal(params, &taskParams); err != nil {
		logging.Sugar.Error("Error unmarshalling task params:", err)
		return
	}
	c.TaskService.UpdateTaskStatus(taskParams.WorkflowHistoryId, taskParams.TaskId, taskParams.Status)
}

func (c *ConsumeMessage) handleWorkflow(params []byte) {
	var workflowParams WorkflowStatusPayload
	if err := json.Unmarshal(params, &workflowParams); err != nil {
		logging.Sugar.Error("Error unmarshalling workflow params:", err)
		return
	}
	c.WorkflowService.UpdateWorkflowHistoryStatus(workflowParams.WorkflowHistoryId, workflowParams.Status)

}

func (c *ConsumeMessage) PrepareMessage(data MessageBody) {

	jsonData, err := json.Marshal(data.Params)
	if err != nil {
		fmt.Println("Error marshalling map to JSON:", err)
		return
	}

	switch data.Action {
	case "workflow_status":
		c.handleWorkflow(jsonData)
	case "task_status":
		c.handleTask(jsonData)

	}

}

func Listen() {

	msgs, err := mq.MQChannel.Consume(
		mq.ReceiverQueue.Name, // queue
		"",                    // consumer
		true,                  // auto-acknowledge
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // arguments
	)

	if err != nil {
		panic("error in consuming a queue")
	}

	var forever chan struct{}

	go func() {
		workflowRepository := repository.NewWorkflowRepository(db.DB)
		taskRepository := repository.NewTaskRepositoryImpl(db.DB)
		workflowService := service.NewWorkflowService(workflowRepository)
		taskService := service.NewTaskServiceImpl(taskRepository, workflowService)
		consumeMessageService := NewConsumeMessage(workflowService, taskService)
		for d := range msgs {
			var message MessageBody

			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				logging.Sugar.Warnf("Error decoding JSON: %v", err)
			}
			logging.Sugar.Infof("Received a message: %s", d.Body)
			consumeMessageService.PrepareMessage(message)
		}
	}()

	logging.Sugar.Info("Listening to message queue")
	<-forever
}
