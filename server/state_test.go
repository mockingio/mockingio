package server

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestServer_ShutdownServer(t *testing.T) {
	t.Run("shutdown function exist", func(t *testing.T) {
		shutdown := false
		state := &mockServerState{
			shutdownServerFn: func() {
				shutdown = true
			},
			MockID: "*mocid*",
			URL:    "https://mocking.io",
			Status: Running,
		}

		state.shutdownServer()

		assert.Truef(t, shutdown, "shutdown function not called")
		assert.Equal(t, Stopped, state.Status)
		assert.Equal(t, "", state.URL)
	})

	t.Run("shutdown function doesn't exist", func(t *testing.T) {
		state := &mockServerState{
			MockID: "*mocid*",
			URL:    "https://mocking.io",
			Status: Running,
		}

		state.shutdownServer()

		assert.Equal(t, Stopped, state.Status)
		assert.Equal(t, "", state.URL)
	})
}
