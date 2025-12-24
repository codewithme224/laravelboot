package plugins

import (
	"fmt"
	"laravelboot/internal/config"
)

type Plugin interface {
	Name() string
	Install(conf *config.Config, projectPath string) error
}

type PluginManager struct {
	plugins []Plugin
}

func NewPluginManager() *PluginManager {
	return &PluginManager{
		plugins: make([]Plugin, 0),
	}
}

func (m *PluginManager) Register(p Plugin) {
	m.plugins = append(m.plugins, p)
}

func (m *PluginManager) RunAll(conf *config.Config, projectPath string) error {
	for _, p := range m.plugins {
		fmt.Printf("🔌 Running plugin: %s\n", p.Name())
		if err := p.Install(conf, projectPath); err != nil {
			return err
		}
	}
	return nil
}
