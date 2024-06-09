package service

import (
	"github.com/yuudev14-workflow/workflow-service/dto"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
)

type WorkflowService interface {
	CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error)
	UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error)
}

type WorkflowServiceImpl struct {
	WorkflowRepository repository.WorkflowRepository
}

func NewWorkflowService(WorkflowRepository repository.WorkflowRepository) WorkflowService {
	return &WorkflowServiceImpl{
		WorkflowRepository: WorkflowRepository,
	}
}

// CreateWorkflow implements WorkflowService.
func (w *WorkflowServiceImpl) CreateWorkflow(workflow dto.WorkflowPayload) (*models.Workflows, error) {
	return w.WorkflowRepository.CreateWorkflow(workflow)
}

// UpdateWorkflow implements WorkflowService.
func (w *WorkflowServiceImpl) UpdateWorkflow(id string, workflow dto.UpdateWorkflowData) (*models.Workflows, error) {
	return w.WorkflowRepository.UpdateWorkflow(id, workflow)
}
