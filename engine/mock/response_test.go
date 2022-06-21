package mock_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/smockyio/smocky/engine/mock"
)

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
