package memory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	
	"github.com/smockyio/smocky/backend/mock/config"
	. "github.com/smockyio/smocky/backend/persistent/memory"
)

func TestMemory_Get(t *testing.T) {
	cfg := &config.Config{
		Port: "1234",
	}

	m := New()

	err := m.Set("*id*", cfg)
	require.NoError(t, err)

	found, err := m.Get("*id*")

	require.NoError(t, err)
	assert.Equal(t, found, cfg)
}
