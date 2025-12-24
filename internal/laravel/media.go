package laravel

import (
	"fmt"
	"os/exec"
)

type MediaSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewMediaSetup(projectPath string, dryRun bool) *MediaSetup {
	return &MediaSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *MediaSetup) Setup() error {
	if m.DryRun {
		fmt.Printf("[Dry Run] Would install spatie/laravel-medialibrary\n")
		return nil
	}

	fmt.Println("🖼️ Installing spatie/laravel-medialibrary...")
	cmd := exec.Command("composer", "require", "spatie/laravel-medialibrary")
	cmd.Dir = m.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install spatie/laravel-medialibrary: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("📦 Publishing migrations...")
	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Spatie\\MediaLibrary\\MediaLibraryServiceProvider", "--tag=migrations")
	cmd.Dir = m.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to publish media migrations: %v\nOutput: %s", err, string(output))
	}

	return nil
}
