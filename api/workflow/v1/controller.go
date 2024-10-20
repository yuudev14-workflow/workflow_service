package workflow_controller_v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	rest "github.com/yuudev14-workflow/workflow-service/pkg/rests"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
	"github.com/yuudev14-workflow/workflow-service/service"
)

type WorkflowController struct {
	WorkflowService        service.WorkflowService
	TaskService            service.TaskService
	EdgeService            service.EdgeService
	WorkflowTriggerService service.WorkflowTriggerService
}

func NewWorkflowController(
	WorkflowService service.WorkflowService,
	TaskService service.TaskService,
	EdgeService service.EdgeService,
	WorkflowTriggerService service.WorkflowTriggerService,
) *WorkflowController {
	return &WorkflowController{
		WorkflowService:        WorkflowService,
		TaskService:            TaskService,
		EdgeService:            EdgeService,
		WorkflowTriggerService: WorkflowTriggerService,
	}
}

func (w *WorkflowController) CreateWorkflow(c *gin.Context) {
	var body dto.WorkflowPayload
	response := rest.Response{C: c}

	check, code, validErr := rest.BindFormAndValidate(c, &body)

	if !check {
		logging.Sugar.Errorf(fmt.Sprintf("%v", validErr))
		response.ResponseError(code, validErr)
		return
	}

	workflow, err := w.WorkflowService.CreateWorkflow(body)

	logging.Sugar.Debug("added workflow...")

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
		logging.Sugar.Errorf(fmt.Sprintf("%v", validErr))
		response.ResponseError(code, validErr)
		return
	}

	workflow, err := w.WorkflowService.UpdateWorkflow(workflowId, body)

	logging.Sugar.Debug("added workflow...")

	if err != nil {
		response.ResponseError(http.StatusBadRequest, err.Error())
		return
	}

	response.ResponseSuccess(workflow)

}

func (w *WorkflowController) UpsertTasks(
	tx *sqlx.Tx,
	workflowUUID uuid.UUID,
	nodes []dto.Task,
) ([]models.Tasks, error) {
	// node to update
	var nodeToUpsert []models.Tasks
	for _, node := range nodes {
		nodeToUpsert = append(nodeToUpsert, models.Tasks{
			Name: node.Name,
			Parameters: func() types.JsonType {
				if node.Parameters != nil {
					return types.JsonType(*node.Parameters)
				}
				return nil
			}(),
			Description:   node.Description,
			Config:        node.Config.Value,
			ConnectorName: node.ConnectorName,
			Operation:     node.Operation,
		})
	}

	logging.Sugar.Debugf("node to add: %v", nodeToUpsert)
	// save the tasks
	if len(nodeToUpsert) > 0 {
		return w.TaskService.UpsertTasks(tx, workflowUUID, nodeToUpsert)
	}
	return nil, nil
}

func (w *WorkflowController) InsertEdges(
	tx *sqlx.Tx,
	workflowUUID uuid.UUID,
	edges map[string][]string,
	tasks []models.Tasks,
) error {
	// node to update
	var edgeToInsert []models.Edges
	tasksMap := make(map[string]uuid.UUID)

	// create a taskmap with name and uuid of the task to easily get the uuid from the edges
	for _, task := range tasks {
		tasksMap[task.Name] = task.ID
	}

	for key, values := range edges {
		for _, val := range values {
			sourceId, sourceIdOk := tasksMap[key]
			destinationID, destinationIdOk := tasksMap[val]
			if sourceIdOk && destinationIdOk {
				edgeToInsert = append(edgeToInsert, models.Edges{
					SourceID:      sourceId.String(),
					DestinationID: destinationID.String(),
				})
			}
		}
	}

	logging.Sugar.Debugf("edges to add: %v", edgeToInsert)
	// save the edges
	if len(edgeToInsert) > 0 {
		_, err := w.EdgeService.InsertEdges(tx, edgeToInsert)
		return err
	}
	return nil
}

func (w *WorkflowController) DeleteTasks(
	tx *sqlx.Tx,
	workflowUUID uuid.UUID,
	nodes []dto.Task,
) error {
	// node to delete
	var nodeToDelete []uuid.UUID
	tasksBodyMap := make(map[string]bool)

	// verify nodes name should be unique
	tasks, tasksErr := w.TaskService.GetTasksByWorkflowId(workflowUUID.String())
	if tasksErr != nil {
		return tasksErr
	}
	logging.Sugar.Debugf("tasks: %v", tasks)

	for _, node := range nodes {
		tasksBodyMap[node.Name] = true
	}
	// 2. if node not in new nodes to be updated, delete
	for _, node := range tasks {
		_, ok := tasksBodyMap[node.Name]
		if !ok {
			nodeToDelete = append(nodeToDelete, node.ID)
		}
	}

	logging.Sugar.Debugf("node to delete: %v", nodeToDelete)
	if len(nodeToDelete) > 0 {
		logging.Sugar.Debugf("node to delete: %v", nodeToDelete)
		err := w.TaskService.DeleteTasks(tx, nodeToDelete)
		return err

	}
	return nil
}

// delete edges that doesnt exist in the body payload
func (w *WorkflowController) DeleteEdges(
	tx *sqlx.Tx,
	workflowUUID uuid.UUID,
	edges map[string][]string,
) error {

	var edgeToDelete []uuid.UUID
	edgesMap := make(map[[2]string]bool)

	// delete all edges from the workflow if nothing is in the payload
	if len(edges) == 0 {
		return w.EdgeService.DeleteAllWorkflowEdges(tx, workflowUUID.String())
	}

	workflowEdges, workflowEdgesErr := w.EdgeService.GetEdgesByWorkflowId(workflowUUID.String())
	logging.Sugar.Debug("workflow edges", workflowEdges)

	if workflowEdgesErr != nil {
		logging.Sugar.Error(workflowEdgesErr)
		return workflowEdgesErr
	}

	// populate the hashmap
	for key, values := range edges {
		for _, val := range values {
			edgesMap[[2]string{key, val}] = true
		}
	}

	// if the edge does not exist in the hashmap, add to the delete lists
	for _, edge := range workflowEdges {
		_, ok := edgesMap[[2]string{edge.SourceTaskName, edge.DestinationTaskName}]
		if !ok {
			edgeToDelete = append(edgeToDelete, edge.ID)
		}
	}

	logging.Sugar.Debugf("edge to delete: %v", edgeToDelete)
	if len(edgeToDelete) > 0 {
		deleteEdgesError := w.EdgeService.DeleteEdges(tx, edgeToDelete)
		return deleteEdgesError

	}
	return nil

}

func validateWorkflowTaskPayload(body dto.UpdateWorkflowtasks) error {
	_, ok := body.Edges["start"]
	if !ok {
		return fmt.Errorf("'Start' doesnt exist in edges")
	}

	for _, node := range body.Nodes {
		if node.Name == "start" {
			return nil
		}
	}

	return fmt.Errorf("'Start' doesnt exist in nodes")
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
		logging.Sugar.Errorf(fmt.Sprintf("%v", validErr))
		response.ResponseError(code, validErr)
		return
	}

	// validate if start node in body payload
	payloadErr := validateWorkflowTaskPayload(body)

	if payloadErr != nil {
		response.ResponseError(code, payloadErr.Error())
		return
	}

	workflowUUID, err := uuid.Parse(workflowId)

	if err != nil {
		response.ResponseError(http.StatusInternalServerError, err.Error())
		return
	}
	deleteEdgesErr := w.DeleteEdges(tx, workflowUUID, body.Edges)
	if deleteEdgesErr != nil {
		logging.Sugar.Error(deleteEdgesErr)
		tx.Rollback()
		response.ResponseError(http.StatusBadRequest, deleteEdgesErr.Error())
		return
	}
	insertedTasks, upsertTasksErr := w.UpsertTasks(tx, workflowUUID, body.Nodes)
	if upsertTasksErr != nil {
		logging.Sugar.Error(upsertTasksErr)
		tx.Rollback()
		response.ResponseError(http.StatusBadRequest, upsertTasksErr.Error())
		return
	}
	deleteTaskError := w.DeleteTasks(tx, workflowUUID, body.Nodes)
	if deleteTaskError != nil {
		logging.Sugar.Error(deleteTaskError)
		tx.Rollback()
		response.ResponseError(http.StatusInternalServerError, deleteTaskError.Error())
		return
	}
	insertEdgeError := w.InsertEdges(tx, workflowUUID, body.Edges, insertedTasks)
	if insertEdgeError != nil {
		logging.Sugar.Error(insertEdgeError)
		tx.Rollback()
		response.ResponseError(http.StatusInternalServerError, insertEdgeError.Error())
		return
	}

	logging.Sugar.Debug("added workflow...")
	commitErr := tx.Commit()

	if commitErr != nil {
		logging.Sugar.Error(commitErr)
		tx.Rollback()
		response.ResponseError(http.StatusInternalServerError, commitErr.Error())
		return
	}

	newTasks, newTaskErr := w.TaskService.GetTasksByWorkflowId(workflowId)
	if newTaskErr != nil {
		logging.Sugar.Errorf("error: ", newTaskErr)
		response.ResponseError(http.StatusBadRequest, newTaskErr.Error())
		return
	}
	newEdges, _ := w.EdgeService.GetEdgesByWorkflowId(workflowId)

	response.Response(http.StatusAccepted, gin.H{
		"tasks": newTasks,
		"edges": newEdges,
	})
}

func (w *WorkflowController) GetTasksByWorkflowId(c *gin.Context) {
	response := rest.Response{C: c}
	workflowId := c.Param("workflow_id")

	_, workflowErr := w.WorkflowService.GetWorkflowById(workflowId)

	if workflowErr != nil {
		logging.Sugar.Error(workflowErr)
		response.ResponseError(http.StatusInternalServerError, workflowErr.Error())
		return
	}

	newTasks, newTaskErr := w.TaskService.GetTasksByWorkflowId(workflowId)
	if newTaskErr != nil {
		logging.Sugar.Errorf("error: ", newTaskErr)
		response.ResponseError(http.StatusBadRequest, newTaskErr.Error())
		return
	}

	response.ResponseSuccess(gin.H{
		"tasks": newTasks,
	})
}

func (w *WorkflowController) Trigger(c *gin.Context) {
	response := rest.Response{C: c}
	workflowId := c.Param("workflow_id")
	triggerErr := w.WorkflowTriggerService.TriggerWorkflow(workflowId)
	if triggerErr != nil {
		logging.Sugar.Errorf("error when sending the message to queue", triggerErr)
		response.ResponseError(http.StatusBadGateway, triggerErr.Error())
		return
	}

	response.Response(http.StatusAccepted, "triggered successfully")
}
