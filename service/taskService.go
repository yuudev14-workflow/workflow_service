package service

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
)

type TaskService interface {
	GetTasksByWorkflowId(workflowId string) []models.Tasks
	UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error)
	DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error
}

type TaskServiceImpl struct {
	TaskRepository repository.TasksRepository
}

func NewTaskServiceImpl(TaskRepository repository.TasksRepository) TaskService {
	return &TaskServiceImpl{
		TaskRepository: TaskRepository,
	}
}

// DeleteTasks implements TaskService.
func (t *TaskServiceImpl) DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error {
	return t.TaskRepository.DeleteTasks(tx, taskIds)
}

// GetTasksByWorkflowId implements TaskService.
func (t *TaskServiceImpl) GetTasksByWorkflowId(workflowId string) []models.Tasks {
	return t.TaskRepository.GetTasksByWorkflowId(workflowId)
}

// UpsertTasks implements TaskService.
func (t *TaskServiceImpl) UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error) {
	return t.TaskRepository.UpsertTasks(tx, workflowId, tasks)
}
