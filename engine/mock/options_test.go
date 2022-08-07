package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWithIDGeneration(t *testing.T) {
	opts := mockOptions{}
	WithIDGeneration()(&opts)
	assert.True(t, opts.idGeneration)
}
