package laravel

import (
	"fmt"
)

type EnterpriseManager struct {
	ProjectPath string
	DryRun      bool
}

func NewEnterpriseManager(projectPath string, dryRun bool) *EnterpriseManager {
	return &EnterpriseManager{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *EnterpriseManager) RunStep(name string) error {
	switch name {
	case "quality":
		return NewQualitySetup(m.ProjectPath, m.DryRun).Setup()
	case "pro-arch":
		return NewProArchSetup(m.ProjectPath, m.DryRun).Setup()
	case "docs-pro":
		return NewDocsProSetup(m.ProjectPath, m.DryRun).Setup()
	case "ci":
		return NewCicdSetup(m.ProjectPath, m.DryRun).Setup()
	case "monitoring":
		return NewMonitoringSetup(m.ProjectPath, m.DryRun).Setup()
	case "tenancy":
		return NewTenancySetup(m.ProjectPath, m.DryRun).Setup()
	case "helpers":
		return NewHelpersSetup(m.ProjectPath, m.DryRun).Setup()
	case "enterprise":
		fmt.Println("👑 Installing complete Enterprise stack...")
		if err := NewQualitySetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewProArchSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewDocsProSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewCicdSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewMonitoringSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewHelpersSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown enterprise feature: %s", name)
	}
}
