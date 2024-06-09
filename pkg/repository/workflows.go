package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type WorkflowRepository interface {
	CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error)
	UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error)
}

type WorkflowRepositoryImpl struct {
	*sqlx.DB
}

func NewWorkflowRepositoryImpl(db *sqlx.DB) WorkflowRepository {
	return &WorkflowRepositoryImpl{
		DB: db,
	}
}

// function for creating a workflow:
func (w *WorkflowRepositoryImpl) CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error) {
	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		`INSERT INTO workflows (name, description, trigger_type) 
		VALUES ($1, $2, $3) 
		RETURNING *`, workflow.Name, workflow.Description, workflow.TriggerType,
	)
}

// updateWorkflow implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error) {

	values, queries := GenerateKeyValueQuery(map[string]types.Nullable[any]{
		"name":         workflow.Name.ToNullableAny(),
		"description":  workflow.Description.ToNullableAny(),
		"trigger_type": workflow.TriggerType.ToNullableAny(),
	}, 2)

	values = append([]any{id}, values...)
	logging.Logger.Debug("yuuuuqueries", queries)

	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		fmt.Sprintf(`UPDATE workflows
		SET %v
		WHERE id = $1
		RETURNING *`, queries), values...,
	)
}
