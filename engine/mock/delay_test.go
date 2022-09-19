package mock_test

import (
	"fmt"
	"testing"

	. "github.com/mockingio/mockingio/engine/mock"
	"github.com/stretchr/testify/assert"
)

func TestDelay_Validate(t *testing.T) {
	tests := []struct {
		delay   Delay
		isValid bool
	}{
		{Delay{Min: 0, Max: 0}, true},
		{Delay{Min: 0, Max: 1}, true},
		{Delay{Min: 1, Max: 1}, true},
		{Delay{Min: 1, Max: 0}, false},
		{Delay{Min: 10, Max: 61}, false},
		{Delay{Min: -1, Max: 0}, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Min: %v, Max: %v, IsValid: %v", tt.delay.Min, tt.delay.Max, tt.isValid), func(t *testing.T) {
			err := tt.delay.Validate()
			if tt.isValid {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}
