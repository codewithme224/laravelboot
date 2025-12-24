package laravel

import (
	"fmt"
	"os/exec"
)

type SearchSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewSearchSetup(projectPath string, dryRun bool) *SearchSetup {
	return &SearchSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *SearchSetup) Setup() error {
	if s.DryRun {
		fmt.Printf("[Dry Run] Would install laravel/scout and typesense/typesense-php\n")
		return nil
	}

	fmt.Println("🔍 Installing laravel/scout...")
	cmd := exec.Command("composer", "require", "laravel/scout", "--with-all-dependencies")
	cmd.Dir = s.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install laravel/scout: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("🔍 Installing typesense/typesense-php and dev-it-me/laravel-scout-typesense-driver...")
	cmd = exec.Command("composer", "require", "typesense/typesense-php", "dev-it-me/laravel-scout-typesense-driver", "--with-all-dependencies")
	cmd.Dir = s.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install typesense driver: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("📦 Publishing scout configuration...")
	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Laravel\\Scout\\ScoutServiceProvider")
	cmd.Dir = s.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to publish scout config: %v\nOutput: %s", err, string(output))
	}

	return nil
}
