package service

import (
	sq "github.com/Masterminds/squirrel"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/models"
	"github.com/yuudev14-workflow/workflow-service/pkg/logging"
	"github.com/yuudev14-workflow/workflow-service/pkg/repository"
)

type EdgeService interface {
	InsertEdges(tx *sqlx.Tx, edges []models.Edges) ([]models.Edges, error)
	DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error
}

type EdgeServiceImpl struct {
	*sqlx.DB
}

func NewEdgeRepositoryImpl(db *sqlx.DB) EdgeService {
	return &EdgeServiceImpl{
		DB: db,
	}
}

// accepts multiple edge structs to be added in the database in a transaction matter
func (e *EdgeServiceImpl) InsertEdges(tx *sqlx.Tx, edges []models.Edges) ([]models.Edges, error) {
	statement := sq.Insert("tasks").Columns("destination_id", "source_id")

	for _, val := range edges {
		statement = statement.Values(val.DestinationID, val.SourceID)
	}

	sql, args, err := statement.ToSql()

	logging.Logger.Debug("UpsertTasks SQL: ", sql)
	logging.Logger.Debug("UpsertTasks Args: ", args)

	if err != nil {
		logging.Logger.Error("Failed to build SQL query", err)
		return nil, err
	}

	return repository.DbExecAndReturnMany[models.Edges](
		tx,
		sql,
		args...,
	)
}

// accepts multiple edge ids to be deleted
func (e *EdgeServiceImpl) DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error {
	sql, args, err := sq.Delete("edges").Where(sq.Eq{"id": edgeIds}).ToSql()
	logging.Logger.Debug("DeleteEdges SQL: ", sql)
	logging.Logger.Debug("DeleteEdges Args: ", args)
	if err != nil {
		logging.Logger.Error("Failed to build SQL query", err)
		return err
	}
	sql = tx.Rebind(sql)
	_, err = tx.Query(sql, args...)
	logging.Logger.Warn(err)

	return err
}
