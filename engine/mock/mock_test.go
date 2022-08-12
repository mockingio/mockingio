package mock_test

import (
	_ "embed"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"

	. "github.com/mockingio/mockingio/engine/mock"
	"github.com/mockingio/mockingio/engine/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig(t *testing.T) {
	t.Run("Load mock from YAML file", func(t *testing.T) {
		cfg, err := FromFile("fixtures/mock.yml")

		assert.True(t, cfg.Validate() == nil)

		var goldenFile = filepath.Join("fixtures", "mock.golden.yml")
		require.NoError(t, err)

		text, _ := yaml.Marshal(cfg)
		test.UpdateGoldenFile(t, goldenFile, text)

		assert.Equal(t, "fixtures/mock.yml", cfg.FilePath)
		assert.Equal(t, test.ReadGoldenFile(t, goldenFile), string(text))
	})

	t.Run("Load mock from YAML file, with ID generation option", func(t *testing.T) {
		mock, err := FromFile("fixtures/mock.yml", WithIDGeneration())
		require.NoError(t, err)

		assert.True(t, mock.ID != "")
		assert.True(t, mock.Routes[0].ID != "")
		assert.True(t, mock.Routes[0].Responses[0].ID != "")
		assert.True(t, mock.Routes[0].Responses[0].Rules[0].ID != "")
	})

	t.Run("When method, status is not presented, use default GET/200 as response", func(t *testing.T) {
		mock, err := FromFile("fixtures/mock_no_method_status.yml")
		require.NoError(t, err)

		assert.Equal(t, "GET", mock.Routes[0].Method)
		assert.Equal(t, 200, mock.Routes[0].Responses[0].Status)
	})

	t.Run("error loading config from YAML file", func(t *testing.T) {
		mock, err := FromFile("")
		assert.Error(t, err)
		assert.Nil(t, mock)
	})

	t.Run("error loading mock from empty yaml", func(t *testing.T) {
		mock, err := FromYaml("")
		assert.Error(t, err)
		assert.Nil(t, mock)
	})

	t.Run("proxy is enabled", func(t *testing.T) {
		mock := &Mock{
			Proxy: &Proxy{
				Enabled: true,
			},
		}
		assert.False(t, New().ProxyEnabled())
		assert.True(t, mock.ProxyEnabled())
	})

	t.Run("TLS is enabled", func(t *testing.T) {
		mock := &Mock{
			TLS: &TLS{
				Enabled: true,
			},
		}
		assert.False(t, New().TLSEnabled())
		assert.True(t, mock.TLSEnabled())
	})
}
