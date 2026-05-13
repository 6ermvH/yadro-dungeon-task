package output

import (
	"bytes"
	"testing"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
	"github.com/stretchr/testify/require"
)

func mustParseClock(t *testing.T, value string) domain.Clock {
	t.Helper()

	clock, err := domain.ParseClock(value)
	require.NoError(t, err)

	return clock
}

func TestFormatEventAllTypes(t *testing.T) {
	t.Parallel()

	f := Formatter{}
	time := mustParseClock(t, "14:00:00")

	cases := []struct {
		eventType domain.OutputEventType
		param     string
		expected  string
	}{
		{domain.OutputRegistered, "", "[14:00:00] Player [1] registered"},
		{domain.OutputEnteredDungeon, "", "[14:00:00] Player [1] entered the dungeon"},
		{domain.OutputKilledMonster, "", "[14:00:00] Player [1] killed the monster"},
		{domain.OutputWentNextFloor, "", "[14:00:00] Player [1] went to the next floor"},
		{domain.OutputWentPreviousFloor, "", "[14:00:00] Player [1] went to the previous floor"},
		{domain.OutputEnteredBossFloor, "", "[14:00:00] Player [1] entered the boss's floor"},
		{domain.OutputKilledBoss, "", "[14:00:00] Player [1] killed the boss"},
		{domain.OutputLeftDungeon, "", "[14:00:00] Player [1] left the dungeon"},
		{domain.OutputCannotContinue, "out of mana", "[14:00:00] Player [1] cannot continue due to [out of mana]"},
		{domain.OutputRestoredHealth, "50", "[14:00:00] Player [1] has restored [50] of health"},
		//nolint:misspell // Required output spelling from the task statement
		{domain.OutputReceivedDamage, "30", "[14:00:00] Player [1] recieved [30] of damage"},
		{domain.OutputDisqualified, "", "[14:00:00] Player [1] is disqualified"},
		{domain.OutputDead, "", "[14:00:00] Player [1] is dead"},
		{domain.OutputImpossibleMove, "5", "[14:00:00] Player [1] makes imposible move [5]"},
	}

	for _, tc := range cases {
		event := domain.OutputEvent{
			Time:     time,
			PlayerID: 1,
			Type:     tc.eventType,
			Param:    tc.param,
		}
		require.Equal(t, tc.expected, f.FormatEvent(event))
	}
}

func TestFormatReportRow(t *testing.T) {
	t.Parallel()

	f := Formatter{}
	row := domain.ReportRow{
		Status:                domain.PlayerStatusSuccess,
		PlayerID:              1,
		TotalTime:             domain.NewDuration(1440),
		AverageFloorClearTime: domain.NewDuration(300),
		BossKillTime:          domain.NewDuration(660),
		HP:                    35,
	}

	expected := "[SUCCESS] 1 [00:24:00, 00:05:00, 00:11:00] HP:35"
	require.Equal(t, expected, f.FormatReportRow(row))
}

func TestFormatReportRowFail(t *testing.T) {
	t.Parallel()

	f := Formatter{}
	row := domain.ReportRow{
		Status:   domain.PlayerStatusFail,
		PlayerID: 2,
		HP:       0,
	}

	expected := "[FAIL] 2 [00:00:00, 00:00:00, 00:00:00] HP:0"
	require.Equal(t, expected, f.FormatReportRow(row))
}

func TestWriteFullResult(t *testing.T) {
	t.Parallel()

	f := Formatter{}
	time := mustParseClock(t, "14:00:00")

	result := domain.Result{
		Events: []domain.OutputEvent{
			{Time: time, PlayerID: 1, Type: domain.OutputRegistered},
		},
		Report: []domain.ReportRow{
			{
				Status:   domain.PlayerStatusDisqual,
				PlayerID: 1,
				HP:       100,
			},
		},
	}

	var buf bytes.Buffer
	err := f.Write(result, &buf)
	require.NoError(t, err)

	output := buf.String()
	require.Contains(t, output, "[14:00:00] Player [1] registered")
	require.Contains(t, output, "Final report:")
	require.Contains(t, output, "[DISQUAL] 1")
}

func TestFormatEventUnknownType(t *testing.T) {
	t.Parallel()

	f := Formatter{}
	event := domain.OutputEvent{
		Time:     mustParseClock(t, "14:00:00"),
		PlayerID: 1,
		Type:     domain.OutputEventType(999),
	}

	result := f.FormatEvent(event)
	require.Equal(t, "[14:00:00] Player [1]", result)
}
