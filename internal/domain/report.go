package domain

type ReportRow struct {
	Status                PlayerStatus
	PlayerID              int
	TotalTime             Duration
	AverageFloorClearTime Duration
	BossKillTime          Duration
	HP                    int
}
