package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Api      APIConfig      `yaml:"api"`
	Database DatabaseConfig `yaml:"database"`
	Session  SessionConfig  `yaml:"session"`
}

type APIConfig struct {
	Port    int    `yaml:"port"`
	Env     string `yaml:"env"`
	Version string `yaml:"version"`
	Name    string `yaml:"name"`
}

type DatabaseConfig struct {
	Dsn          string `yaml:"dsn"`
	MaxOpenConns int    `yaml:"maxOpenConns"`
	MaxIdleConns int    `yaml:"maxIdleConns"`
	MaxIdleTime  string `yaml:"maxIdleTime"`
}

type SessionConfig struct {
	Lifetime string `yaml:"lifetime"`
}

func New(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
