package dto

import (
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type UpdateWorkflowData struct {
	Name        types.Nullable[string] `json:"name,omitempty"`
	Description types.Nullable[string] `json:"description,omitempty"`
	TriggerType types.Nullable[string] `json:"trigger_type,omitempty"`
}

type Task struct {
	Name        string  `db:"name" json:"name"`
	Description string  `db:"description" json:"description"`
	Parameters  *string `db:"parameters" json:"parameters"`
}

type UpdateWorkflowtasks struct {
	Nodes []Task              `json:"nodes"`
	Edges map[string][]string `json:"edges"`
}

type WorkflowPayload struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	TriggerType string  `json:"trigger_type" binding:"required"`
}
