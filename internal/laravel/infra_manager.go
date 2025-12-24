package laravel

import (
	"fmt"
)

type InfraManager struct {
	ProjectPath string
	DryRun      bool
}

func NewInfraManager(projectPath string, dryRun bool) *InfraManager {
	return &InfraManager{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *InfraManager) RunStep(name string) error {
	switch name {
	case "docker":
		return NewDockerSetup(m.ProjectPath, m.DryRun).Setup()
	case "security":
		return NewSecuritySetup(m.ProjectPath, m.DryRun).Setup()
	case "rate-limit":
		return NewRateLimitSetup(m.ProjectPath, m.DryRun).Setup()
	case "health":
		return NewHealthSetup(m.ProjectPath, m.DryRun).Setup()
	case "infra":
		fmt.Println("🚀 Hardening infrastructure and security...")
		if err := NewDockerSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewSecuritySetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewRateLimitSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewHealthSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown infra feature: %s", name)
	}
}
