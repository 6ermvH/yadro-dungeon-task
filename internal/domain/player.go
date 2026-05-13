package domain

const maxHealth = 100

type PlayerStatus string

const (
	PlayerStatusSuccess PlayerStatus = "SUCCESS"
	PlayerStatusFail    PlayerStatus = "FAIL"
	PlayerStatusDisqual PlayerStatus = "DISQUAL"
)

type Player struct {
	id int

	hp int

	registered bool
	inDungeon  bool
	finished   bool

	disqualified bool
	dead         bool

	currentFloor int
	onBossFloor  bool
	bossKilled   bool

	enteredAt  Clock
	finishedAt Clock

	floorStartedAt               Clock
	monstersKilledOnCurrentFloor int
	clearedFloors                int
	totalFloorClearDuration      Duration

	bossStartedAt    Clock
	bossKillDuration Duration
}

func NewPlayer(id int) *Player {
	return &Player{
		id: id,
		hp: maxHealth,
	}
}

// Getters

func (p *Player) ID() int                           { return p.id }
func (p *Player) HP() int                           { return p.hp }
func (p *Player) IsRegistered() bool                { return p.registered }
func (p *Player) IsInDungeon() bool                 { return p.inDungeon }
func (p *Player) IsFinished() bool                  { return p.finished }
func (p *Player) IsDisqualified() bool              { return p.disqualified }
func (p *Player) IsDead() bool                      { return p.dead }
func (p *Player) CurrentFloor() int                 { return p.currentFloor }
func (p *Player) IsOnBossFloor() bool               { return p.onBossFloor }
func (p *Player) IsBossKilled() bool                { return p.bossKilled }
func (p *Player) EnteredAt() Clock                  { return p.enteredAt }
func (p *Player) FinishedAt() Clock                 { return p.finishedAt }
func (p *Player) ClearedFloors() int                { return p.clearedFloors }
func (p *Player) TotalFloorClearDuration() Duration { return p.totalFloorClearDuration }
func (p *Player) BossKillDuration() Duration        { return p.bossKillDuration }
func (p *Player) MonstersKilledOnCurrentFloor() int { return p.monstersKilledOnCurrentFloor }

// Mutators

func (p *Player) Register() {
	p.registered = true
}

func (p *Player) EnterDungeon(at Clock) {
	p.inDungeon = true
	p.currentFloor = 1
	p.enteredAt = at
	p.floorStartedAt = at
}

func (p *Player) LeaveDungeon(at Clock) {
	p.inDungeon = false
	p.finished = true
	p.finishedAt = at
}

func (p *Player) MarkDisqualified(at Clock) {
	p.disqualified = true
	p.finished = true
	p.inDungeon = false
	p.finishedAt = at
}

func (p *Player) MarkDead(at Clock) {
	p.dead = true
	p.finished = true
	p.inDungeon = false
	p.finishedAt = at
	p.hp = 0
}

func (p *Player) MarkExpired(at Clock) {
	p.finished = true
	p.inDungeon = false
	p.finishedAt = at
}

func (p *Player) KillMonster(monstersPerFloor int, at Clock) {
	p.monstersKilledOnCurrentFloor++
	if monstersPerFloor == 0 || p.monstersKilledOnCurrentFloor >= monstersPerFloor {
		p.clearedFloors++
		p.totalFloorClearDuration = p.totalFloorClearDuration.Add(DurationBetween(p.floorStartedAt, at))
	}
}

func (p *Player) MoveNextFloor(at Clock) {
	p.currentFloor++
	p.monstersKilledOnCurrentFloor = 0
	p.floorStartedAt = at
}

func (p *Player) MovePreviousFloor(at Clock) {
	if p.currentFloor > 1 {
		p.currentFloor--
	}

	p.monstersKilledOnCurrentFloor = 0
	p.floorStartedAt = at
	p.onBossFloor = false
}

func (p *Player) EnterBossFloor(at Clock) {
	p.onBossFloor = true
	p.bossStartedAt = at
}

func (p *Player) KillBoss(at Clock) {
	p.bossKilled = true
	p.bossKillDuration = DurationBetween(p.bossStartedAt, at)
}

func (p *Player) RestoreHealth(health int) {
	p.hp += health
	if p.hp > maxHealth {
		p.hp = maxHealth
	}
}

func (p *Player) TakeDamage(damage int) bool {
	p.hp -= damage
	if p.hp <= 0 {
		p.hp = 0
		return true
	}

	return false
}

func (p *Player) HasEnteredDungeon() bool {
	return !p.enteredAt.IsZero() || p.inDungeon || (p.finished && p.currentFloor > 0)
}

func (p *Player) ClearedDungeon(dungeon Dungeon) bool {
	return p.clearedFloors >= dungeon.MonsterFloors() && p.bossKilled
}

func (p *Player) Status(dungeon Dungeon) PlayerStatus {
	if p.disqualified || !p.registered {
		return PlayerStatusDisqual
	}

	if p.ClearedDungeon(dungeon) && !p.dead {
		return PlayerStatusSuccess
	}

	return PlayerStatusFail
}

func MaxHealth() int {
	return maxHealth
}
