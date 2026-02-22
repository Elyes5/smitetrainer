package riotclient

type MatchResponse struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     MatchInfo     `json:"info"`
}

type MatchMetadata struct {
	MatchID      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type MatchInfo struct {
	GameCreation int64  `json:"gameCreation"`
	GameDuration int64  `json:"gameDuration"`
	GameVersion  string `json:"gameVersion"`
}

type TimelineResponse struct {
	Metadata MatchMetadata `json:"metadata"`
	Info     TimelineInfo  `json:"info"`
}

type TimelineInfo struct {
	Frames []TimelineFrame `json:"frames"`
}

type TimelineFrame struct {
	Timestamp int64           `json:"timestamp"`
	Events    []TimelineEvent `json:"events"`
}

type TimelineEvent struct {
	Type                    string   `json:"type"`
	Timestamp               int64    `json:"timestamp"`
	MonsterType             string   `json:"monsterType"`
	MonsterSubType          string   `json:"monsterSubType"`
	KillerID                int      `json:"killerId"`
	AssistingParticipantIDs []int    `json:"assistingParticipantIds"`
	Position                Position `json:"position"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
