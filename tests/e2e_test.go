package tests

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/6ermvH/yadro-dungeon-task/internal/app"
	"github.com/stretchr/testify/require"
)

func TestE2E(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		config string
		events string
		output string
	}{
		{
			name:   "sample from task",
			config: "sample_config.json",
			events: "sample_events",
			output: "sample_output",
		},
		{
			name:   "cleared floor re-entry blocks monster kills",
			config: "cleared_floor_config.json",
			events: "cleared_floor_events",
			output: "cleared_floor_output",
		},
		{
			name:   "dungeon expiry and death",
			config: "expiry_config.json",
			events: "expiry_events",
			output: "expiry_output",
		},
		{
			name:   "damage and heal edge cases",
			config: "damage_heal_config.json",
			events: "damage_heal_events",
			output: "damage_heal_output",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			configPath := filepath.Join("..", "testdata", tc.config)
			eventsPath := filepath.Join("..", "testdata", tc.events)
			outputPath := filepath.Join("..", "testdata", tc.output)

			expected, err := os.ReadFile(outputPath)
			require.NoError(t, err)

			var stdout bytes.Buffer
			var stderr bytes.Buffer

			code := app.Run([]string{configPath, eventsPath}, &stdout, &stderr)
			require.Zero(t, code, stderr.String())
			require.Equal(t, withTrailingNewline(string(expected)), stdout.String())
		})
	}
}

func withTrailingNewline(value string) string {
	if value == "" || value[len(value)-1] == '\n' {
		return value
	}

	return value + "\n"
}
