package config

import (
	"log/slog"
)

type Config struct {
	Log *LogConfig `yaml:"log"`
	Net *NetConfig `yaml:"net"`
}

type LogConfig struct {
	UpdaterLevel slog.Level `yaml:"notes_level"`
}

type NetConfig struct {
	ServerPort int `yaml:"server_port"`
	ClientPort int `yaml:"client_port"`
}
