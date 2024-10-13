package repository

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
)

type EdgeRepository interface {
	InsertEdges(tx *sqlx.Tx, edges []models.Edges) ([]models.Edges, error)
	DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error
	DeleteAllWorkflowEdges(tx *sqlx.Tx, workflowId string) error
	GetEdgesByWorkflowId(workflowId string) ([]Edges, error)
}

type Edges struct {
	ID                  uuid.UUID `db:"id" json:"id"`
	DestinationID       string    `db:"destination_id" json:"destination_id"`
	SourceID            string    `db:"source_id" json:"source_id"`
	DestinationTaskName string    `db:"destination_task_name" json:"destination_task_name"`
	SourceTaskName      string    `db:"source_task_name" json:"source_task_name"`
}

type EdgeRepositoryImpl struct {
	*sqlx.DB
}

func NewEdgeRepositoryImpl(db *sqlx.DB) EdgeRepository {
	return &EdgeRepositoryImpl{
		DB: db,
	}
}

// GetNodesByWorkflowId implements EdgeRepository.
func (e *EdgeRepositoryImpl) GetEdgesByWorkflowId(workflowId string) ([]Edges, error) {
	sql, args, err := sq.
		Select("e.*, t1.name AS source_task_name, t2.name AS destination_task_name").
		From("edges e").Join("tasks t1 ON e.source_id = t1.id").
		Join("tasks t2 ON e.destination_id = t2.id").
		Where(sq.Eq{"t1.workflow_id": workflowId}).
		ToSql()
	logging.Sugar.Debug("GetNodesByWorkflowId SQL: ", sql)
	logging.Sugar.Debug("GetNodesByWorkflowId Args: ", args)
	if err != nil {
		logging.Sugar.Error("Failed to build SQL query", err)
		return nil, err
	}
	return DbExecAndReturnMany[Edges](
		e.DB,
		sql,
		args...,
	)
}

// accepts multiple edge structs to be added in the database in a transaction matter
// Do nothing if there's already existing source and destination combined
func (e *EdgeRepositoryImpl) InsertEdges(tx *sqlx.Tx, edges []models.Edges) ([]models.Edges, error) {
	statement := sq.Insert("edges").Columns("destination_id", "source_id")

	for _, val := range edges {
		statement = statement.Values(val.DestinationID, val.SourceID)
	}

	sql, args, err := statement.Suffix(`ON CONFLICT (destination_id, source_id) DO NOTHING`).ToSql()

	logging.Sugar.Debug("InsertEdges SQL: ", sql)
	logging.Sugar.Debug("InsertEdges Args: ", args)

	if err != nil {
		logging.Sugar.Error("Failed to build SQL query", err)
		return nil, err
	}

	return DbExecAndReturnMany[models.Edges](
		tx,
		sql,
		args...,
	)
}

// accepts multiple edge ids to be deleted
func (e *EdgeRepositoryImpl) DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error {
	sql, args, err := sq.Delete("edges").Where(sq.Eq{"id": edgeIds}).ToSql()
	logging.Sugar.Debug("DeleteEdges SQL: ", sql)
	logging.Sugar.Debug("DeleteEdges Args: ", args)
	if err != nil {
		logging.Sugar.Error("Failed to build SQL query", err)
		return err
	}
	sql = tx.Rebind(sql)
	_, err = tx.Exec(sql, args...)
	if err != nil {
		logging.Sugar.Warn(err)
	}

	return err
}

// accepts multiple edge ids to be deleted
func (e *EdgeRepositoryImpl) DeleteAllWorkflowEdges(tx *sqlx.Tx, workflowId string) error {

	// Main delete query with the subquery in both conditions (destination_id and source_id)
	deleteQuery := sq.Delete("edges").Where(`
		destination_id IN  (SELECT id FROM tasks WHERE workflow_id = ?) OR 
		source_id IN (SELECT id FROM tasks WHERE workflow_id = ?)`, workflowId, workflowId)
	// Convert the query to SQL
	sql, args, err := deleteQuery.ToSql()
	logging.Sugar.Debug("DeleteAllWorkflowEdges SQL: ", sql)
	logging.Sugar.Debug("DeleteAllWorkflowEdges Args: ", args)
	if err != nil {
		logging.Sugar.Error("Failed to build SQL query", err)
		return err
	}
	sql = tx.Rebind(sql)
	_, err = tx.Exec(sql, args...)
	if err != nil {
		logging.Sugar.Warn(err)
	}

	return err
}
