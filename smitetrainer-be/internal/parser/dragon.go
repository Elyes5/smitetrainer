package parser

import (
	"errors"
	"math"
	"sort"
	"strings"

	"smitetrainer-be/internal/riotclient"
)

var ErrDragonNotFound = errors.New("requested dragon kill was not found")

const (
	dragonPitX = 9866
	dragonPitY = 4414
)

type DragonKillEvent struct {
	KillMs                  int64
	DragonSubType           string
	KillerID                int
	AssistingParticipantIDs []int
	Position                riotclient.Position
}

type FightMarkers struct {
	ChampionKillsInWindow int     `json:"championKillsInWindow"`
	ChampionKillsNearPit  int     `json:"championKillsNearPit"`
	DistanceThreshold     float64 `json:"distanceThreshold"`
}

func FindDragonKill(timeline riotclient.TimelineResponse, dragonIndex int) (DragonKillEvent, error) {
	if dragonIndex < 1 {
		return DragonKillEvent{}, ErrDragonNotFound
	}

	type orderedDragonEvent struct {
		timestamp int64
		event     riotclient.TimelineEvent
	}

	events := make([]orderedDragonEvent, 0, 8)
	for _, frame := range timeline.Info.Frames {
		for _, ev := range frame.Events {
			if ev.Type != "ELITE_MONSTER_KILL" {
				continue
			}
			if strings.ToUpper(ev.MonsterType) != "DRAGON" {
				continue
			}
			if strings.EqualFold(ev.MonsterSubType, "ELDER_DRAGON") {
				continue
			}

			ts := ev.Timestamp
			if ts == 0 {
				ts = frame.Timestamp
			}

			events = append(events, orderedDragonEvent{
				timestamp: ts,
				event:     ev,
			})
		}
	}

	if len(events) < dragonIndex {
		return DragonKillEvent{}, ErrDragonNotFound
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].timestamp < events[j].timestamp
	})

	selected := events[dragonIndex-1]
	assists := make([]int, len(selected.event.AssistingParticipantIDs))
	copy(assists, selected.event.AssistingParticipantIDs)

	return DragonKillEvent{
		KillMs:                  selected.timestamp,
		DragonSubType:           selected.event.MonsterSubType,
		KillerID:                selected.event.KillerID,
		AssistingParticipantIDs: assists,
		Position:                selected.event.Position,
	}, nil
}

func ExtractFightMarkers(timeline riotclient.TimelineResponse, startMs, endMs int64, distanceThreshold float64) FightMarkers {
	markers := FightMarkers{
		DistanceThreshold: distanceThreshold,
	}

	if endMs < startMs {
		return markers
	}

	for _, frame := range timeline.Info.Frames {
		for _, ev := range frame.Events {
			if ev.Type != "CHAMPION_KILL" {
				continue
			}

			ts := ev.Timestamp
			if ts == 0 {
				ts = frame.Timestamp
			}
			if ts < startMs || ts > endMs {
				continue
			}

			markers.ChampionKillsInWindow++
			if distanceThreshold <= 0 {
				continue
			}

			dist := distanceToDragonPit(ev.Position)
			if dist <= distanceThreshold {
				markers.ChampionKillsNearPit++
			}
		}
	}

	return markers
}

func distanceToDragonPit(pos riotclient.Position) float64 {
	dx := float64(pos.X - dragonPitX)
	dy := float64(pos.Y - dragonPitY)
	return math.Hypot(dx, dy)
}
