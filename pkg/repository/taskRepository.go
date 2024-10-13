package repository

import (
	"encoding/json"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

type TaskRepository interface {
	GetTasksByWorkflowId(workflowId string) []models.Tasks
	UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error)
	DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error
}

type TaskRepositoryImpl struct {
	*sqlx.DB
}

func NewTaskRepositoryImpl(db *sqlx.DB) TaskRepository {
	return &TaskRepositoryImpl{
		DB: db,
	}
}

// get tasks by workflow id
func (t *TaskRepositoryImpl) GetTasksByWorkflowId(workflowId string) []models.Tasks {
	result, _ := DbExecAndReturnMany[models.Tasks](
		t,
		queries.GET_TASK_BY_WORKFLOW_ID,
		workflowId,
	)
	return result
}

// upsert tasks. insert multiple tasks.
// if task does not exist yet add the task in the database
// else update the content of the task
func (t *TaskRepositoryImpl) UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error) {

	statement := sq.Insert("tasks").Columns("workflow_id", "name", "description", "parameters")

	for _, val := range tasks {
		parameters, _ := json.Marshal(val.Parameters)
		statement = statement.Values(workflowId, val.Name, val.Description, parameters)
	}

	sql, args, err := statement.Suffix(`
		ON CONFLICT (workflow_id, name) DO UPDATE
   	SET description = EXCLUDED.description,
       parameters = EXCLUDED.parameters,
       updated_at = NOW()
		RETURNING *`).ToSql()

	logging.Sugar.Debug("UpsertTasks SQL: ", sql)
	logging.Sugar.Debug("UpsertTasks Args: ", args)

	if err != nil {
		logging.Sugar.Error("Failed to build SQL query", err)
		return nil, err
	}

	return DbExecAndReturnMany[models.Tasks](
		tx,
		sql,
		args...,
	)
}

// Delete multiple tasks based on the taskIds
func (t *TaskRepositoryImpl) DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error {
	sql, args, err := sq.Delete("tasks").Where(sq.Eq{"id": taskIds}).ToSql()
	logging.Sugar.Debug("DeleteTasks SQL: ", sql)
	logging.Sugar.Debug("DeleteTasks Args: ", args)
	if err != nil {
		logging.Sugar.Error("Failed to build SQL query", err)
		return err
	}
	sql = tx.Rebind(sql)
	_, err = tx.Query(sql, args...)
	logging.Sugar.Warn(err)

	return err
}
