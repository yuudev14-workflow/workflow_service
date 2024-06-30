package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
	"github.com/yuudev14-workflow/workflow-service/models"
)

type EdgeRepository interface {
	InsertEdges(tx *sqlx.Tx, edges []models.Edges) ([]models.Edges, error)
	DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error
}

type EdgeRepositoryImpl struct {
	*sqlx.DB
}

func NewEdgeRepositoryImpl(db *sqlx.DB) EdgeRepository {
	return &EdgeRepositoryImpl{
		DB: db,
	}
}

// accepts multiple edge structs to be added in the database in a transaction matter
func (e *EdgeRepositoryImpl) InsertEdges(tx *sqlx.Tx, edges []models.Edges) ([]models.Edges, error) {
	var values []string

	for _, val := range edges {
		values = append(values,
			fmt.Sprintf(`(%v, %v)`, val.DestinationID, val.SourceID),
		)
	}

	valueQuery := strings.Join(values, ",")

	return DbExecAndReturnMany[models.Edges](
		tx,
		queries.INSERT_EDGES, valueQuery,
	)
}

// accepts multiple edge ids to be deleted
func (e *EdgeRepositoryImpl) DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error {
	stringIds := make([]string, len(edgeIds))
	for i, u := range edgeIds {
		stringIds[i] = u.String()
	}

	_, err := tx.Exec(queries.DELETE_EDGES, strings.Join(stringIds, ","))

	return err
}
