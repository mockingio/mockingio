package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRule_Clone(t *testing.T) {
	rule := Rule{
		ID:       newID(),
		Target:   "header",
		Modifier: "Authorization",
		Value:    "Bearer...",
		Operator: "equal",
	}

	clone := rule.Clone()
	assert.True(t, clone.Validate() == nil)
	assert.NotEqual(t, rule.ID, clone.ID)

	// TODO: Should compare pointer for all property
}

func TestRule_Validate(t *testing.T) {
	tests := []struct {
		name  string
		rule  Rule
		error bool
	}{
		{"valid rule", Rule{Target: "header", Modifier: "Authorization", Value: "Bearer...", Operator: "equal"}, false},
		{"invalid route, missing target", Rule{Target: "", Modifier: "Authorization", Value: "Bearer...", Operator: "equal"}, true},
		{"invalid route, missing value", Rule{Target: "cookie", Modifier: "Authorization", Value: "", Operator: "equal"}, true},
		{"invalid route, missing operator", Rule{Target: "body", Modifier: "Authorization", Value: "Bearer...", Operator: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			assert.Equal(t, tt.error, err != nil)
		})
	}
}
