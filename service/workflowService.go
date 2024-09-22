package service

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type WorkflowService interface {
	CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error)
	UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error)
}

type WorkflowServiceImpl struct {
	*sqlx.DB
}

func NewWorkflowService(db *sqlx.DB) WorkflowService {
	return &WorkflowServiceImpl{
		DB: db,
	}
}

// function for creating a workflow:
func (w *WorkflowServiceImpl) CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error) {
	return repository.DbExecAndReturnOne[models.Workflows](
		w.DB,
		queries.INSERT_WORKFLOW, workflow.Name, workflow.Description, workflow.TriggerType,
	)
}

// updateWorkflow implements WorkflowRepository.
func (w *WorkflowServiceImpl) UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error) {

	data := repository.GenerateKeyValueQuery(map[string]types.Nullable[any]{
		"name":         workflow.Name.ToNullableAny(),
		"description":  workflow.Description.ToNullableAny(),
		"trigger_type": workflow.TriggerType.ToNullableAny(),
	})

	sql, args, err := sq.Update("workflows").SetMap(data).Where(sq.Eq{"id": id}).Suffix("RETURNING *").ToSql()

	logging.Logger.Debug("UpdateWorkflow SQL: ", sql)
	logging.Logger.Debug("UpdateWorkflow Args: ", args)
	if err != nil {
		logging.Logger.Error("Failed to build SQL query", err)
		return nil, err
	}
	return repository.DbExecAndReturnOne[models.Workflows](
		w.DB,
		sql,
		args...,
	)
}
