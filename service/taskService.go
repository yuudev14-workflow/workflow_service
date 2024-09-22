package service

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
)

type TaskService interface {
	GetTasksByWorkflowId(workflowId string) []models.Tasks
	UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error)
	DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error
}

type TaskServiceImpl struct {
	*sqlx.DB
}

func NewTaskServiceImpl(db *sqlx.DB) TaskService {
	return &TaskServiceImpl{
		DB: db,
	}
}

// get tasks by workflow id
func (t *TaskServiceImpl) GetTasksByWorkflowId(workflowId string) []models.Tasks {
	result, _ := repository.DbExecAndReturnMany[models.Tasks](
		t,
		queries.GET_TASK_BY_WORKFLOW_ID,
		workflowId,
	)
	return result
}

// upsert tasks. insert multiple tasks.
// if task does not exist yet add the task in the database
// else update the content of the task
func (t *TaskServiceImpl) UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error) {

	statement := sq.Insert("tasks").Columns("workflow_id", "name", "description")

	for _, val := range tasks {
		statement = statement.Values(workflowId, val.Name, val.Description)
	}

	sql, args, err := statement.Suffix(`
		ON CONFLICT (workflow_id, name) DO UPDATE
   	SET description = EXCLUDED.description,
       parameters = EXCLUDED.parameters,
       updated_at = NOW()`).ToSql()

	logging.Logger.Debug("UpsertTasks SQL: ", sql)
	logging.Logger.Debug("UpsertTasks Args: ", args)

	if err != nil {
		logging.Logger.Error("Failed to build SQL query", err)
		return nil, err
	}

	return repository.DbExecAndReturnMany[models.Tasks](
		tx,
		sql,
		args...,
	)
}

// Delete multiple tasks based on the taskIds
func (t *TaskServiceImpl) DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error {
	sql, args, err := sq.Delete("tasks").Where(sq.Eq{"id": taskIds}).ToSql()
	logging.Logger.Debug("DeleteTasks SQL: ", sql)
	logging.Logger.Debug("DeleteTasks Args: ", args)
	if err != nil {
		logging.Logger.Error("Failed to build SQL query", err)
		return err
	}
	sql = tx.Rebind(sql)
	_, err = tx.Query(sql, args...)
	logging.Logger.Warn(err)

	return err
}
