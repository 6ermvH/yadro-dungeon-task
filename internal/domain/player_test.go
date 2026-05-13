package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPlayerDefaults(t *testing.T) {
	t.Parallel()

	p := NewPlayer(42)
	require.Equal(t, 42, p.ID())
	require.Equal(t, maxHealth, p.HP())
	require.False(t, p.IsRegistered())
	require.False(t, p.IsInDungeon())
	require.False(t, p.IsFinished())
	require.False(t, p.IsDead())
	require.False(t, p.IsDisqualified())
	require.Equal(t, 0, p.CurrentFloor())
}

func TestPlayerRegister(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)
	p.Register()
	require.True(t, p.IsRegistered())
}

func TestPlayerEnterAndLeaveDungeon(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "14:00:00")
	leaveAt := mustParseClock(t, "15:00:00")

	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(at)

	require.True(t, p.IsInDungeon())
	require.Equal(t, 1, p.CurrentFloor())
	require.Equal(t, at, p.EnteredAt())
	require.True(t, p.HasEnteredDungeon())

	p.LeaveDungeon(leaveAt)
	require.False(t, p.IsInDungeon())
	require.True(t, p.IsFinished())
	require.Equal(t, leaveAt, p.FinishedAt())
}

func TestPlayerHealthRestore(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)
	p.TakeDamage(30)
	require.Equal(t, 70, p.HP())

	p.RestoreHealth(10)
	require.Equal(t, 80, p.HP())

	// Cannot exceed max
	p.RestoreHealth(50)
	require.Equal(t, maxHealth, p.HP())
}

func TestPlayerTakeDamageReturnsDead(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)

	dead := p.TakeDamage(50)
	require.False(t, dead)
	require.Equal(t, 50, p.HP())

	dead = p.TakeDamage(50)
	require.True(t, dead)
	require.Equal(t, 0, p.HP())
}

func TestPlayerTakeDamageOverkill(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)
	dead := p.TakeDamage(200)
	require.True(t, dead)
	require.Equal(t, 0, p.HP())
}

func TestPlayerMarkDead(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "14:00:00")
	p := NewPlayer(1)
	p.MarkDead(at)

	require.True(t, p.IsDead())
	require.True(t, p.IsFinished())
	require.False(t, p.IsInDungeon())
	require.Equal(t, 0, p.HP())
}

func TestPlayerMarkDisqualified(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "14:00:00")
	p := NewPlayer(1)
	p.MarkDisqualified(at)

	require.True(t, p.IsDisqualified())
	require.True(t, p.IsFinished())
	require.False(t, p.IsInDungeon())
	require.Equal(t, at, p.FinishedAt())
}

func TestPlayerMarkExpired(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "16:00:00")
	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(mustParseClock(t, "14:00:00"))
	p.MarkExpired(at)

	require.True(t, p.IsFinished())
	require.False(t, p.IsInDungeon())
	require.Equal(t, at, p.FinishedAt())
}

func TestPlayerFloorNavigation(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "14:00:00")
	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(at)

	require.Equal(t, 1, p.CurrentFloor())

	p.MoveNextFloor(mustParseClock(t, "14:10:00"))
	require.Equal(t, 2, p.CurrentFloor())
	require.Equal(t, 0, p.MonstersKilledOnCurrentFloor())

	p.MovePreviousFloor(mustParseClock(t, "14:15:00"))
	require.Equal(t, 1, p.CurrentFloor())
	require.False(t, p.IsOnBossFloor())
}

func TestPlayerMovePreviousFloorAtFirstFloor(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "14:00:00")
	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(at)

	p.MovePreviousFloor(mustParseClock(t, "14:10:00"))
	require.Equal(t, 1, p.CurrentFloor())
}

func TestPlayerKillMonsterClearsFloor(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "14:00:00")
	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(at)

	p.KillMonster(2, mustParseClock(t, "14:01:00"))
	require.Equal(t, 1, p.MonstersKilledOnCurrentFloor())
	require.Equal(t, 0, p.ClearedFloors())

	p.KillMonster(2, mustParseClock(t, "14:02:00"))
	require.Equal(t, 1, p.ClearedFloors())
	require.False(t, p.TotalFloorClearDuration().IsZero())
}

func TestPlayerKillMonsterZeroMonstersPerFloor(t *testing.T) {
	t.Parallel()

	at := mustParseClock(t, "14:00:00")
	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(at)

	// If monstersPerFloor == 0, floor is cleared immediately
	p.KillMonster(0, mustParseClock(t, "14:01:00"))
	require.Equal(t, 1, p.ClearedFloors())
}

func TestPlayerBossFloor(t *testing.T) {
	t.Parallel()

	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(mustParseClock(t, "14:00:00"))

	bossEnterAt := mustParseClock(t, "14:30:00")
	p.EnterBossFloor(bossEnterAt)
	require.True(t, p.IsOnBossFloor())

	bossKillAt := mustParseClock(t, "14:40:00")
	p.KillBoss(bossKillAt)
	require.True(t, p.IsBossKilled())
	require.Equal(t, "00:10:00", p.BossKillDuration().String())
}

func TestPlayerStatusSuccess(t *testing.T) {
	t.Parallel()

	dungeon := NewDungeon(DungeonConfig{
		Floors: 2, Monsters: 1, OpenAt: mustParseClock(t, "14:00:00"), Duration: 2,
	})

	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(mustParseClock(t, "14:00:00"))
	p.KillMonster(1, mustParseClock(t, "14:01:00"))
	p.MoveNextFloor(mustParseClock(t, "14:02:00"))
	p.EnterBossFloor(mustParseClock(t, "14:02:00"))
	p.KillBoss(mustParseClock(t, "14:05:00"))

	require.True(t, p.ClearedDungeon(dungeon))
	require.Equal(t, PlayerStatusSuccess, p.Status(dungeon))
}

func TestPlayerStatusFail(t *testing.T) {
	t.Parallel()

	dungeon := NewDungeon(DungeonConfig{
		Floors: 2, Monsters: 1, OpenAt: mustParseClock(t, "14:00:00"), Duration: 2,
	})

	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(mustParseClock(t, "14:00:00"))
	p.LeaveDungeon(mustParseClock(t, "14:10:00"))

	require.Equal(t, PlayerStatusFail, p.Status(dungeon))
}

func TestPlayerStatusDisqual(t *testing.T) {
	t.Parallel()

	dungeon := NewDungeon(DungeonConfig{
		Floors: 2, Monsters: 1, OpenAt: mustParseClock(t, "14:00:00"), Duration: 2,
	})

	// Not registered
	p1 := NewPlayer(1)
	require.Equal(t, PlayerStatusDisqual, p1.Status(dungeon))

	// Disqualified
	p2 := NewPlayer(2)
	p2.Register()
	p2.MarkDisqualified(mustParseClock(t, "14:00:00"))
	require.Equal(t, PlayerStatusDisqual, p2.Status(dungeon))
}

func TestPlayerStatusDeadAfterClearing(t *testing.T) {
	t.Parallel()

	dungeon := NewDungeon(DungeonConfig{
		Floors: 2, Monsters: 1, OpenAt: mustParseClock(t, "14:00:00"), Duration: 2,
	})

	p := NewPlayer(1)
	p.Register()
	p.EnterDungeon(mustParseClock(t, "14:00:00"))
	p.KillMonster(1, mustParseClock(t, "14:01:00"))
	p.MoveNextFloor(mustParseClock(t, "14:02:00"))
	p.EnterBossFloor(mustParseClock(t, "14:02:00"))
	p.KillBoss(mustParseClock(t, "14:05:00"))
	p.MarkDead(mustParseClock(t, "14:06:00"))

	require.Equal(t, PlayerStatusFail, p.Status(dungeon))
}

func TestMaxHealth(t *testing.T) {
	t.Parallel()
	require.Equal(t, 100, MaxHealth())
}

func TestHasEnteredDungeon(t *testing.T) {
	t.Parallel()

	// Never entered
	p := NewPlayer(1)
	require.False(t, p.HasEnteredDungeon())

	// Currently in dungeon
	p.EnterDungeon(mustParseClock(t, "14:00:00"))
	require.True(t, p.HasEnteredDungeon())
}
