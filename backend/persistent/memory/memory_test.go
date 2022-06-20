package memory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smockyio/smocky/backend/mock/config"
	. "github.com/smockyio/smocky/backend/persistent/memory"
)

func TestMemory_GetSetConfig(t *testing.T) {
	cfg := &config.Config{
		Port: "1234",
	}

	m := New()

	err := m.SetConfig(context.Background(), "*id*", cfg)
	require.NoError(t, err)

	value, err := m.GetConfig(context.Background(), "*id*")

	require.NoError(t, err)
	assert.Equal(t, value, cfg)
}

func TestMemory_GetInt(t *testing.T) {
	m := New()

	err := m.Set(context.Background(), "*id*", 200)
	require.NoError(t, err)

	value, err := m.GetInt(context.Background(), "*id*")

	require.NoError(t, err)
	assert.Equal(t, value, 200)
}

func TestMemory_Increase(t *testing.T) {
	m := New()

	err := m.Set(context.Background(), "*id*", 200)
	require.NoError(t, err)

	val, err := m.Increase(context.Background(), "*id*")
	require.NoError(t, err)
	assert.Equal(t, 201, val)

	i, err := m.GetInt(context.Background(), "*id*")
	require.NoError(t, err)
	assert.Equal(t, 201, i)
}

func TestMemory_SetGetActiveSession(t *testing.T) {
	m := New()

	err := m.SetActiveSession(context.Background(), "mockid", "123456")
	require.NoError(t, err)

	v, err := m.GetActiveSession(context.Background(), "mockid")
	require.NoError(t, err)
	assert.Equal(t, "123456", v)
}
