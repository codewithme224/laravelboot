package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type ActivityLogSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewActivityLogSetup(projectPath string, dryRun bool) *ActivityLogSetup {
	return &ActivityLogSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (a *ActivityLogSetup) Setup() error {
	if a.DryRun {
		fmt.Printf("[Dry Run] Would install spatie/laravel-activitylog\n")
		return nil
	}

	fmt.Println("📝 Installing spatie/laravel-activitylog...")
	cmd := exec.Command("composer", "require", "spatie/laravel-activitylog", "--with-all-dependencies")
	cmd.Dir = a.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install spatie/laravel-activitylog: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("📦 Publishing migrations...")
	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Spatie\\Activitylog\\ActivitylogServiceProvider", "--tag=activitylog-migrations")
	cmd.Dir = a.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to publish activitylog migrations: %v\nOutput: %s", err, string(output))
	}

	if err := a.setupTrait(); err != nil {
		return err
	}

	return nil
}

func (a *ActivityLogSetup) setupTrait() error {
	content := `<?php

namespace App\Support\Concerns;

use Spatie\Activitylog\LogOptions;
use Spatie\Activitylog\Traits\LogsActivity;

trait InteractsWithActivityLog
{
    use LogsActivity;

    public function getActivitylogOptions(): LogOptions
    {
        return LogOptions::defaults()
            ->logAll()
            ->logOnlyDirty()
            ->useLogName(str(class_basename($this))->plural()->lower());
    }
}
`
	dir := filepath.Join(a.ProjectPath, "app/Support/Concerns")
	if !a.DryRun {
		os.MkdirAll(dir, 0755)
		path := filepath.Join(dir, "InteractsWithActivityLog.php")
		return os.WriteFile(path, []byte(content), 0644)
	}
	return nil
}
