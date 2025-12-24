package laravel

import (
	"fmt"
	"os"
)

type AuthManager struct {
	DryRun bool
}

func NewAuthManager(dryRun bool) *AuthManager {
	return &AuthManager{DryRun: dryRun}
}

func (m *AuthManager) AddAuth() error {
	fmt.Println("🔐 Adding Authentication and Database features...")

	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	// In a real scenario, we might want to verify if we are in a Laravel project
	// For Phase 2, we assume the user runs this from the project root.

	// 1. Install Sanctum/API
	sanctum := NewSanctumInstaller(projectPath, m.DryRun)
	if err := sanctum.Install(); err != nil {
		return err
	}

	// 2. Pagination & Support
	pagination := NewPaginationSetup(projectPath, m.DryRun)
	if err := pagination.Setup(); err != nil {
		return err
	}

	// 3. Auth Logic
	auth := NewAuthSetup(projectPath, m.DryRun)
	if err := auth.Setup(); err != nil {
		return err
	}

	// 4. Database
	db := NewDatabaseSetup(projectPath, m.DryRun)
	if err := db.RunMigrations(); err != nil {
		return err
	}

	fmt.Println("\n✅ Authentication and API features added successfully!")
	return nil
}
