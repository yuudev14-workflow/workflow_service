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
	TaskRepository repository.TaskRepository
}

func NewTaskServiceImpl(TaskService repository.TaskRepository) TaskService {
	return &TaskServiceImpl{
		TaskRepository: TaskService,
	}
}

// get tasks by workflow id
func (t *TaskServiceImpl) GetTasksByWorkflowId(workflowId string) []models.Tasks {
	return t.TaskRepository.GetTasksByWorkflowId(workflowId)
}

// upsert tasks. insert multiple tasks.
// if task does not exist yet add the task in the database
// else update the content of the task
func (t *TaskServiceImpl) UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error) {
	return t.TaskRepository.UpsertTasks(tx, workflowId, tasks)
}

// Delete multiple tasks based on the taskIds
func (t *TaskServiceImpl) DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error {
	return t.TaskRepository.DeleteTasks(tx, taskIds)
}
