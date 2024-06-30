package repository

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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

// InsertEdges implements EdgeRepository.
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
		`INSERT INTO tasks (destination_id, source_id)
		VALUES $1
		RETURNING *`, valueQuery,
	)
}

// accepts multiple edge ids to be deleted
func (e *EdgeRepositoryImpl) DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error {
	stringIds := make([]string, len(edgeIds))
	for i, u := range edgeIds {
		stringIds[i] = u.String()
	}

	_, err := tx.Exec(`DELETE FROM tasks WHERE id in ()`, strings.Join(stringIds, ","))

	return err
}
