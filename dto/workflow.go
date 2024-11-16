package dto

import (
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type WorkflowFilter struct {
	Name *string `form:"name"`
}

type UpdateWorkflowData struct {
	Name        types.Nullable[string] `json:"name,omitempty"`
	Description types.Nullable[string] `json:"description,omitempty"`
	TriggerType types.Nullable[string] `json:"trigger_type,omitempty"`
}

type UpdateWorkflowHistoryData struct {
	Status types.Nullable[string] `json:"status,omitempty"`
	Error  types.Nullable[string] `json:"error,omitempty"`
	Result interface{}            `json:"result,omitempty"`
}

type UpdateTaskHistoryData struct {
	Status types.Nullable[string] `json:"status,omitempty"`
	Error  types.Nullable[string] `json:"error,omitempty"`
	Result interface{}            `json:"result,omitempty"`
}

type Task struct {
	Name          string                  `db:"name" json:"name"`
	Description   string                  `db:"description" json:"description"`
	Parameters    *map[string]interface{} `db:"parameters" json:"parameters,omitempty"`
	ConnectorName string                  `db:"connector_name" json:"connector_name"`
	Operation     string                  `db:"operation" json:"operation"`
	Config        types.Nullable[string]  `json:"config,omitempty"`
}

type UpdateWorkflowtasks struct {
	Nodes []Task              `json:"nodes"`
	Edges map[string][]string `json:"edges"`
}

type UpdateWorkflowTaskHistoryStatus struct {
	Status string `json:"status" binding:"required"`
}

type WorkflowPayload struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
	TriggerType string  `json:"trigger_type" binding:"required"`
}
