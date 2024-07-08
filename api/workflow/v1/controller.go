package workflow_controller_v1

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"github.com/yuudev14-workflow/workflow-service/db"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/mq"
	rest "github.com/yuudev14-workflow/workflow-service/pkg/rests"
	"github.com/yuudev14-workflow/workflow-service/pkg/utils"
	"github.com/yuudev14-workflow/workflow-service/service"
)

type WorkflowController struct {
	WorkflowService service.WorkflowService
	TaskService     service.TaskService
}

func NewWorkflowController(WorkflowService service.WorkflowService, TaskService service.TaskService) *WorkflowController {
	return &WorkflowController{
		WorkflowService: WorkflowService,
		TaskService:     TaskService,
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

func (w *WorkflowController) UpdateWorkflowTasks(c *gin.Context) {
	var body dto.UpdateWorkflowtasks
	response := rest.Response{C: c}
	workflowId := c.Param("workflow_id")
	tx, err := db.DB.Beginx()
	if err != nil {
		tx.Rollback()
		response.ResponseError(http.StatusInternalServerError, err)
		return
	}

	check, code, validErr := rest.BindFormAndValidate(c, &body)

	if !check {
		logging.Logger.Errorf(fmt.Sprintf("%v", validErr))
		response.ResponseError(code, validErr)
		return
	}

	// verify nodes name should be unique
	tasks := w.TaskService.GetTasksByWorkflowId(workflowId)
	logging.Logger.Debugf("tasks: %v", tasks)
	tasksMap := make(map[string]models.Tasks)
	tasksBodyMap := make(map[string]bool)

	for _, task := range tasks {
		tasksMap[task.Name] = task
	}
	for _, node := range body.Nodes {
		tasksBodyMap[node.Name] = true
	}

	// // node to insert
	// var nodeToInsert []dto.Task
	// // node to update
	// var nodeToUpdate []dto.Task

	// node to update
	var nodeToUpsert []models.Tasks
	// node to delete
	var nodeToDelete []uuid.UUID

	// 1. check if node exists in the existing nodes else update
	for _, node := range body.Nodes {
		// _, ok := tasksMap[node.Name]
		// if ok {
		// 	nodeToUpdate = append(nodeToUpdate, node)
		// } else {
		// 	nodeToInsert = append(nodeToInsert, node)
		// }
		nodeToUpsert = append(nodeToUpsert, models.Tasks{
			Name:        node.Name,
			Parameters:  node.Parameters,
			Description: node.Description,
		})
	}
	// 2. if node not in new nodes to be updated, delete

	for _, node := range tasks {
		_, ok := tasksBodyMap[node.Name]
		if !ok {
			nodeToDelete = append(nodeToDelete, node.ID)
		}
	}

	workflowUUID, err := uuid.Parse(workflowId)

	if err != nil {
		response.ResponseError(http.StatusInternalServerError, err)
		return
	}

	logging.Logger.Debugf("node to add: %v", nodeToUpsert)
	// save the tasks
	if len(nodeToUpsert) > 0 {
		w.TaskService.UpsertTasks(tx, workflowUUID, nodeToUpsert)
	}

	if len(nodeToDelete) > 0 {
		w.TaskService.DeleteTasks(tx, nodeToDelete)

	}

	// get all edges
	// convert edges into ids
	// 1. check if node exists in the existing nodes else update
	// 2. if node not in new nodes to be updated, delete

	// save everything

	logging.Logger.Debug("added workflow...")
	commitErr := tx.Commit()

	if commitErr != nil {
		response.ResponseError(http.StatusInternalServerError, commitErr)
		return
	}

	newTasks := w.TaskService.GetTasksByWorkflowId(workflowId)
	response.Response(http.StatusAccepted, newTasks)
	return
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
