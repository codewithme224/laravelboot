package laravel

import (
	"fmt"
	"os/exec"
)

type TenancySetup struct {
	ProjectPath string
	DryRun      bool
}

func NewTenancySetup(projectPath string, dryRun bool) *TenancySetup {
	return &TenancySetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (t *TenancySetup) Setup() error {
	if t.DryRun {
		fmt.Printf("[Dry Run] Would install Multi-Tenancy support (stancl/tenancy)\n")
		return nil
	}

	fmt.Println("🏢 Installing stancl/tenancy...")
	cmd := exec.Command("composer", "require", "stancl/tenancy", "--with-all-dependencies")
	cmd.Dir = t.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install stancl/tenancy: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("⚙️ Initializing tenancy...")
	initCmd := exec.Command("php", "artisan", "tenancy:install")
	initCmd.Dir = t.ProjectPath
	if output, err := initCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize tenancy: %v\nOutput: %s", err, string(output))
	}

	return nil
}
