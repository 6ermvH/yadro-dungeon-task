package simulator

import (
	"sort"

	"github.com/6ermvH/yadro-dungeon-task/internal/domain"
)

type Simulator struct {
	dungeon domain.Dungeon
	players map[int]*domain.Player
	events  []domain.OutputEvent
}

func New(dungeon domain.Dungeon) *Simulator {
	return &Simulator{
		dungeon: dungeon,
		players: make(map[int]*domain.Player),
	}
}

func (s *Simulator) Run(events []domain.Event) domain.Result {
	for _, event := range events {
		s.Handle(event)
	}

	return domain.Result{
		Events: s.events,
		Report: s.report(),
	}
}

func (s *Simulator) Handle(event domain.Event) {
	player := s.players[event.PlayerID]
	if player != nil && player.IsFinished() {
		return
	}

	if player != nil && player.IsInDungeon() && event.Time.AfterOrEqual(s.dungeon.CloseAt()) {
		player.MarkExpired(s.dungeon.CloseAt())
		return
	}

	switch event.Type {
	case domain.EventRegister:
		s.handleRegister(event)
	case domain.EventEnterDungeon:
		s.handleEnterDungeon(event)
	case domain.EventKillMonster:
		s.handleKillMonster(event)
	case domain.EventNextFloor:
		s.handleNextFloor(event)
	case domain.EventPreviousFloor:
		s.handlePreviousFloor(event)
	case domain.EventEnterBossFloor:
		s.handleEnterBossFloor(event)
	case domain.EventKillBoss:
		s.handleKillBoss(event)
	case domain.EventLeaveDungeon:
		s.handleLeaveDungeon(event)
	case domain.EventCannotContinue:
		s.handleCannotContinue(event)
	case domain.EventRestoreHealth:
		s.handleRestoreHealth(event)
	case domain.EventReceiveDamage:
		s.handleReceiveDamage(event)
	}
}

func (s *Simulator) player(id int) *domain.Player {
	player := s.players[id]
	if player == nil {
		player = domain.NewPlayer(id)
		s.players[id] = player
	}

	return player
}

func (s *Simulator) appendEvent(event domain.Event, eventType domain.OutputEventType, param string) {
	s.events = append(s.events, domain.OutputEvent{
		Time:     event.Time,
		PlayerID: event.PlayerID,
		Type:     eventType,
		Param:    param,
	})
}

func (s *Simulator) appendImpossibleMove(event domain.Event) {
	s.appendEvent(event, domain.OutputImpossibleMove, event.Type.String())
}

func (s *Simulator) appendDisqualified(event domain.Event, player *domain.Player) {
	player.MarkDisqualified(event.Time)
	s.appendEvent(event, domain.OutputDisqualified, "")
}

func (s *Simulator) report() []domain.ReportRow {
	ids := make([]int, 0, len(s.players))
	for id := range s.players {
		ids = append(ids, id)
	}
	sort.Ints(ids)

	rows := make([]domain.ReportRow, 0, len(ids))
	for _, id := range ids {
		player := s.players[id]
		rows = append(rows, s.reportRow(player))
	}

	return rows
}

func (s *Simulator) reportRow(player *domain.Player) domain.ReportRow {
	totalTime := domain.Duration{}
	if player.CurrentFloor() > 0 {
		finishedAt := player.FinishedAt()
		if !player.IsFinished() {
			finishedAt = s.dungeon.CloseAt()
		}

		totalTime = domain.DurationBetween(player.EnteredAt(), finishedAt)
	}

	return domain.ReportRow{
		Status:                player.Status(s.dungeon),
		PlayerID:              player.ID(),
		TotalTime:             totalTime,
		AverageFloorClearTime: player.TotalFloorClearDuration().Div(player.ClearedFloors()),
		BossKillTime:          player.BossKillDuration(),
		HP:                    player.HP(),
	}
}
