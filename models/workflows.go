package models

import (
	"github.com/google/uuid"
)

type Workflows struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
}

type WorkflowRepository interface {
}
