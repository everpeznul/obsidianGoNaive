package config

import (
	"fmt"
	"log/slog"
	"time"
)

type Config struct {
	DB  DBConfig  `yaml:"db"`
	Log LogConfig `yaml:"log"`
	Net NetConfig `yaml:"net"`
}

type DBConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`

	MaxOpenConns    int           `yaml:"max_open_conns"`
	MaxIdleConns    int           `yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `yaml:"conn_max_lifetime"`
}

func (c DBConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

type LogConfig struct {
	NotesLevel slog.Level `yaml:"notes_level"`
}

type NetConfig struct {
	ServerPort int `yaml:"server_port"`
}
