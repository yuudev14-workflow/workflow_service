package service

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/yuudev14-workflow/workflow-service/db/queries"
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
	var values []string

	for _, val := range edges {
		values = append(values,
			fmt.Sprintf(`(%v, %v)`, val.DestinationID, val.SourceID),
		)
	}

	valueQuery := strings.Join(values, ",")

	statement := fmt.Sprintf(queries.INSERT_EDGES, valueQuery)
	logging.Logger.Debugf("insert edge query: %v", statement)

	return repository.DbExecAndReturnMany[models.Edges](
		tx,
		statement,
	)
}

// accepts multiple edge ids to be deleted
func (e *EdgeServiceImpl) DeleteEdges(tx *sqlx.Tx, edgeIds []uuid.UUID) error {
	stringIds := make([]string, len(edgeIds))
	for i, u := range edgeIds {
		stringIds[i] = u.String()
	}

	_, err := tx.Exec(queries.DELETE_EDGES, strings.Join(stringIds, ","))

	return err
}
