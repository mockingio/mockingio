package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
