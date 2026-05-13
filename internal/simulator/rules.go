package simulator

import (
	"strconv"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
)

func (s *Simulator) registeredPlayer(event domain.Event) (*domain.Player, bool) {
	player := s.players[event.PlayerID]
	if player == nil || !player.IsRegistered() {
		player = s.player(event.PlayerID)
		s.appendDisqualified(event, player)
		return nil, false
	}

	return player, true
}

func (s *Simulator) activePlayer(event domain.Event) (*domain.Player, bool) {
	player, ok := s.registeredPlayer(event)
	if !ok {
		return nil, false
	}

	if !player.IsInDungeon() || player.IsFinished() {
		s.appendImpossibleMove(event)
		return nil, false
	}

	return player, true
}

func (s *Simulator) canEnterDungeon(player *domain.Player, event domain.Event) bool {
	return player.IsRegistered() &&
		!player.IsInDungeon() &&
		!player.IsFinished() &&
		player.CurrentFloor() == 0 &&
		s.dungeon.IsOpenAt(event.Time)
}

func (s *Simulator) canKillMonster(player *domain.Player) bool {
	if !player.IsInDungeon() || player.IsOnBossFloor() || s.dungeon.MonstersPerFloor() <= 0 {
		return false
	}

	if player.CurrentFloor() > s.dungeon.MonsterFloors() {
		return false
	}

	if player.CurrentFloor() <= player.ClearedFloors() {
		return false
	}

	return player.MonstersKilledOnCurrentFloor() < s.dungeon.MonstersPerFloor()
}

func (s *Simulator) canMoveNextFloor(player *domain.Player) bool {
	if !player.IsInDungeon() || player.IsOnBossFloor() || player.CurrentFloor() >= s.dungeon.BossFloor() {
		return false
	}

	if player.CurrentFloor() <= player.ClearedFloors() {
		return true
	}

	return player.MonstersKilledOnCurrentFloor() >= s.dungeon.MonstersPerFloor()
}

func (s *Simulator) canMovePreviousFloor(player *domain.Player) bool {
	return player.IsInDungeon() && player.CurrentFloor() > 1
}

func (s *Simulator) canEnterBossFloor(player *domain.Player) bool {
	return player.IsInDungeon() &&
		!player.IsOnBossFloor() &&
		player.CurrentFloor() == s.dungeon.BossFloor() &&
		player.ClearedFloors() >= s.dungeon.MonsterFloors()
}

func (s *Simulator) canKillBoss(player *domain.Player) bool {
	return player.IsInDungeon() && player.IsOnBossFloor() && !player.IsBossKilled()
}

func eventParamInt(event domain.Event) (int, bool) {
	value, err := strconv.Atoi(event.Param)
	return value, err == nil
}
