package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type WorkflowRepository interface {
	createWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error)
	updateWorkflow(id string, workflow dto.WorkflowPayload) (*models.Workflows, error)
}

type WorkflowRepositoryImpl struct {
	*sqlx.DB
}

func NewWorkflowRepositoryImple(db *sqlx.DB) WorkflowRepository {
	return &WorkflowRepositoryImpl{
		DB: db,
	}
}

// function for creating a workflow:
func (w *WorkflowRepositoryImpl) createWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error) {
	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		`INSERT INTO workflows (name, description, trigger_type) 
		VALUES ($1, $2, $3) 
		RETURNING *`, workflow.Name, workflow.Description, workflow.TriggerType,
	)
}

// updateWorkflow implements WorkflowRepository.
func (w *WorkflowRepositoryImpl) updateWorkflow(id string, workflow dto.WorkflowPayload) (*models.Workflows, error) {

	setQuery := GenerateKeyValueQuery(map[string]types.Nullable[any]{
		"name":         workflow.Name.ToNullableAny(),
		"description":  workflow.Description.ToNullableAny(),
		"trigger_type": workflow.TriggerType.ToNullableAny(),
	})

	return DbExecAndReturnOne[models.Workflows](
		w.DB,
		fmt.Sprintf(`UPDATE workflows
		SET %v
		WHERE id = $1
		RETURNING *`, setQuery), id,
	)
}
