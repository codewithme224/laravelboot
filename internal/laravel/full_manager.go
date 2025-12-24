package laravel

import "fmt"

type FullStackManager struct {
	DryRun bool
}

func NewFullStackManager(dryRun bool) *FullStackManager {
	return &FullStackManager{DryRun: dryRun}
}

func (m *FullStackManager) AddAll() error {
	fmt.Println("🌟 Installing the COMPLETE LaravelBoot Stack...")

	// 1. Auth & Database
	if err := NewAuthManager(m.DryRun).AddAuth(); err != nil {
		return err
	}

	// 2. Platform (Roles, Media, Search, Activity)
	if err := NewPlatformManager(m.DryRun).RunStep("platform"); err != nil {
		return err
	}

	// 3. Infrastructure (Docker, Security, Rate Limit, Health)
	if err := NewInfraManager(m.DryRun).RunStep("infra"); err != nil {
		return err
	}

	// 4. Enterprise (Quality, Pro-Arch, Docs, CI, monitoring)
	if err := NewEnterpriseManager(m.DryRun).RunStep("enterprise"); err != nil {
		return err
	}

	fmt.Println("\n🏆 CONGRATULATIONS! Your project is now fully loaded and production-ready.")
	return nil
}
