package simulator

import (
	"testing"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestUnregisteredPlayerDisqualifiedOnce(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:10:00", 3, domain.EventEnterDungeon, ""),
		testEvent(t, "14:12:00", 3, domain.EventKillMonster, ""),
	})

	require.Len(t, result.Events, 1)
	require.Equal(t, domain.OutputDisqualified, result.Events[0].Type)
	require.Equal(t, domain.PlayerStatusDisqual, result.Report[0].Status)
}

func TestHealingCannotExceedMaxHealth(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventReceiveDamage, "20"),
		testEvent(t, "14:12:00", 1, domain.EventRestoreHealth, "90"),
	})

	require.Equal(t, domain.MaxHealth(), result.Report[0].HP)
}

func TestDungeonCloseEndsActiveChallenge(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
	})

	row := result.Report[0]
	require.Equal(t, domain.PlayerStatusFail, row.Status)
	require.Equal(t, "01:55:00", row.TotalTime.String())
}

func TestRegisterAlreadyRegistered(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:01:00", 1, domain.EventRegister, ""),
	})

	require.Len(t, result.Events, 2)
	require.Equal(t, domain.OutputRegistered, result.Events[0].Type)
	require.Equal(t, domain.OutputImpossibleMove, result.Events[1].Type)
}

func TestEnterDungeonBeforeOpen(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:04:00", 1, domain.EventEnterDungeon, ""),
	})

	require.Len(t, result.Events, 2)
	require.Equal(t, domain.OutputImpossibleMove, result.Events[1].Type)
}

func TestEnterDungeonAfterClose(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "16:05:00", 1, domain.EventEnterDungeon, ""),
	})

	require.Equal(t, domain.OutputImpossibleMove, result.Events[1].Type)
}

func TestCannotEnterDungeonTwice(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventEnterDungeon, ""),
	})

	require.Equal(t, domain.OutputImpossibleMove, result.Events[2].Type)
}

func TestKillMonsterOnBossFloorImpossible(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:12:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:13:00", 1, domain.EventNextFloor, ""),
		testEvent(t, "14:13:00", 1, domain.EventEnterBossFloor, ""),
		testEvent(t, "14:14:00", 1, domain.EventKillMonster, ""),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestKillMoreMonstersThanAllowed(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:12:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:13:00", 1, domain.EventKillMonster, ""), // 3rd — impossible, only 2 per floor
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestNextFloorBeforeClearingMonsters(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:12:00", 1, domain.EventNextFloor, ""), // only 1 of 2 killed
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestPreviousFloorAtFirstFloor(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventPreviousFloor, ""),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestEnterBossFloorWithoutClearingAllFloors(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventEnterBossFloor, ""),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestKillBossWithoutBeingOnBossFloor(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventKillBoss, ""),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestPlayerDeath(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventReceiveDamage, "100"),
	})

	require.Equal(t, domain.OutputReceivedDamage, result.Events[2].Type)
	require.Equal(t, domain.OutputDead, result.Events[3].Type)
	require.Equal(t, domain.PlayerStatusFail, result.Report[0].Status)
	require.Equal(t, 0, result.Report[0].HP)
}

func TestEventsAfterDeathIgnored(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventReceiveDamage, "100"),
		testEvent(t, "14:12:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:13:00", 1, domain.EventRestoreHealth, "50"),
	})

	// Only register + enter + damage + dead
	require.Len(t, result.Events, 4)
}

func TestCannotContinueDisqualifies(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:15:00", 1, domain.EventCannotContinue, "out of mana"),
	})

	require.Equal(t, domain.OutputCannotContinue, result.Events[2].Type)
	require.Equal(t, "out of mana", result.Events[2].Param)
	require.Equal(t, domain.PlayerStatusDisqual, result.Report[0].Status)
}

func TestLeaveDungeon(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:20:00", 1, domain.EventLeaveDungeon, ""),
	})

	require.Equal(t, domain.OutputLeftDungeon, result.Events[2].Type)
	require.Equal(t, domain.PlayerStatusFail, result.Report[0].Status)
	require.Equal(t, "00:10:00", result.Report[0].TotalTime.String())
}

func TestSuccessfulDungeonRun(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:12:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:13:00", 1, domain.EventNextFloor, ""),
		testEvent(t, "14:13:00", 1, domain.EventEnterBossFloor, ""),
		testEvent(t, "14:20:00", 1, domain.EventKillBoss, ""),
		testEvent(t, "14:25:00", 1, domain.EventLeaveDungeon, ""),
	})

	row := result.Report[0]
	require.Equal(t, domain.PlayerStatusSuccess, row.Status)
	require.Equal(t, "00:15:00", row.TotalTime.String())
	require.Equal(t, "00:02:00", row.AverageFloorClearTime.String())
	require.Equal(t, "00:07:00", row.BossKillTime.String())
	require.Equal(t, domain.MaxHealth(), row.HP)
}

func TestInvalidDamageParam(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventReceiveDamage, "abc"),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestInvalidHealthParam(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventRestoreHealth, "xyz"),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestNegativeDamageParam(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventReceiveDamage, "-10"),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestNegativeHealthParam(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventRestoreHealth, "-5"),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestDungeonExpiryDuringEvent(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "16:05:00", 1, domain.EventKillMonster, ""), // after close
	})

	// Event after close is silently handled (player marked expired)
	require.Len(t, result.Events, 2) // only register + enter
	require.Equal(t, domain.PlayerStatusFail, result.Report[0].Status)
}

func TestMultiplePlayersReportSortedByID(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 3, domain.EventRegister, ""),
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:00:00", 2, domain.EventRegister, ""),
	})

	require.Len(t, result.Report, 3)
	require.Equal(t, 1, result.Report[0].PlayerID)
	require.Equal(t, 2, result.Report[1].PlayerID)
	require.Equal(t, 3, result.Report[2].PlayerID)
}

func TestPreviousFloorValidFromSecondFloor(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:12:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:13:00", 1, domain.EventNextFloor, ""),
		testEvent(t, "14:14:00", 1, domain.EventPreviousFloor, ""),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputWentPreviousFloor, last.Type)
}

func TestActionOutsideDungeonImpossible(t *testing.T) {
	t.Parallel()

	result := newTestSimulator(t).Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventKillMonster, ""),
	})

	last := result.Events[len(result.Events)-1]
	require.Equal(t, domain.OutputImpossibleMove, last.Type)
}

func TestCannotKillMonstersOnClearedFloor(t *testing.T) {
	t.Parallel()

	// Need 3 floors (2 monster floors + boss) so we can clear floor 1, go to 2, return to 1
	openAt, err := domain.ParseClock("14:05:00")
	require.NoError(t, err)

	sim := New(domain.NewDungeon(domain.DungeonConfig{
		Floors:   3,
		Monsters: 1,
		OpenAt:   openAt,
		Duration: 2,
	}))

	result := sim.Run([]domain.Event{
		testEvent(t, "14:00:00", 1, domain.EventRegister, ""),
		testEvent(t, "14:10:00", 1, domain.EventEnterDungeon, ""),
		testEvent(t, "14:11:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:12:00", 1, domain.EventNextFloor, ""),
		testEvent(t, "14:13:00", 1, domain.EventPreviousFloor, ""),
		testEvent(t, "14:14:00", 1, domain.EventKillMonster, ""),
		testEvent(t, "14:15:00", 1, domain.EventNextFloor, ""),
	})

	// Kill on cleared floor is impossible
	require.Equal(t, domain.OutputImpossibleMove, result.Events[5].Type)
	// But advancing from cleared floor is allowed
	require.Equal(t, domain.OutputWentNextFloor, result.Events[6].Type)
}

func newTestSimulator(t *testing.T) *Simulator {
	t.Helper()

	openAt, err := domain.ParseClock("14:05:00")
	require.NoError(t, err)

	return New(domain.NewDungeon(domain.DungeonConfig{
		Floors:   2,
		Monsters: 2,
		OpenAt:   openAt,
		Duration: 2,
	}))
}

func testEvent(t *testing.T, at string, playerID int, eventType domain.EventType, param string) domain.Event {
	t.Helper()

	eventTime, err := domain.ParseClock(at)
	require.NoError(t, err)

	return domain.Event{
		Time:     eventTime,
		PlayerID: playerID,
		Type:     eventType,
		Param:    param,
	}
}
