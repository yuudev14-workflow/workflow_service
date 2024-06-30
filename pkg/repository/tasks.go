package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/models"
)

type TasksRepository interface {
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

// upsert tasks. insert multiple tasks.
// if task does not exist yet add the task in the database
// else update the content of the task
func (t *TasksRepositoryImpl) UpsertTasks(tx *sqlx.Tx, workflowId uuid.UUID, tasks []models.Tasks) ([]models.Tasks, error) {
	var values []string

	for _, val := range tasks {
		values = append(values,
			fmt.Sprintf(`(%v, %v, %v, %v)`, workflowId, val.Name, val.Description, val.Parameters),
		)
	}

	valueQuery := strings.Join(values, ",")
	statement := `
	INSERT INTO tasks (workflow_id, name, description, parameters)
	VALUES $1
	ON DUPLICATE KEY UPDATE
		name = VALUES(name),
		description = VALUES(description),
		parameters = VALUES(parameters);
	`

	return DbExecAndReturnMany[models.Tasks](
		tx,
		statement, valueQuery,
	)
}

// Delete multiple tasks based on the taskIds
func (t *TasksRepositoryImpl) DeleteTasks(tx *sqlx.Tx, taskIds []uuid.UUID) error {
	stringIds := make([]string, len(taskIds))
	for i, u := range taskIds {
		stringIds[i] = u.String()
	}

	_, err := tx.Exec(`DELETE FROM tasks WHERE id in ()`, strings.Join(stringIds, ","))

	return err
}
