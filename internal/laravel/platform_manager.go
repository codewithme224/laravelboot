package laravel

import (
	"fmt"
	"os"
)

type PlatformManager struct {
	DryRun bool
}

func NewPlatformManager(dryRun bool) *PlatformManager {
	return &PlatformManager{DryRun: dryRun}
}

func (m *PlatformManager) RunStep(name string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	switch name {
	case "roles":
		return NewRolesSetup(projectPath, m.DryRun).Setup()
	case "media":
		return NewMediaSetup(projectPath, m.DryRun).Setup()
	case "activity-log":
		return NewActivityLogSetup(projectPath, m.DryRun).Setup()
	case "search":
		return NewSearchSetup(projectPath, m.DryRun).Setup()
	case "platform":
		fmt.Println("🚀 Installing complete platform stack...")
		if err := NewRolesSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewMediaSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewActivityLogSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewSearchSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown platform feature: %s", name)
	}
}
