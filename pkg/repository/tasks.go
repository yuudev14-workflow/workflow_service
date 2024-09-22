package repository

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

type TasksRepository interface {
	GetTasksByWorkflowId(workflowId string) []models.Tasks
	UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error)
	DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error
}

type TasksRepositoryImpl struct {
	*sqlx.DB
}

func NewTaskRepositoryImpl(db *sqlx.DB) TasksRepository {
	return &TasksRepositoryImpl{
		DB: db,
	}
}

// get tasks by workflow id
func (t *TasksRepositoryImpl) GetTasksByWorkflowId(workflowId string) []models.Tasks {
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
func (t *TasksRepositoryImpl) UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error) {
	var values []string

	for _, val := range tasks {
		values = append(values,
			fmt.Sprintf(`('%v', '%v', '%v')`, workflowId, val.Name, val.Description),
		)
	}

	valueQuery := strings.Join(values, ",")

	statement := fmt.Sprintf(queries.UPSERT_TASK, valueQuery)
	// insertBuilder := sq.Insert("tasks").
	// 	Columns("workflow_id", "name", "description")
	// for _, v := range tasks {
	// 	insertBuilder = insertBuilder.Values(workflowId.String(),
	// 		v.Name,
	// 		v.Description,
	// 	)
	// }

	// conflictBuilder := insertBuilder.Suffix(`
	// ON CONFLICT (workflow_id, name) DO UPDATE
	//   SET description = EXCLUDED.description,
	//       parameters = EXCLUDED.parameters,
	//       updated_at = NOW();`)

	// Build the SQL query and arguments
	// sqlQuery, args, err := insertBuilder.ToSql()
	// if err != nil {
	// 	return nil, err
	// }

	return DbExecAndReturnMany[models.Tasks](
		tx,
		statement,
	)
}

// Delete multiple tasks based on the taskIds
func (t *TasksRepositoryImpl) DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error {
	sql, args, err := sq.Delete("tasks").Where(sq.Eq{"id": taskIds}).ToSql()
	logging.Logger.Debug(sql, args, err)
	sql = tx.Rebind(sql)
	_, err = tx.Query(sql, args...)
	logging.Logger.Warn(err)

	return err
}
