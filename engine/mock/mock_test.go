package mock_test

import (
	_ "embed"
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	. "github.com/smockyio/smocky/engine/mock"
	"github.com/smockyio/smocky/engine/test"
)

func TestConfig(t *testing.T) {
	t.Run("Load config from YAML file", func(t *testing.T) {
		cfg, err := FromYamlFile("fixtures/mock.yml")

		assert.True(t, cfg.Validate() == nil)

		var goldenFile = filepath.Join("fixtures", "mock.golden.yml")
		require.NoError(t, err)

		text, _ := yaml.Marshal(cfg)
		test.UpdateGoldenFile(t, goldenFile, text)

		assert.Equal(t, test.ReadGoldenFile(t, goldenFile), string(text))
	})

	t.Run("Load config from JSON file", func(t *testing.T) {
		cfg, err := FromYamlFile("fixtures/mock.json")

		assert.True(t, cfg.Validate() == nil)

		var goldenFile = filepath.Join("fixtures", "mock.golden.json")
		require.NoError(t, err)

		text, _ := json.MarshalIndent(cfg, "", "  ")

		test.UpdateGoldenFile(t, goldenFile, text)

		assert.Equal(t, test.ReadGoldenFile(t, goldenFile), string(text))
	})

	t.Run("error loading config from YAML file", func(t *testing.T) {
		cfg, err := FromYamlFile("")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})

	t.Run("error loading mock from empty yaml", func(t *testing.T) {
		cfg, err := FromYaml("")
		assert.Error(t, err)
		assert.Nil(t, cfg)
	})
}
