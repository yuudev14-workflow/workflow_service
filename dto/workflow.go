package dto

import (
	"database/sql"
	"encoding/json"

	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type UpdateWorkflowPayload struct {
	Name        types.Nullable[string] `json:"name,omitempty"`
	Description types.Nullable[string] `json:"description,omitempty"`
	TriggerType types.Nullable[string] `json:"trigger_type,omitempty"`
}

type WorkflowPayload struct {
	Name        string         `json:"name" binding:"required"`
	Description sql.NullString `json:"description,omitempty"`
	TriggerType string         `json:"trigger_type" binding:"required"`
}

func (ns *WorkflowPayload) UnmarshalJSON(data []byte) error {
	var raw struct {
		Name        string `json:"name"`
		Description string `json:"description,omitempty"`
		TriggerType string `json:"trigger_type"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	ns.Name = raw.Name
	ns.TriggerType = raw.TriggerType
	ns.Description = sql.NullString{String: raw.Description, Valid: raw.Description != ""}

	return nil
}
