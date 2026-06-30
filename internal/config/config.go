package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port            int
	ConfigSecret    string
	CacheTTL        time.Duration
	MaxCacheEntries int
	Debug           bool
	PrefetchEnabled bool
	PrefetchMaxSize int64
}

func Load() *Config {
	return &Config{
		Port:            getEnvInt("PORT", 7001),
		ConfigSecret:    os.Getenv("CONFIG_SECRET"),
		CacheTTL:        getEnvDuration("CACHE_TTL", 6*time.Hour),
		MaxCacheEntries: getEnvInt("MAX_CACHE_ENTRIES", 500),
		Debug:           getEnvBool("DEBUG", false),
		PrefetchEnabled: getEnvBool("PREFETCH_ENABLED", true),
		PrefetchMaxSize: getEnvInt64("PREFETCH_MAX_SIZE", 150*1024*1024),
	}
}

func getEnvInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvInt64(key string, defaultVal int64) int64 {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.ParseInt(val, 10, 64); err == nil {
			return i
		}
	}
	return defaultVal
}

func getEnvBool(key string, defaultVal bool) bool {
	if val := os.Getenv(key); val != "" {
		if b, err := strconv.ParseBool(val); err == nil {
			return b
		}
	}
	return defaultVal
}

func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if d, err := time.ParseDuration(val); err == nil {
			return d
		}
	}
	return defaultVal
}
