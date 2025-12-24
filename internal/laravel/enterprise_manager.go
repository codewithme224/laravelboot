package laravel

import (
	"fmt"
	"os"
)

type EnterpriseManager struct {
	DryRun bool
}

func NewEnterpriseManager(dryRun bool) *EnterpriseManager {
	return &EnterpriseManager{DryRun: dryRun}
}

func (m *EnterpriseManager) RunStep(name string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	switch name {
	case "quality":
		return NewQualitySetup(projectPath, m.DryRun).Setup()
	case "pro-arch":
		return NewProArchSetup(projectPath, m.DryRun).Setup()
	case "docs-pro":
		return NewDocsProSetup(projectPath, m.DryRun).Setup()
	case "ci":
		return NewCicdSetup(projectPath, m.DryRun).Setup()
	case "monitoring":
		return NewMonitoringSetup(projectPath, m.DryRun).Setup()
	case "enterprise":
		fmt.Println("👑 Installing complete Enterprise stack...")
		if err := NewQualitySetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewProArchSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewDocsProSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewCicdSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewMonitoringSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown enterprise feature: %s", name)
	}
}
