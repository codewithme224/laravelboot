package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type SpatieQueryBuilder struct {
	ProjectPath string
	DryRun      bool
}

func NewSpatieQueryBuilder(projectPath string, dryRun bool) *SpatieQueryBuilder {
	return &SpatieQueryBuilder{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *SpatieQueryBuilder) Install() error {
	if s.DryRun {
		fmt.Printf("[Dry Run] Would run: composer require spatie/laravel-query-builder\n")
		return nil
	}

	cmd := exec.Command("composer", "require", "spatie/laravel-query-builder")
	cmd.Dir = s.ProjectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install spatie/laravel-query-builder: %v\nOutput: %s", err, string(output))
	}

	return nil
}

func (s *SpatieQueryBuilder) CreateExample() error {
	content := `<?php

namespace App\Domain\Users\QueryBuilders;

use App\Models\User;
use Spatie\QueryBuilder\QueryBuilder;

class UserQueryBuilder extends QueryBuilder
{
    public function __construct()
    {
        parent::__construct(User::query());

        $this->allowedFilters(['name', 'email'])
            ->allowedSorts(['name', 'created_at']);
    }
}
`
	path := filepath.Join(s.ProjectPath, "app/Domain/Users/QueryBuilders/UserQueryBuilder.php")
	if s.DryRun {
		fmt.Printf("[Dry Run] Would create example QueryBuilder: %s\n", path)
		return nil
	}

	return os.WriteFile(path, []byte(content), 0644)
}
