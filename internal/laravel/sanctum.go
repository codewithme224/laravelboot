package laravel

import (
	"fmt"
	"os/exec"
)

type SanctumInstaller struct {
	ProjectPath string
	DryRun      bool
}

func NewSanctumInstaller(projectPath string, dryRun bool) *SanctumInstaller {
	return &SanctumInstaller{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *SanctumInstaller) Install() error {
	if s.DryRun {
		fmt.Printf("[Dry Run] Would run: php artisan install:api\n")
		return nil
	}

	cmd := exec.Command("php", "artisan", "install:api", "--no-interaction")
	cmd.Dir = s.ProjectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install API (Sanctum): %v\nOutput: %s", err, string(output))
	}

	return nil
}
