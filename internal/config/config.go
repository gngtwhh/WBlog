package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// 全局变量，加载后通过 config.Cfg 访问
var Cfg *Config

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
	App      AppConfig      `json:"app"`
}

type ServerConfig struct {
	Port    string `json:"port"`
	RunMode string `json:"run_mode"`
	// ReadTimeout  int    `json:"read_timeout"`  // second
	// WriteTimeout int    `json:"write_timeout"` // second
}

type DatabaseConfig struct {
	DSN string `json:"dsn"`
	// MaxOpenConns int    `json:"max_open_conns"`
	// MaxIdleConns int    `json:"max_idle_conns"`
}

type AppConfig struct {
	TemplateDir string `json:"template_dir"`
	StaticDir   string `json:"static_dir"`
	LogFile     string `json:"log_file"`
	JwtSecret   string `json:"jwt_secret"`
}

func Load(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("config file not exists: %s", filePath)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	cfg := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(cfg); err != nil {
		return fmt.Errorf("parse file to config failed: %w", err)
	}

	if cfg.App.JwtSecret == "" {
		return fmt.Errorf("jwt secret is empty")
	}

	Cfg = cfg
	return nil
}
