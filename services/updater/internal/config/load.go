package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

func LoadFileConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return &Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return &Config{}, fmt.Errorf("parse yaml %q: %w", path, err)
	}

	/*
		if err := cfg.Validate(); err != nil {
			return &Config{}, err
		}
	*/
	return &cfg, nil
}

func LoadEnvConfig() *Config {
	return &Config{
		Log: &LogConfig{
			UpdaterLevel: slog.Level(getEnvInt("UPDATER_LOG_LEVEL")),
		}, Net: &NetConfig{
			ServerPort: getEnvInt("UPDATER_SERVER_PORT"),
			ClientPort: getEnvInt("NOTES_SERVER_PORT"),
		},
	}
}

func getEnv(k string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return ""
}
func getEnvInt(k string) int {
	v := os.Getenv(k)
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}
func getEnvDur(k string, def time.Duration) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
