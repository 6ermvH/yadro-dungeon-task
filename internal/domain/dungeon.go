package domain

type DungeonConfig struct {
	Floors   int
	Monsters int
	OpenAt   Clock
	Duration int
}

type Dungeon struct {
	floors        int
	monsters      int
	openAt        Clock
	closeAt       Clock
	monsterFloors int
}

func NewDungeon(config DungeonConfig) Dungeon {
	return Dungeon{
		floors:        config.Floors,
		monsters:      config.Monsters,
		openAt:        config.OpenAt,
		closeAt:       config.OpenAt.AddHours(config.Duration),
		monsterFloors: config.Floors - 1,
	}
}

func (d Dungeon) Floors() int {
	return d.floors
}

func (d Dungeon) MonstersPerFloor() int {
	return d.monsters
}

func (d Dungeon) OpenAt() Clock {
	return d.openAt
}

func (d Dungeon) CloseAt() Clock {
	return d.closeAt
}

func (d Dungeon) MonsterFloors() int {
	if d.monsterFloors < 0 {
		return 0
	}

	return d.monsterFloors
}

func (d Dungeon) BossFloor() int {
	return d.floors
}

func (d Dungeon) IsOpenAt(at Clock) bool {
	return at.AfterOrEqual(d.openAt) && at.Before(d.closeAt)
}
