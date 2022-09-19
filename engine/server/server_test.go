package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/persistent"
	"github.com/mockingio/mockingio/engine/persistent/memory"
)

var (
	//go:embed certs/rootCA.pem
	defaultRootCAPEM []byte
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

func TestServer_TLS(t *testing.T) {
	certPath, _ := filepath.Abs("./fixtures/certs/cert.pem")
	keyPath, _ := filepath.Abs("./fixtures/certs/key.pem")

	var shutdownFns []func()
	defer func() {
		for _, fn := range shutdownFns {
			fn()
		}
	}()

	tests := []struct {
		name string
		mock *mock.Mock
	}{
		{
			name: "TLS with default config",
			mock: &mock.Mock{
				ID: "*mock-id-1*",
				TLS: &mock.TLS{
					Enabled: true,
				},
			},
		},
		{
			name: "TLS with default config since custom config missing key",
			mock: &mock.Mock{
				ID: "*mock-id-1*",
				TLS: &mock.TLS{
					Enabled:     true,
					PEMCertPath: "*",
				},
			},
		},
		{
			name: "TLS with default config since custom config missing cert",
			mock: &mock.Mock{
				ID: "*mock-id-1*",
				TLS: &mock.TLS{
					Enabled:    true,
					PEMKeyPath: "*",
				},
			},
		},
		{
			name: "TLS with custom config",
			mock: &mock.Mock{
				ID: "*mock-id-1*",
				TLS: &mock.TLS{
					Enabled:     true,
					PEMCertPath: certPath,
					PEMKeyPath:  keyPath,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := memory.New()
			_ = db.SetMock(context.Background(), tt.mock)

			server := New(db)
			state, err := server.NewMockServerByID(context.Background(), "*mock-id-1*")
			if state != nil {
				shutdownFns = append(shutdownFns, state.shutdownServer)
			}

			assert.NoError(t, err)

			cert, err := getTLSCert(tt.mock)
			require.NoError(t, err)

			roots := x509.NewCertPool()
			roots.AppendCertsFromPEM(defaultRootCAPEM)
			_, err = tls.Dial("tcp", strings.ReplaceAll(state.URL, "https://", ""), &tls.Config{
				Certificates: []tls.Certificate{*cert},
				RootCAs:      roots,
			})
			assert.NoError(t, err)
		})
	}

}

func setupDatabase() persistent.EngineDB {
	db := memory.New()

	_ = db.SetMock(context.Background(), &mock.Mock{ID: "*mock-id-1*"})
	_ = db.SetMock(context.Background(), &mock.Mock{ID: "*mock-id-2*"})

	return db
}
