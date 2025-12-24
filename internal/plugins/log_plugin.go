package plugins

import (
	"fmt"
	"laravelboot/internal/config"
)

type LogPlugin struct{}

func (p *LogPlugin) Name() string {
	return "Logger"
}

func (p *LogPlugin) Install(conf *config.Config, projectPath string) error {
	fmt.Printf("📝 Logger Plugin: Setup complete for %s\n", conf.ProjectName)
	return nil
}
