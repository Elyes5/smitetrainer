package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"smitetrainer-be/internal/config"
	"smitetrainer-be/internal/parser"
	"smitetrainer-be/internal/riotclient"
	"smitetrainer-be/internal/sim"
)

type Handler struct {
	riot *riotclient.Client
	cfg  config.Config
}

func New(riot *riotclient.Client, cfg config.Config) *Handler {
	return &Handler{
		riot: riot,
		cfg:  cfg,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", h.healthz)
	mux.HandleFunc("/v1/matches/", h.getDragonHP)
}

func (h *Handler) healthz(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) getDragonHP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	matchID, dragonIndex, ok := parseDragonRoute(r.URL.Path)
	if !ok {
		writeError(w, http.StatusNotFound, "endpoint not found")
		return
	}

	tickMs := parseInt64OrDefault(r.URL.Query().Get("tickMs"), h.cfg.DefaultTickMs)
	windowMs := parseInt64OrDefault(r.URL.Query().Get("windowMs"), h.cfg.DefaultWindowMs)
	if tickMs < 100 {
		tickMs = 100
	}
	if windowMs < tickMs {
		windowMs = tickMs
	}
	if windowMs > 120000 {
		windowMs = 120000
	}

	var (
		match    riotclient.MatchResponse
		timeline riotclient.TimelineResponse
		matchErr error
		tlErr    error
	)

	ctx, cancel := context.WithTimeout(r.Context(), h.cfg.HTTPTimeout)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		match, matchErr = h.riot.GetMatch(ctx, matchID)
	}()
	go func() {
		defer wg.Done()
		timeline, tlErr = h.riot.GetTimeline(ctx, matchID)
	}()
	wg.Wait()

	if matchErr != nil {
		status, message := mapUpstreamError("match", matchErr)
		writeError(w, status, message)
		return
	}
	if tlErr != nil {
		status, message := mapUpstreamError("timeline", tlErr)
		writeError(w, status, message)
		return
	}

	killEvent, err := parser.FindDragonKill(timeline, dragonIndex)
	if err != nil {
		if errors.Is(err, parser.ErrDragonNotFound) {
			writeError(w, http.StatusNotFound, "requested dragon index not found in timeline")
			return
		}
		writeError(w, http.StatusInternalServerError, "failed to parse timeline")
		return
	}

	opts := sim.NormalizeOptions(sim.Options{
		TickMs:          tickMs,
		WindowMs:        windowMs,
		BurstMs:         h.cfg.DefaultBurstMs,
		BurstStartHPPct: 0.22,
	})

	levelEstimate := sim.EstimateDragonLevel(killEvent.KillMs)
	hp0 := sim.BaseElementalDragonHP(levelEstimate)
	startMs, series, modelName := sim.BuildSeries(killEvent.KillMs, hp0, opts)
	markers := parser.ExtractFightMarkers(timeline, startMs, killEvent.KillMs, h.cfg.FightRadius)

	dragonType := killEvent.DragonSubType
	if dragonType == "" {
		dragonType = "DRAGON"
	}

	patchBasis := sim.PatchBasisText()
	if strings.TrimSpace(match.Info.GameVersion) != "" {
		patchBasis = fmt.Sprintf("%s (match.gameVersion=%s)", patchBasis, match.Info.GameVersion)
	}

	response := dragonHPResponse{
		MatchID:     matchID,
		DragonIndex: dragonIndex,
		DragonType:  dragonType,
		KillMs:      killEvent.KillMs,
		StartMs:     startMs,
		KillEvent: killEventDTO{
			KillerID:                killEvent.KillerID,
			AssistingParticipantIDs: killEvent.AssistingParticipantIDs,
			Position:                killEvent.Position,
		},
		HPModel: hpModelDTO{
			IsEstimated:      true,
			PatchBasis:       patchBasis,
			LevelEstimate:    levelEstimate,
			HP0:              hp0,
			TickMs:           opts.TickMs,
			WindowMs:         opts.WindowMs,
			BurstMs:          opts.BurstMs,
			BurstStartHPPct:  opts.BurstStartHPPct,
			Model:            modelName,
			SourceGameLength: match.Info.GameDuration,
		},
		Series:  series,
		Markers: markers,
	}

	writeJSON(w, http.StatusOK, response)
}

func parseDragonRoute(path string) (string, int, bool) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) != 6 {
		return "", 0, false
	}
	if parts[0] != "v1" || parts[1] != "matches" || parts[3] != "dragons" || parts[5] != "hp" {
		return "", 0, false
	}

	matchID := parts[2]
	if matchID == "" {
		return "", 0, false
	}

	dragonIndex, err := strconv.Atoi(parts[4])
	if err != nil || dragonIndex < 1 {
		return "", 0, false
	}
	return matchID, dragonIndex, true
}

func parseInt64OrDefault(raw string, fallback int64) int64 {
	if strings.TrimSpace(raw) == "" {
		return fallback
	}
	n, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return fallback
	}
	return n
}

func mapUpstreamError(resource string, err error) (int, string) {
	var apiErr *riotclient.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case http.StatusNotFound:
			if resource == "timeline" {
				return http.StatusNotFound, "timeline not available for this match"
			}
			return http.StatusNotFound, "match not found"
		case http.StatusTooManyRequests:
			return http.StatusTooManyRequests, "riot rate limit exceeded"
		default:
			if apiErr.StatusCode >= 500 {
				return http.StatusBadGateway, "riot service temporarily unavailable"
			}
			return http.StatusBadGateway, fmt.Sprintf("riot api error while fetching %s", resource)
		}
	}
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return http.StatusGatewayTimeout, "upstream request timed out"
	}
	return http.StatusBadGateway, fmt.Sprintf("failed to fetch %s", resource)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
