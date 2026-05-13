package simulator

import (
	"strconv"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
)

func (s *Simulator) handleRegister(event domain.Event) {
	player := s.player(event.PlayerID)
	if player.IsRegistered() || player.IsFinished() {
		s.appendImpossibleMove(event)
		return
	}

	player.Register()
	s.appendEvent(event, domain.OutputRegistered, "")
}

func (s *Simulator) handleEnterDungeon(event domain.Event) {
	player, ok := s.registeredPlayer(event)
	if !ok {
		return
	}

	if !s.canEnterDungeon(player, event) {
		s.appendImpossibleMove(event)
		return
	}

	player.EnterDungeon(event.Time)
	s.appendEvent(event, domain.OutputEnteredDungeon, "")
}

func (s *Simulator) handleKillMonster(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	if !s.canKillMonster(player) {
		s.appendImpossibleMove(event)
		return
	}

	player.KillMonster(s.dungeon.MonstersPerFloor(), event.Time)
	s.appendEvent(event, domain.OutputKilledMonster, "")
}

func (s *Simulator) handleNextFloor(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	if !s.canMoveNextFloor(player) {
		s.appendImpossibleMove(event)
		return
	}

	player.MoveNextFloor(event.Time)
	s.appendEvent(event, domain.OutputWentNextFloor, "")
}

func (s *Simulator) handlePreviousFloor(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	if !s.canMovePreviousFloor(player) {
		s.appendImpossibleMove(event)
		return
	}

	player.MovePreviousFloor(event.Time)
	s.appendEvent(event, domain.OutputWentPreviousFloor, "")
}

func (s *Simulator) handleEnterBossFloor(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	if !s.canEnterBossFloor(player) {
		s.appendImpossibleMove(event)
		return
	}

	player.EnterBossFloor(event.Time)
	s.appendEvent(event, domain.OutputEnteredBossFloor, "")
}

func (s *Simulator) handleKillBoss(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	if !s.canKillBoss(player) {
		s.appendImpossibleMove(event)
		return
	}

	player.KillBoss(event.Time)
	s.appendEvent(event, domain.OutputKilledBoss, "")
}

func (s *Simulator) handleLeaveDungeon(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	player.LeaveDungeon(event.Time)
	s.appendEvent(event, domain.OutputLeftDungeon, "")
}

func (s *Simulator) handleCannotContinue(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	player.MarkDisqualified(event.Time)
	s.appendEvent(event, domain.OutputCannotContinue, event.Param)
}

func (s *Simulator) handleRestoreHealth(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	health, ok := eventParamInt(event)
	if !ok || health < 0 {
		s.appendImpossibleMove(event)
		return
	}

	player.RestoreHealth(health)
	s.appendEvent(event, domain.OutputRestoredHealth, strconv.Itoa(health))
}

func (s *Simulator) handleReceiveDamage(event domain.Event) {
	player, ok := s.activePlayer(event)
	if !ok {
		return
	}

	damage, ok := eventParamInt(event)
	if !ok || damage < 0 {
		s.appendImpossibleMove(event)
		return
	}

	dead := player.TakeDamage(damage)
	s.appendEvent(event, domain.OutputReceivedDamage, strconv.Itoa(damage))
	if dead {
		player.MarkDead(event.Time)
		s.appendEvent(event, domain.OutputDead, "")
	}
}
