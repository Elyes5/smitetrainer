package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port             string
	RiotAPIKey       string
	HTTPTimeout      time.Duration
	RiotMaxAttempts  int
	CacheBackend     string
	CacheTTL         time.Duration
	CacheMaxEntries  int
	RedisAddr        string
	RedisUsername    string
	RedisPassword    string
	RedisDB          int
	RedisDialTimeout time.Duration
	RedisIOTimeout   time.Duration
	DefaultTickMs    int64
	DefaultWindowMs  int64
	DefaultBurstMs   int64
	FightRadius      float64
}

func Load() (Config, error) {
	loadDotEnv(".env")

	cfg := Config{
		Port:             getEnvString("PORT", "8080"),
		HTTPTimeout:      time.Duration(getEnvInt("RIOT_HTTP_TIMEOUT_SECONDS", 8)) * time.Second,
		RiotMaxAttempts:  getEnvInt("RIOT_MAX_ATTEMPTS", 5),
		CacheBackend:     strings.ToLower(getEnvString("CACHE_BACKEND", "memory")),
		CacheTTL:         time.Duration(getEnvInt("CACHE_TTL_SECONDS", 300)) * time.Second,
		CacheMaxEntries:  getEnvInt("CACHE_MAX_ENTRIES", 256),
		RedisAddr:        getEnvString("REDIS_ADDR", "127.0.0.1:6379"),
		RedisUsername:    getEnvString("REDIS_USERNAME", ""),
		RedisPassword:    getEnvString("REDIS_PASSWORD", ""),
		RedisDB:          getEnvInt("REDIS_DB", 0),
		RedisDialTimeout: time.Duration(getEnvInt("REDIS_DIAL_TIMEOUT_SECONDS", 3)) * time.Second,
		RedisIOTimeout:   time.Duration(getEnvInt("REDIS_IO_TIMEOUT_SECONDS", 3)) * time.Second,
		DefaultTickMs:    int64(getEnvInt("DEFAULT_TICK_MS", 500)),
		DefaultWindowMs:  int64(getEnvInt("DEFAULT_WINDOW_MS", 20000)),
		DefaultBurstMs:   int64(getEnvInt("DEFAULT_BURST_MS", 1200)),
		FightRadius:      getEnvFloat("FIGHT_RADIUS", 2500.0),
	}

	cfg.RiotAPIKey = strings.TrimSpace(os.Getenv("RIOT_API_KEY"))
	cfg.RiotAPIKey = strings.Trim(cfg.RiotAPIKey, "\"")
	if cfg.RiotAPIKey == "" {
		return Config{}, fmt.Errorf("RIOT_API_KEY is required")
	}

	if cfg.RiotMaxAttempts < 1 {
		cfg.RiotMaxAttempts = 1
	}
	if cfg.CacheBackend != "memory" && cfg.CacheBackend != "redis" {
		return Config{}, fmt.Errorf("CACHE_BACKEND must be one of: memory, redis")
	}
	if cfg.CacheTTL < 0 {
		cfg.CacheTTL = 0
	}
	if cfg.CacheMaxEntries < 1 {
		cfg.CacheMaxEntries = 1
	}
	if cfg.RedisDB < 0 {
		return Config{}, fmt.Errorf("REDIS_DB must be >= 0")
	}
	if cfg.RedisDialTimeout <= 0 {
		cfg.RedisDialTimeout = 3 * time.Second
	}
	if cfg.RedisIOTimeout <= 0 {
		cfg.RedisIOTimeout = 3 * time.Second
	}
	if cfg.CacheBackend == "redis" && strings.TrimSpace(cfg.RedisAddr) == "" {
		return Config{}, fmt.Errorf("REDIS_ADDR is required when CACHE_BACKEND=redis")
	}
	if cfg.DefaultTickMs < 100 {
		cfg.DefaultTickMs = 100
	}
	if cfg.DefaultWindowMs < cfg.DefaultTickMs {
		cfg.DefaultWindowMs = cfg.DefaultTickMs
	}
	if cfg.DefaultBurstMs < 0 {
		cfg.DefaultBurstMs = 0
	}
	if cfg.FightRadius < 0 {
		cfg.FightRadius = 0
	}

	return cfg, nil
}

func getEnvString(key, fallback string) string {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	return val
}

func getEnvInt(key string, fallback int) int {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvFloat(key string, fallback float64) float64 {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	n, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return fallback
	}
	return n
}

func loadDotEnv(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		idx := strings.Index(line, "=")
		if idx <= 0 {
			continue
		}

		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		value = strings.Trim(value, "\"")

		if key == "" {
			continue
		}
		if _, exists := os.LookupEnv(key); exists {
			continue
		}
		_ = os.Setenv(key, value)
	}
}
