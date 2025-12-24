package laravel

import (
	"fmt"
	"os/exec"
)

type Installer struct {
	DryRun bool
}

func NewInstaller(dryRun bool) *Installer {
	return &Installer{DryRun: dryRun}
}

func (i *Installer) CheckDependencies() error {
	deps := []string{"php", "composer"}
	for _, dep := range deps {
		if _, err := exec.LookPath(dep); err != nil {
			return fmt.Errorf("%s is not installed or not in PATH", dep)
		}
	}
	return nil
}

func (i *Installer) HasLaravelInstaller() bool {
	_, err := exec.LookPath("laravel")
	return err == nil
}

func (i *Installer) CreateProject(name string) error {
	var cmd *exec.Cmd
	if i.HasLaravelInstaller() {
		fmt.Printf("Using Laravel installer to create %s...\n", name)
		cmd = exec.Command("laravel", "new", name, "--no-interaction")
	} else {
		fmt.Printf("Laravel installer not found. Using composer to create %s...\n", name)
		cmd = exec.Command("composer", "create-project", "laravel/laravel", name)
	}

	if i.DryRun {
		fmt.Printf("[Dry Run] Would run: %s\n", cmd.String())
		return nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create project: %v\nOutput: %s", err, string(output))
	}

	return nil
}
