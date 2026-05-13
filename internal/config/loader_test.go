package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoaderSuccess(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
		"Floors": 3,
		"Monsters": 5,
		"OpenAt": "14:00:00",
		"Duration": 2
	}`), 0o600))

	cfg, err := Loader{}.Load(path)
	require.NoError(t, err)
	require.Equal(t, 3, cfg.Floors)
	require.Equal(t, 5, cfg.Monsters)
	require.Equal(t, 2, cfg.Duration)
	require.Equal(t, "14:00:00", cfg.OpenAt.String())
}

func TestLoaderRejectsInvalidConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
		"Floors": 1,
		"Monsters": -1,
		"OpenAt": "bad-time",
		"Duration": 0
	}`), 0o600))

	_, err := Loader{}.Load(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "validate config")
}

func TestLoaderRejectsInvalidJSON(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(`not json`), 0o600))

	_, err := Loader{}.Load(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "parse config json")
}

func TestLoaderRejectsMissingFile(t *testing.T) {
	t.Parallel()

	_, err := Loader{}.Load("/nonexistent/config.json")
	require.Error(t, err)
	require.Contains(t, err.Error(), "read config")
}

func TestLoaderRejectsInvalidOpenAt(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
		"Floors": 2,
		"Monsters": 1,
		"OpenAt": "99:99:99",
		"Duration": 1
	}`), 0o600))

	_, err := Loader{}.Load(path)
	require.Error(t, err)
}

func TestLoaderMinimalValidConfig(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(`{
		"Floors": 2,
		"Monsters": 0,
		"OpenAt": "00:00:00",
		"Duration": 1
	}`), 0o600))

	cfg, err := Loader{}.Load(path)
	require.NoError(t, err)
	require.Equal(t, 2, cfg.Floors)
	require.Equal(t, 0, cfg.Monsters)
}
