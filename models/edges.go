package models

import (
	"github.com/google/uuid"
)

type Edges struct {
	ID            uuid.UUID `db:"id" json:"id"`
	DestinationID string    `db:"destination_id" json:"destination_id"`
	SourceID      string    `db:"source_id" json:"source_id"`
}
