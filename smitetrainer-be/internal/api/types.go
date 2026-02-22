package api

import (
	"smitetrainer-be/internal/parser"
	"smitetrainer-be/internal/riotclient"
	"smitetrainer-be/internal/sim"
)

type dragonHPResponse struct {
	MatchID     string              `json:"matchId"`
	DragonIndex int                 `json:"dragonIndex"`
	DragonType  string              `json:"dragonType"`
	KillMs      int64               `json:"killMs"`
	StartMs     int64               `json:"startMs"`
	KillEvent   killEventDTO        `json:"killEvent"`
	HPModel     hpModelDTO          `json:"hpModel"`
	Series      []sim.HPPoint       `json:"series"`
	Markers     parser.FightMarkers `json:"markers"`
}

type killEventDTO struct {
	KillerID                int                 `json:"killerId"`
	AssistingParticipantIDs []int               `json:"assistingParticipantIds"`
	Position                riotclient.Position `json:"position"`
}

type hpModelDTO struct {
	IsEstimated      bool    `json:"isEstimated"`
	PatchBasis       string  `json:"patchBasis"`
	LevelEstimate    int     `json:"levelEstimate"`
	HP0              int     `json:"hp0"`
	TickMs           int64   `json:"tickMs"`
	WindowMs         int64   `json:"windowMs"`
	BurstMs          int64   `json:"burstMs"`
	BurstStartHPPct  float64 `json:"burstStartHpPct"`
	Model            string  `json:"model"`
	SourceGameLength int64   `json:"sourceGameLength"`
}
