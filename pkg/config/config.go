package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Bot      `json:"bot"`
	Database `json:"database"`
}

type Bot struct {
	Token         string `json:"token"`
	DebugMode     bool   `json:"debug_mode"`
	UpdateOffset  int    `json:"update_offset"`
	UpdateTimeout int    `json:"update_timeout"`
}

type Database struct {
	Driver string `json:"driver"`
	Path   string `json:"path"`
}

func GetConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	var cfg Config
	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}
