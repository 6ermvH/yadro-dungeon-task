package domain

import "fmt"

type EventType int

const (
	EventRegister EventType = iota + 1
	EventEnterDungeon
	EventKillMonster
	EventNextFloor
	EventPreviousFloor
	EventEnterBossFloor
	EventKillBoss
	EventLeaveDungeon
	EventCannotContinue
	EventRestoreHealth
	EventReceiveDamage
)

type OutputEventType int

const (
	OutputRegistered OutputEventType = iota + 1
	OutputEnteredDungeon
	OutputKilledMonster
	OutputWentNextFloor
	OutputWentPreviousFloor
	OutputEnteredBossFloor
	OutputKilledBoss
	OutputLeftDungeon
	OutputCannotContinue
	OutputRestoredHealth
	OutputReceivedDamage
	OutputDisqualified
	OutputDead
	OutputImpossibleMove
)

type Event struct {
	Time     Clock
	PlayerID int
	Type     EventType
	Param    string
}

type OutputEvent struct {
	Time     Clock
	PlayerID int
	Type     OutputEventType
	Param    string
}

func (t EventType) Int() int {
	return int(t)
}

func (t EventType) String() string {
	return fmt.Sprintf("%d", t)
}

func ParseEventType(id int) (EventType, error) {
	if id < int(EventRegister) || id > int(EventReceiveDamage) {
		return 0, fmt.Errorf("unknown event type %d", id)
	}

	return EventType(id), nil
}
