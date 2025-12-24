package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type Architecture struct {
	ProjectPath string
	DryRun      bool
}

func NewArchitecture(projectPath string, dryRun bool) *Architecture {
	return &Architecture{ProjectPath: projectPath, DryRun: dryRun}
}

func (a *Architecture) SetupFolders() error {
	dirs := []string{
		"app/Domain",
		"app/Domain/Shared",
		"app/Domain/Users",
		"app/Domain/Users/Actions",
		"app/Domain/Users/Models",
		"app/Domain/Users/Resources",
		"app/Domain/Users/QueryBuilders",
	}

	for _, dir := range dirs {
		path := filepath.Join(a.ProjectPath, dir)
		if a.DryRun {
			fmt.Printf("[Dry Run] Would create directory: %s\n", path)
			continue
		}
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %v", path, err)
		}
	}

	return nil
}
