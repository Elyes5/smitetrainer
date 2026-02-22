package riotclient

import (
	"fmt"
	"strings"
)

func RegionalRouteFromMatchID(matchID string) (string, error) {
	parts := strings.SplitN(matchID, "_", 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid match id format: %s", matchID)
	}

	platform := strings.ToUpper(strings.TrimSpace(parts[0]))
	switch platform {
	case "AMERICAS", "EUROPE", "ASIA", "SEA":
		return strings.ToLower(platform), nil
	case "BR1", "LA1", "LA2", "NA1":
		return "americas", nil
	case "EUN1", "EUW1", "ME1", "RU", "TR1":
		return "europe", nil
	case "JP1", "KR":
		return "asia", nil
	case "OCE1", "PH2", "SG2", "TH2", "TW2", "VN2":
		return "sea", nil
	default:
		return "", fmt.Errorf("unsupported platform in match id: %s", platform)
	}
}
