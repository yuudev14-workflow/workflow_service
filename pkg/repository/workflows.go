package repository

import "github.com/jmoiron/sqlx"

type WorkflowRepository interface {
}

type WorkflowRepositoryImpl struct {
	*sqlx.DB
}

func NewWorkflowRepositoryImple(db *sqlx.DB) WorkflowRepository {
	return &WorkflowRepositoryImpl{
		DB: db,
	}
}
