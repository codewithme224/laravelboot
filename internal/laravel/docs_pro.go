package laravel

import (
	"fmt"
	"os/exec"
)

type DocsProSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewDocsProSetup(projectPath string, dryRun bool) *DocsProSetup {
	return &DocsProSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (d *DocsProSetup) Setup() error {
	if d.DryRun {
		fmt.Printf("[Dry Run] Would install dedoc/scramble for automated API docs\n")
		return nil
	}

	fmt.Println("📚 Installing dedoc/scramble...")
	cmd := exec.Command("composer", "require", "dedoc/scramble", "--with-all-dependencies")
	cmd.Dir = d.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install Scramble: %v\nOutput: %s", err, string(output))
	}

	return nil
}
