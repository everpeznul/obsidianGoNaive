package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

func Load(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse yaml %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	// Минимальная валидация: чтобы падать при старте, а не “где-то потом”.
	if c.DB.DBName == "" {
		return fmt.Errorf("invalid db config: host/port/user/name are required")
	}

	return nil
}

func LoadDBFromEnv() Config {
	return Config{DBConfig{
		Host:     getEnv("POSTGRES_HOST"),
		Port:     getEnvInt("POSTGRES_PORT"),
		User:     getEnv("POSTGRES_USER"),
		Password: getEnv("POSTGRES_PASSWORD"),
		DBName:   getEnv("POSTGRES_DB"),
		SSLMode:  getEnv("POSTGRES_SSLMODE"),
	}}
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
