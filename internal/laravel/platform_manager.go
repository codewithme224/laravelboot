package laravel

import (
	"fmt"
)

type PlatformManager struct {
	ProjectPath string
	DryRun      bool
}

func NewPlatformManager(projectPath string, dryRun bool) *PlatformManager {
	return &PlatformManager{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *PlatformManager) RunStep(name string) error {
	switch name {
	case "roles":
		return NewRolesSetup(m.ProjectPath, m.DryRun).Setup()
	case "media":
		return NewMediaSetup(m.ProjectPath, m.DryRun).Setup()
	case "activity-log":
		return NewActivityLogSetup(m.ProjectPath, m.DryRun).Setup()
	case "search":
		return NewSearchSetup(m.ProjectPath, m.DryRun).Setup()
	case "platform":
		fmt.Println("🚀 Installing complete platform stack...")
		if err := NewRolesSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewMediaSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewActivityLogSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewSearchSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown platform feature: %s", name)
	}
}
