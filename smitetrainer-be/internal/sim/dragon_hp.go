package sim

import "math"

const (
	patchBasisText          = "26.1 elemental drake HP: 3625 + 375 growth (lvl6-18)"
	defaultBurstStartHPPct  = 0.22
	defaultSimulationTickMs = int64(500)
	defaultWindowMs         = int64(20000)
	defaultBurstMs          = int64(1200)
)

type Options struct {
	TickMs          int64
	WindowMs        int64
	BurstMs         int64
	BurstStartHPPct float64
}

type HPPoint struct {
	TMs   int64   `json:"tMs"`
	HP    int     `json:"hp"`
	HPPct float64 `json:"hpPct"`
}

func PatchBasisText() string {
	return patchBasisText
}

func NormalizeOptions(opts Options) Options {
	if opts.TickMs < 100 {
		opts.TickMs = defaultSimulationTickMs
	}
	if opts.WindowMs < opts.TickMs {
		opts.WindowMs = defaultWindowMs
	}
	if opts.BurstMs < 0 {
		opts.BurstMs = 0
	}
	if opts.BurstStartHPPct <= 0 || opts.BurstStartHPPct >= 1 {
		opts.BurstStartHPPct = defaultBurstStartHPPct
	}
	return opts
}

func DefaultOptions() Options {
	return Options{
		TickMs:          defaultSimulationTickMs,
		WindowMs:        defaultWindowMs,
		BurstMs:         defaultBurstMs,
		BurstStartHPPct: defaultBurstStartHPPct,
	}
}

func EstimateDragonLevel(killMs int64) int {
	gameMinutes := float64(killMs) / 60000.0
	if gameMinutes < 0 {
		gameMinutes = 0
	}
	if gameMinutes > 40 {
		gameMinutes = 40
	}
	level := int(math.Round(6 + (12 * gameMinutes / 40)))
	return clampInt(level, 6, 18)
}

func BaseElementalDragonHP(level int) int {
	level = clampInt(level, 6, 18)
	return 3625 + 375*(level-6)
}

func BuildSeries(killMs int64, hp0 int, opts Options) (int64, []HPPoint, string) {
	opts = NormalizeOptions(opts)

	startMs := killMs - opts.WindowMs
	if startMs < 0 {
		startMs = 0
	}

	if killMs <= startMs {
		return startMs, []HPPoint{
			{
				TMs:   killMs,
				HP:    0,
				HPPct: 0,
			},
		}, "linear"
	}

	useBurst := opts.BurstMs > 0 && opts.BurstMs < (killMs-startMs)
	modelName := "linear"
	if useBurst {
		modelName = "linear+endBurst"
	}

	points := make([]HPPoint, 0, int((killMs-startMs)/opts.TickMs)+2)
	for t := startMs; t < killMs; t += opts.TickMs {
		hpPct := calculateHPPct(t, startMs, killMs, opts, useBurst)
		points = append(points, HPPoint{
			TMs:   t,
			HP:    int(math.Round(float64(hp0) * hpPct)),
			HPPct: round4(hpPct),
		})
	}

	points = append(points, HPPoint{
		TMs:   killMs,
		HP:    0,
		HPPct: 0,
	})

	return startMs, points, modelName
}

func calculateHPPct(t, startMs, killMs int64, opts Options, useBurst bool) float64 {
	if !useBurst {
		frac := float64(t-startMs) / float64(killMs-startMs)
		return clampFloat(1-frac, 0, 1)
	}

	burstStart := killMs - opts.BurstMs
	if t <= burstStart {
		preDuration := burstStart - startMs
		if preDuration <= 0 {
			frac := float64(t-startMs) / float64(killMs-startMs)
			return clampFloat(1-frac, 0, 1)
		}
		preFrac := float64(t-startMs) / float64(preDuration)
		return clampFloat(1-preFrac*(1-opts.BurstStartHPPct), 0, 1)
	}

	burstDuration := killMs - burstStart
	if burstDuration <= 0 {
		return 0
	}
	burstFrac := float64(t-burstStart) / float64(burstDuration)
	hpPct := opts.BurstStartHPPct * (1 - burstFrac)
	return clampFloat(hpPct, 0, 1)
}

func clampInt(v, minV, maxV int) int {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func clampFloat(v, minV, maxV float64) float64 {
	if v < minV {
		return minV
	}
	if v > maxV {
		return maxV
	}
	return v
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
