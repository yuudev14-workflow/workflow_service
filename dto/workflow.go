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
	Name          string                  `json:"name"`
	Description   string                  `json:"description"`
	Parameters    *map[string]interface{} `json:"parameters,omitempty"`
	ConnectorName string                  `json:"connector_name"`
	Operation     string                  `json:"operation"`
	Config        types.Nullable[string]  `json:"config,omitempty"`
	X             int                     `form:"x,default=0"`
	Y             int                     `form:"y,default=0"`
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
}
