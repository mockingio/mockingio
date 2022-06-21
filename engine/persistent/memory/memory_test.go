package memory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smockyio/smocky/engine/mock"
	. "github.com/smockyio/smocky/engine/persistent/memory"
)

func TestMemory_GetSetConfig(t *testing.T) {
	cfg := &mock.Mock{
		Port: "1234",
		ID:   "*id*",
	}

	m := New()

	err := m.SetConfig(context.Background(), cfg)
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

	val, err := m.Increment(context.Background(), "*id*")
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

func TestMemory_GetConfigs(t *testing.T) {
	cfg1 := &mock.Mock{
		Port: "1234",
		ID:   "*id1*",
	}

	cfg2 := &mock.Mock{
		Port: "1234",
		ID:   "*id2*",
	}

	m := New()
	_ = m.SetConfig(context.Background(), cfg1)
	_ = m.SetConfig(context.Background(), cfg2)

	configs, err := m.GetConfigs(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 2, len(configs))
}
