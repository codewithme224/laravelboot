package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectName  string   `yaml:"project_name"`
	Database     string   `yaml:"database"`     // mysql, postgres, sqlite, mongo
	Auth         string   `yaml:"auth"`         // sanctum, passport
	Features     []string `yaml:"features"`     // roles, media, search, activity-log
	Infra        []string `yaml:"infra"`        // docker, health, security, rate-limit
	Enterprise   []string `yaml:"enterprise"`   // quality, pro-arch, docs-pro, ci, monitoring
	Architecture string   `yaml:"architecture"` // domain-based, standard
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *Config) Save(path string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func DefaultConfig() *Config {
	return &Config{
		ProjectName:  "myapp",
		Database:     "mysql",
		Auth:         "sanctum",
		Features:     []string{"roles", "media"},
		Infra:        []string{"docker", "security", "health"},
		Enterprise:   []string{"quality"},
		Architecture: "domain-based",
	}
}
