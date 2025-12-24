package laravel

import (
	"fmt"
	"os"
)

type InfraManager struct {
	DryRun bool
}

func NewInfraManager(dryRun bool) *InfraManager {
	return &InfraManager{DryRun: dryRun}
}

func (m *InfraManager) RunStep(name string) error {
	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	switch name {
	case "docker":
		return NewDockerSetup(projectPath, m.DryRun).Setup()
	case "security":
		return NewSecuritySetup(projectPath, m.DryRun).Setup()
	case "rate-limit":
		return NewRateLimitSetup(projectPath, m.DryRun).Setup()
	case "health":
		return NewHealthSetup(projectPath, m.DryRun).Setup()
	case "infra":
		fmt.Println("🚀 Hardening infrastructure and security...")
		if err := NewDockerSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewSecuritySetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewRateLimitSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewHealthSetup(projectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown infra feature: %s", name)
	}
}
