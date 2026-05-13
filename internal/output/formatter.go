package output

import (
	"fmt"
	"io"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
)

type Formatter struct{}

func (f Formatter) Write(result domain.Result, writer io.Writer) error {
	for _, event := range result.Events {
		if _, err := fmt.Fprintln(writer, f.FormatEvent(event)); err != nil {
			return err
		}
	}

	if _, err := fmt.Fprintln(writer, "Final report:"); err != nil {
		return err
	}

	for _, row := range result.Report {
		if _, err := fmt.Fprintln(writer, f.FormatReportRow(row)); err != nil {
			return err
		}
	}

	return nil
}

func (Formatter) FormatEvent(event domain.OutputEvent) string {
	prefix := fmt.Sprintf("[%s] Player [%d]", event.Time, event.PlayerID)

	switch event.Type {
	case domain.OutputRegistered:
		return prefix + " registered"
	case domain.OutputEnteredDungeon:
		return prefix + " entered the dungeon"
	case domain.OutputKilledMonster:
		return prefix + " killed the monster"
	case domain.OutputWentNextFloor:
		return prefix + " went to the next floor"
	case domain.OutputWentPreviousFloor:
		return prefix + " went to the previous floor"
	case domain.OutputEnteredBossFloor:
		return prefix + " entered the boss's floor"
	case domain.OutputKilledBoss:
		return prefix + " killed the boss"
	case domain.OutputLeftDungeon:
		return prefix + " left the dungeon"
	case domain.OutputCannotContinue:
		return fmt.Sprintf("%s cannot continue due to [%s]", prefix, event.Param)
	case domain.OutputRestoredHealth:
		return fmt.Sprintf("%s has restored [%s] of health", prefix, event.Param)
	case domain.OutputReceivedDamage:
		//nolint:misspell // Required output spelling from the task statement
		return fmt.Sprintf("%s recieved [%s] of damage", prefix, event.Param)
	case domain.OutputDisqualified:
		return prefix + " is disqualified"
	case domain.OutputDead:
		return prefix + " is dead"
	case domain.OutputImpossibleMove:
		return fmt.Sprintf("%s makes imposible move [%s]", prefix, event.Param)
	default:
		return prefix
	}
}

func (Formatter) FormatReportRow(row domain.ReportRow) string {
	return fmt.Sprintf(
		"[%s] %d [%s, %s, %s] HP:%d",
		row.Status,
		row.PlayerID,
		row.TotalTime,
		row.AverageFloorClearTime,
		row.BossKillTime,
		row.HP,
	)
}
