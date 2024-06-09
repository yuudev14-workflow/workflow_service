package dto

import (
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type WorkflowPayload struct {
	Name        types.Nullable[string] `json:"name,omitempty"`
	Description types.Nullable[string] `json:"description,omitempty"`
	TriggerType types.Nullable[string] `json:"trigger_type,omitempty"`
}
