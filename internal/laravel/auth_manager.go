package laravel

import (
	"fmt"
)

type AuthManager struct {
	ProjectPath string
	DryRun      bool
}

func NewAuthManager(projectPath string, dryRun bool) *AuthManager {
	return &AuthManager{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *AuthManager) AddAuth() error {
	fmt.Println("🔐 Adding Authentication and Database features...")

	// 1. Install Sanctum/API
	sanctum := NewSanctumInstaller(m.ProjectPath, m.DryRun)
	if err := sanctum.Install(); err != nil {
		return err
	}

	// 2. Pagination & Support
	pagination := NewPaginationSetup(m.ProjectPath, m.DryRun)
	if err := pagination.Setup(); err != nil {
		return err
	}

	// 3. Auth Logic
	auth := NewAuthSetup(m.ProjectPath, m.DryRun)
	if err := auth.Setup(); err != nil {
		return err
	}

	// 4. Database
	db := NewDatabaseSetup(m.ProjectPath, m.DryRun)
	if err := db.RunMigrations(); err != nil {
		return err
	}

	fmt.Println("\n✅ Authentication and API features added successfully!")
	return nil
}
