package models

import (
	"time"

	"github.com/google/uuid"
)

type Workflows struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description string    `db:"description" json:"description"`
	TriggerType string    `json:"trigger_type" db:"trigger_type"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}