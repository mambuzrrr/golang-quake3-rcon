package examples

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Config struct {
	Address   string `json:"address"`
	Password  string `json:"password"`
	Debug     bool   `json:"debug"`
	TimeoutMs int    `json:"timeout_ms"`
	QuietMs   int    `json:"quiet_ms"`
}

func LoadConfig(path string) (Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read %s: %w", path, err)
	}

	var cfg Config
	if err := json.Unmarshal(b, &cfg); err != nil {
		return Config{}, fmt.Errorf("parse %s: %w", path, err)
	}

	cfg.Address = strings.TrimSpace(cfg.Address)
	if cfg.Address == "" {
		return Config{}, fmt.Errorf("%s: address is empty", path)
	}
	if !strings.Contains(cfg.Address, ":") {
		return Config{}, fmt.Errorf("%s: address must include port (ip:port)", path)
	}
	if strings.TrimSpace(cfg.Password) == "" {
		return Config{}, fmt.Errorf("%s: password is empty", path)
	}

	return cfg, nil
}

func MsOrDefault(ms int, def int) time.Duration {
	if ms <= 0 {
		ms = def
	}
	return time.Duration(ms) * time.Millisecond
}
