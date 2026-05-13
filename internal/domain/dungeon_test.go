package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewDungeon(t *testing.T) {
	t.Parallel()

	openAt := mustParseClock(t, "14:00:00")
	d := NewDungeon(DungeonConfig{
		Floors: 3, Monsters: 5, OpenAt: openAt, Duration: 2,
	})

	require.Equal(t, 3, d.Floors())
	require.Equal(t, 5, d.MonstersPerFloor())
	require.Equal(t, openAt, d.OpenAt())
	require.Equal(t, "16:00:00", d.CloseAt().String())
	require.Equal(t, 2, d.MonsterFloors())
	require.Equal(t, 3, d.BossFloor())
}

func TestDungeonMonsterFloorsMinZero(t *testing.T) {
	t.Parallel()

	d := NewDungeon(DungeonConfig{
		Floors: 0, Monsters: 0, OpenAt: mustParseClock(t, "14:00:00"), Duration: 1,
	})

	require.Equal(t, 0, d.MonsterFloors())
}

func TestDungeonIsOpenAt(t *testing.T) {
	t.Parallel()

	openAt := mustParseClock(t, "14:00:00")
	d := NewDungeon(DungeonConfig{
		Floors: 2, Monsters: 1, OpenAt: openAt, Duration: 2,
	})

	require.False(t, d.IsOpenAt(mustParseClock(t, "13:59:59")))
	require.True(t, d.IsOpenAt(mustParseClock(t, "14:00:00")))
	require.True(t, d.IsOpenAt(mustParseClock(t, "15:00:00")))
	require.True(t, d.IsOpenAt(mustParseClock(t, "15:59:59")))
	require.False(t, d.IsOpenAt(mustParseClock(t, "16:00:00")))
}

func TestParseEventType(t *testing.T) {
	t.Parallel()

	et, err := ParseEventType(1)
	require.NoError(t, err)
	require.Equal(t, EventRegister, et)

	et, err = ParseEventType(11)
	require.NoError(t, err)
	require.Equal(t, EventReceiveDamage, et)

	_, err = ParseEventType(0)
	require.Error(t, err)

	_, err = ParseEventType(12)
	require.Error(t, err)

	_, err = ParseEventType(-1)
	require.Error(t, err)
}

func TestEventTypeString(t *testing.T) {
	t.Parallel()

	require.Equal(t, "1", EventRegister.String())
	require.Equal(t, "5", EventPreviousFloor.String())
	require.Equal(t, 5, EventPreviousFloor.Int())
}
