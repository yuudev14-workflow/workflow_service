package models

import (
	"time"

	"github.com/google/uuid"
)

type TaskHistory struct {
	ID                uuid.UUID `db:"id" json:"id"`
	WorkflowHistoryID string    `db:"workflow_history_id" json:"workflow_history_id"`
	TaskID            string    `db:"task_id" json:"task_id"`
	Status            string    `db:"status" json:"status"`
	TriggeredAt       time.Time `db:"triggered_at" json:"triggered_at"`
}