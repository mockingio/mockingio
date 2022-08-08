package mock

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse_Clone(t *testing.T) {
	rule := Response{
		ID:              "",
		Status:          http.StatusOK,
		RuleAggregation: Or,
	}

	clone := rule.Clone()
	assert.True(t, clone.Validate() == nil)
	assert.NotEqual(t, rule.ID, clone.ID)

	// TODO: Should compare pointer for all property
}

func TestResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request Response
		error   bool
	}{
		{"valid status 200", Response{Status: http.StatusOK, RuleAggregation: Or}, false},
		{"invalid status 1000", Response{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			assert.Equal(t, tt.error, err != nil)
		})
	}
}
