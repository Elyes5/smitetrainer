package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"smitetrainer-be/internal/api"
	"smitetrainer-be/internal/cache"
	"smitetrainer-be/internal/config"
	"smitetrainer-be/internal/riotclient"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	memoryCache := cache.NewByteLRU(cfg.CacheMaxEntries, cfg.CacheTTL)
	var responseCache cache.Store = memoryCache

	if cfg.CacheBackend == "redis" {
		redisStore, err := cache.NewRedisStore(cache.RedisConfig{
			Addr:        cfg.RedisAddr,
			Username:    cfg.RedisUsername,
			Password:    cfg.RedisPassword,
			DB:          cfg.RedisDB,
			DialTimeout: cfg.RedisDialTimeout,
			IOTimeout:   cfg.RedisIOTimeout,
			DefaultTTL:  cfg.CacheTTL,
		})
		if err != nil {
			log.Printf("redis cache config invalid (%v), falling back to memory", err)
		} else {
			responseCache = cache.NewFallbackStore(redisStore, memoryCache)
			log.Printf("cache backend: redis (memory fallback enabled)")
		}
	} else {
		log.Printf("cache backend: memory")
	}
	defer func() {
		if err := responseCache.Close(); err != nil {
			log.Printf("cache close error: %v", err)
		}
	}()

	riot := riotclient.New(cfg.RiotAPIKey, cfg.HTTPTimeout, cfg.RiotMaxAttempts, responseCache, cfg.CacheTTL)
	handler := api.New(riot, cfg)

	mux := http.NewServeMux()
	handler.Register(mux)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           withCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
