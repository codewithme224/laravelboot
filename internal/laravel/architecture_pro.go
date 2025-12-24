package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ProArchSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewProArchSetup(projectPath string, dryRun bool) *ProArchSetup {
	return &ProArchSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (p *ProArchSetup) Setup() error {
	if p.DryRun {
		fmt.Printf("[Dry Run] Would install spatie/laravel-data and setup Action patterns\n")
		return nil
	}

	fmt.Println("🏗️ Installing spatie/laravel-data...")
	cmd := exec.Command("composer", "require", "spatie/laravel-data", "--with-all-dependencies")
	cmd.Dir = p.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install spatie/laravel-data: %v\nOutput: %s", err, string(output))
	}

	if err := p.createBaseAction(); err != nil {
		return err
	}

	return nil
}

func (p *ProArchSetup) createBaseAction() error {
	content := `<?php

namespace App\Support\Actions;

trait AsAction
{
    public static function run(...$arguments)
    {
        return app(static::class)->handle(...$arguments);
    }
}
`
	dir := filepath.Join(p.ProjectPath, "app/Support/Actions")
	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "AsAction.php")
	return os.WriteFile(path, []byte(content), 0644)
}
