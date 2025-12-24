package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type QualitySetup struct {
	ProjectPath string
	DryRun      bool
}

func NewQualitySetup(projectPath string, dryRun bool) *QualitySetup {
	return &QualitySetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (q *QualitySetup) Setup() error {
	if q.DryRun {
		fmt.Printf("[Dry Run] Would install Pint, Larastan, and Pest\n")
		return nil
	}

	fmt.Println("🧹 Installing Laravel Pint...")
	if err := q.runComposerRequire("laravel/pint", true); err != nil {
		return err
	}

	fmt.Println("🔬 Installing Larastan...")
	if err := q.runComposerRequire("phpstan/phpstan nunomaduro/larastan", true); err != nil {
		return err
	}
	if err := q.createPhpStanConfig(); err != nil {
		return err
	}

	fmt.Println("🧪 Installing Pest...")
	if err := q.runComposerRequire("pestphp/pest pestphp/pest-plugin-laravel", true); err != nil {
		return err
	}

	cmd := exec.Command("php", "artisan", "pest:install", "--no-interaction")
	cmd.Dir = q.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to initialize Pest: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func (q *QualitySetup) runComposerRequire(packages string, dev bool) error {
	args := []string{"require"}
	if dev {
		args = append(args, "--dev")
	}
	// Split packages if multiple are provided
	packageList := strings.Fields(packages)

	args = append(args, packageList...)
	cmd := exec.Command("composer", args...)
	cmd.Dir = q.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install %s: %v\nOutput: %s", packages, err, string(output))
	}
	return nil
}

func (q *QualitySetup) createPhpStanConfig() error {
	content := `includes:
    - ./vendor/nunomaduro/larastan/extension.neon

parameters:
    paths:
        - app/
    level: 5
    ignoreErrors:
    excludePaths:
`
	path := filepath.Join(q.ProjectPath, "phpstan.neon")
	return os.WriteFile(path, []byte(content), 0644)
}
