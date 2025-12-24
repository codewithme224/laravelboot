package laravel

import (
	"fmt"
	"os/exec"
)

type MonitoringSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewMonitoringSetup(projectPath string, dryRun bool) *MonitoringSetup {
	return &MonitoringSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *MonitoringSetup) Setup() error {
	if m.DryRun {
		fmt.Printf("[Dry Run] Would install Laravel Telescope and Pulse\n")
		return nil
	}

	fmt.Println("🔭 Installing Laravel Telescope...")
	cmd := exec.Command("composer", "require", "laravel/telescope", "--dev", "--with-all-dependencies")
	cmd.Dir = m.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install Telescope: %v\nOutput: %s", err, string(output))
	}

	cmd = exec.Command("php", "artisan", "telescope:install")
	cmd.Dir = m.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize Telescope: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("📈 Installing Laravel Pulse...")
	cmd = exec.Command("composer", "require", "laravel/pulse", "--with-all-dependencies")
	cmd.Dir = m.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install Pulse: %v\nOutput: %s", err, string(output))
	}

	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Laravel\\Pulse\\PulseServiceProvider")
	cmd.Dir = m.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize Pulse: %v\nOutput: %s", err, string(output))
	}

	return nil
}
