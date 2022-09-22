package mock

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponse_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request Response
		error   bool
	}{
		{"valid status 200", Response{Status: http.StatusOK, RuleAggregation: Or}, false},
		{"default, no status", Response{}, false},
		{"invalid status", Response{Status: 9999}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			assert.Equal(t, tt.error, err != nil, err)
		})
	}
}
