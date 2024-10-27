package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
)

type WorkflowService interface {
	GetWorkflowById(id string) (*models.Workflows, error)
	CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error)
	UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error)
	CreateWorkflowHistory(tx *sqlx.Tx, id string) (*models.WorkflowHistory, error)
	UpdateWorkflowHistoryStatus(workflowHistoryId string, status string) (*models.WorkflowHistory, error)
}

type WorkflowServiceImpl struct {
	WorkflowRepository repository.WorkflowRepository
}

func NewWorkflowService(WorkflowRepository repository.WorkflowRepository) WorkflowService {
	return &WorkflowServiceImpl{
		WorkflowRepository: WorkflowRepository,
	}
}

// CreateWorkflowHistory implements WorkflowService.
func (w *WorkflowServiceImpl) CreateWorkflowHistory(tx *sqlx.Tx, id string) (*models.WorkflowHistory, error) {
	return w.WorkflowRepository.CreateWorkflowHistory(tx, id)
}

// GetWorkflowById implements WorkflowService.
func (w *WorkflowServiceImpl) GetWorkflowById(id string) (*models.Workflows, error) {
	workflow, workflowErr := w.WorkflowRepository.GetWorkflowById(id)
	if workflowErr != nil {
		return nil, workflowErr
	}

	if workflow == nil {
		return nil, fmt.Errorf("user is not found")
	}
	return workflow, nil
}

// function for creating a workflow:
func (w *WorkflowServiceImpl) CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error) {
	return w.WorkflowRepository.CreateWorkflow(workflow)
}

// updateWorkflow implements WorkflowRepository.
func (w *WorkflowServiceImpl) UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error) {
	return w.WorkflowRepository.UpdateWorkflow(id, workflow)
}

// UpdateWorkflowHistoryStatus implements WorkflowRepository.
func (w *WorkflowServiceImpl) UpdateWorkflowHistoryStatus(workflowHistoryId string, status string) (*models.WorkflowHistory, error) {
	res, err := w.WorkflowRepository.UpdateWorkflowHistoryStatus(workflowHistoryId, status)

	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, fmt.Errorf("no workflow status was updated")
	}

	return res, nil
}
