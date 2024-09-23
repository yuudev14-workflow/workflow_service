package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Tasks struct {
	ID          uuid.UUID       `db:"id" json:"id"`
	WorkflowID  string          `db:"workflow_id" json:"workflow_id"`
	Status      string          `db:"status" json:"status"`
	Name        string          `db:"name" json:"name"`
	Description string          `db:"description" json:"description"`
	Parameters  json.RawMessage `db:"parameters" json:"parameters"`
	CreatedAt   time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at" json:"updated_at"`
}
