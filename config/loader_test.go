package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/silvan-talos/tlp/config"
)

func TestLoadFromYAML(t *testing.T) {
	t.Run("using example config", loadFromExampleYAML)
	t.Run("config file not found", configFileNotFound)
	t.Run("invalid configuration struct", invalidConfigStruct)
}

func loadFromExampleYAML(t *testing.T) {
	t.Parallel()

	var cfg config.Config
	err := config.LoadFromYAML("config_example.yml", &cfg)
	if err != nil {
		t.Fatalf("load from yaml err: %s", err)
	}
	require.EqualValues(t, config.Config{LogLevel: "debug"}, cfg, "configuration should be correctly loaded")
}

func configFileNotFound(t *testing.T) {
	t.Parallel()

	var cfg config.Config
	err := config.LoadFromYAML("config_not_found.yml", &cfg)
	if err == nil {
		t.Fatal("loading was expected to fail")
	}
	require.ErrorIs(t, err, os.ErrNotExist)
}

func invalidConfigStruct(t *testing.T) {
	t.Parallel()

	cfg := struct {
		Address string `yaml:"address" validate:"required"`
	}{}
	// load valid file
	err := config.LoadFromYAML("config_example.yml", &cfg)
	if err == nil {
		t.Fatal("loading was expected to fail")
	}
	require.ErrorContains(t, err, "invalid configuration:")
}
