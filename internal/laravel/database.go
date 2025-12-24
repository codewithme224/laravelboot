package laravel

import (
	"fmt"
	"os/exec"
)

type DatabaseSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewDatabaseSetup(projectPath string, dryRun bool) *DatabaseSetup {
	return &DatabaseSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (d *DatabaseSetup) RunMigrations() error {
	if d.DryRun {
		fmt.Printf("[Dry Run] Would run: php artisan migrate\n")
		return nil
	}

	cmd := exec.Command("php", "artisan", "migrate", "--force")
	cmd.Dir = d.ProjectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run migrations: %v\nOutput: %s", err, string(output))
	}

	return nil
}
