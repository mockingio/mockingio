package mock_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	. "github.com/smockyio/smocky/engine/mock"
)

func TestRoute_Validate(t *testing.T) {
	validResponse := []Response{{Status: http.StatusOK}}

	tests := []struct {
		name  string
		route Route
		error bool
	}{
		{"valid route", Route{Request: "POST /", Responses: validResponse}, false},
		{"invalid route, missing request", Route{Responses: validResponse}, true},
		{"invalid route, invalid request", Route{Request: "", Responses: validResponse}, true},
		{"invalid route, missing response", Route{Request: "POST /"}, true},
		{"invalid route, invalid response", Route{Request: "POST /", Responses: []Response{}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.route.Validate()
			assert.Equal(t, tt.error, err != nil)
		})
	}
}
