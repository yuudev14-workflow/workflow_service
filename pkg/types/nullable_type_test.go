package types_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yuudev14-workflow/workflow-service/pkg/types"
)

type TestStruct struct {
	Name types.Nullable[string] `json:"name"`
}

func TestNullableType(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{
			value:    `{}`,
			expected: false,
		},
		{
			value:    `{"name": null}`,
			expected: true,
		},
	}

	for _, test := range tests {
		var testVal TestStruct
		json.Unmarshal([]byte(test.value), &testVal)
		assert.Equal(t, test.expected, testVal.Name.Set)
	}
}
