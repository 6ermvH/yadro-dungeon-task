package parser

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestParseLineBasic(t *testing.T) {
	t.Parallel()

	event, err := EventParser{}.ParseLine("[14:00:00] 1 1")
	require.NoError(t, err)
	require.Equal(t, 1, event.PlayerID)
	require.Equal(t, domain.EventRegister, event.Type)
	require.Equal(t, "", event.Param)
}

func TestParseLineWithMultiWordParam(t *testing.T) {
	t.Parallel()

	event, err := EventParser{}.ParseLine("[15:00:00] 4 9 out of mana")
	require.NoError(t, err)
	require.Equal(t, 4, event.PlayerID)
	require.Equal(t, domain.EventCannotContinue, event.Type)
	require.Equal(t, "out of mana", event.Param)
}

func TestParseLineWithNumericParam(t *testing.T) {
	t.Parallel()

	event, err := EventParser{}.ParseLine("[14:27:00] 2 11 60")
	require.NoError(t, err)
	require.Equal(t, 2, event.PlayerID)
	require.Equal(t, domain.EventReceiveDamage, event.Type)
	require.Equal(t, "60", event.Param)
}

func TestParseLineRejectsInvalidEventType(t *testing.T) {
	t.Parallel()

	_, err := EventParser{}.ParseLine("[15:00:00] 4 99")
	require.Error(t, err)
}

func TestParseLineRejectsTooFewFields(t *testing.T) {
	t.Parallel()

	_, err := EventParser{}.ParseLine("[15:00:00] 4")
	require.Error(t, err)
	require.Contains(t, err.Error(), "at least 3 fields")
}

func TestParseLineRejectsInvalidTimeFormat(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
	}{
		{"no brackets", "15:00:00 4 1"},
		{"missing close bracket", "[15:00:00 4 1"},
		{"wrong length", "[5:00:00] 4 1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := EventParser{}.ParseLine(tc.input)
			require.Error(t, err)
		})
	}
}

func TestParseLineRejectsInvalidPlayerID(t *testing.T) {
	t.Parallel()

	_, err := EventParser{}.ParseLine("[15:00:00] abc 1")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid player id")
}

func TestParseLineRejectsInvalidEventID(t *testing.T) {
	t.Parallel()

	_, err := EventParser{}.ParseLine("[15:00:00] 1 abc")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid event id")
}

func TestParseFileSuccess(t *testing.T) {
	t.Parallel()

	content := "[14:00:00] 1 1\n[14:10:00] 1 2\n"
	path := filepath.Join(t.TempDir(), "events")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	events, err := EventParser{}.ParseFile(path)
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.Equal(t, domain.EventRegister, events[0].Type)
	require.Equal(t, domain.EventEnterDungeon, events[1].Type)
}

func TestParseFileSkipsEmptyLines(t *testing.T) {
	t.Parallel()

	content := "[14:00:00] 1 1\n\n\n[14:10:00] 1 2\n"
	path := filepath.Join(t.TempDir(), "events")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	events, err := EventParser{}.ParseFile(path)
	require.NoError(t, err)
	require.Len(t, events, 2)
}

func TestParseFileErrorOnInvalidLine(t *testing.T) {
	t.Parallel()

	content := "[14:00:00] 1 1\nbad line\n"
	path := filepath.Join(t.TempDir(), "events")
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	_, err := EventParser{}.ParseFile(path)
	require.Error(t, err)
	require.Contains(t, err.Error(), "line 2")
}

func TestParseFileErrorOnMissingFile(t *testing.T) {
	t.Parallel()

	_, err := EventParser{}.ParseFile("/nonexistent/path")
	require.Error(t, err)
	require.Contains(t, err.Error(), "open events file")
}
