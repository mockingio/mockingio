package server

import (
	"context"
	"testing"

	"github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
	"github.com/mockingio/mockingio/engine/persistent/memory"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	// no crash test
	_ = memory.New()
}

func TestServer_NewMockServerByID(t *testing.T) {
	t.Run("mock not found", func(t *testing.T) {
		server := New(setupDatabase())

		_, err := server.NewMockServerByID(context.Background(), "*random-id*")
		assert.Error(t, err)
	})

	t.Run("mock exists, server started successfully", func(t *testing.T) {
		server := New(setupDatabase())
		state, err := server.NewMockServerByID(context.Background(), "*mock-id-1*")
		if state != nil {
			defer state.shutdownServer()
		}

		assert.Equal(t, "running", state.Status)
		assert.True(t, state.URL != "")
		assert.NoError(t, err)
	})
}

func TestServer_StopMockServer(t *testing.T) {
	server := New(setupDatabase())
	state, _ := server.NewMockServerByID(context.Background(), "*mock-id-1*")

	assert.Equal(t, "running", state.Status)
	newState, _ := server.StopMockServer(state.MockID)
	assert.Equal(t, "stopped", newState.Status)
	assert.Equal(t, "", newState.URL)

	_, err := server.StopMockServer("random id")
	assert.Error(t, err)
}

func TestServer_StopAllServers(t *testing.T) {
	server := New(setupDatabase())
	_, _ = server.NewMockServerByID(context.Background(), "*mock-id-1*")
	_, _ = server.NewMockServerByID(context.Background(), "*mock-id-2*")
	states := server.GetMockServerStates()

	assert.Equal(t, 2, len(states))
	assert.Equal(t, "running", states["*mock-id-1*"].Status)
	assert.Equal(t, "running", states["*mock-id-2*"].Status)

	server.StopAllServers()

	assert.Equal(t, "stopped", states["*mock-id-1*"].Status)
	assert.Equal(t, "stopped", states["*mock-id-2*"].Status)
}

func TestServer_GetMockServerURLs(t *testing.T) {
	server := New(setupDatabase())
	state1, _ := server.NewMockServerByID(context.Background(), "*mock-id-1*")
	defer server.StopAllServers()

	urls := server.GetMockServerURLs()
	assert.Equal(t, []string{state1.URL}, urls)
}

func setupDatabase() persistent.Persistent {
	db := memory.New()

	_ = db.SetMock(context.Background(), &mock.Mock{ID: "*mock-id-1*"})
	_ = db.SetMock(context.Background(), &mock.Mock{ID: "*mock-id-2*"})

	return db
}
