package repository

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type WorkflowRepository interface {
	GetWorkflowById(id string) (*models.Workflows, error)
	CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error)
	UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error)
}

type WorkflowRepositoryImpl struct {
	*sqlx.DB
}

func NewWorkflowRepository(db *sqlx.DB) WorkflowRepository {
	return &WorkflowRepositoryImpl{
		DB: db,
	}
}

// GetWorkflowById implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) GetWorkflowById(id string) (*models.Workflows, error) {
	sql, args, err := sq.Select("*").From("workflows").Where("id = ?", id).ToSql()
	logging.Sugar.Debugw("GetWorkflowById statement", "sql", sql, "args", args)
	if err != nil {
		logging.Sugar.Error("Error in GetWorkflowById", err)
		return nil, err
	}
	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		sql,
		args...,
	)
}

// function for creating a workflow:
func (w *WorkflowRepositoryImpl) CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error) {
	return DbExecAndReturnOne[models.Workflows](
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

	sql, args, err := sq.Update("workflows").SetMap(data).Where(sq.Eq{"id": id}).Suffix("RETURNING *").ToSql()

	logging.Sugar.Debug("UpdateWorkflow SQL: ", sql)
	logging.Sugar.Debug("UpdateWorkflow Args: ", args)
	if err != nil {
		logging.Sugar.Error("Failed to build SQL query", err)
		return nil, err
	}
	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		sql,
		args...,
	)
}
