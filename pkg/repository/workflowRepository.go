package repository

import (
	"encoding/json"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type WorkflowRepository interface {
	GetWorkflows(offset int, limit int) ([]models.Workflows, error)
	GetWorkflowById(id string) (*models.Workflows, error)
	CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error)
	UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error)
	CreateWorkflowHistory(tx *sqlx.Tx, id string) (*models.WorkflowHistory, error)
	UpdateWorkflowHistoryStatus(workflow_history_id string, status string) (*models.WorkflowHistory, error)
	UpdateWorkflowHistory(workflowHistoryId string, workflowHistory dto.UpdateWorkflowHistoryData) (*models.WorkflowHistory, error)
}

type WorkflowRepositoryImpl struct {
	*sqlx.DB
}

func NewWorkflowRepository(db *sqlx.DB) WorkflowRepository {
	return &WorkflowRepositoryImpl{
		DB: db,
	}
}

// GetWorkflows implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) GetWorkflows(offset int, limit int) ([]models.Workflows, error) {
	statement := sq.Select("*").From("workflows").Offset(uint64(offset)).Limit(uint64(limit))
	return DbExecAndReturnMany[models.Workflows](
		w.DB,
		statement,
	)
}

// CreateWorkflowHistory implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) CreateWorkflowHistory(tx *sqlx.Tx, id string) (*models.WorkflowHistory, error) {
	statement := sq.Insert("workflow_history").Columns("workflow_id", "triggered_at").Values(id, time.Now()).Suffix("RETURNING *")
	return DbExecAndReturnOne[models.WorkflowHistory](
		tx,
		statement,
	)
}

// GetWorkflowById implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) GetWorkflowById(id string) (*models.Workflows, error) {
	statement := sq.Select("*").From("workflows").Where("id = ?", id)
	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		statement,
	)
}

// function for creating a workflow:
func (w *WorkflowRepositoryImpl) CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error) {

	return DbExecAndReturnOneOld[models.Workflows](
		w.DB,
		queries.INSERT_WORKFLOW, workflow.Name, workflow.Description, workflow.TriggerType,
	)
}

// updateWorkflow implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error) {

	data := GenerateKeyValueQuery(map[string]types.Nullable[any]{
		"name":         workflow.Name.ToNullableAny(),
		"description":  workflow.Description.ToNullableAny(),
		"trigger_type": workflow.TriggerType.ToNullableAny(),
	})

	statement := sq.Update("workflows").SetMap(data).Where(sq.Eq{"id": id}).Suffix("RETURNING *")

	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		statement,
	)
}

// UpdateWorkflowHistory implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) UpdateWorkflowHistory(workflowHistoryId string, workflowHistory dto.UpdateWorkflowHistoryData) (*models.WorkflowHistory, error) {
	data := GenerateKeyValueQuery(map[string]types.Nullable[any]{
		"status": workflowHistory.Status.ToNullableAny(),
		"error":  workflowHistory.Error.ToNullableAny(),
	})

	jsonData, err := json.Marshal(workflowHistory.Result)
	if err != nil {
		return nil, err
	}

	data["result"] = jsonData
	statement := sq.Update("workflow_history").SetMap(data).Where(sq.Eq{"id": workflowHistoryId}).Suffix("RETURNING *")
	return DbExecAndReturnOne[models.WorkflowHistory](
		w.DB,
		statement,
	)
}

// UpdateWorkflowHistoryStatus implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) UpdateWorkflowHistoryStatus(workflowHistoryId string, status string) (*models.WorkflowHistory, error) {
	statement := sq.Update("workflow_history").Set("status", status).Where(sq.Eq{"id": workflowHistoryId}).Suffix("RETURNING *")
	return DbExecAndReturnOne[models.WorkflowHistory](
		w.DB,
		statement,
	)
}
