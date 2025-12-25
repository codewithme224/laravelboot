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
		fmt.Printf("[Dry Run] Would install spatie/laravel-query-builder with config and service\n")
		return nil
	}

	fmt.Println("🔍 Installing spatie/laravel-query-builder...")
	cmd := exec.Command("composer", "require", "spatie/laravel-query-builder", "--with-all-dependencies")
	cmd.Dir = s.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install spatie/laravel-query-builder: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("⚙️ Publishing config...")
	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Spatie\\QueryBuilder\\QueryBuilderServiceProvider", "--tag=query-builder-config")
	cmd.Dir = s.ProjectPath
	_ = cmd.Run()

	fmt.Println("🔧 Creating QueryBuilderService...")
	if err := s.createQueryBuilderService(); err != nil {
		return err
	}

	fmt.Println("🔧 Creating example QueryBuilder...")
	if err := s.CreateExample(); err != nil {
		return err
	}

	return nil
}

func (s *SpatieQueryBuilder) createQueryBuilderService() error {
	content := `<?php

namespace App\Services;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Http\Request;
use Spatie\QueryBuilder\QueryBuilder;
use Spatie\QueryBuilder\AllowedFilter;
use Spatie\QueryBuilder\AllowedSort;

class QueryBuilderService
{
    /**
     * Create a QueryBuilder instance for the given model.
     */
    public static function for(string $modelClass, ?Request $request = null): QueryBuilder
    {
        return QueryBuilder::for($modelClass, $request ?? request());
    }

    /**
     * Apply common filters to a query builder.
     */
    public static function withCommonFilters(QueryBuilder $query, array $searchableFields = []): QueryBuilder
    {
        $filters = [
            AllowedFilter::exact('id'),
            AllowedFilter::partial('created_at'),
        ];

        foreach ($searchableFields as $field) {
            $filters[] = AllowedFilter::partial($field);
        }

        return $query->allowedFilters($filters);
    }

    /**
     * Apply common sorts to a query builder.
     */
    public static function withCommonSorts(QueryBuilder $query, array $additionalSorts = []): QueryBuilder
    {
        $sorts = array_merge(['id', 'created_at', 'updated_at'], $additionalSorts);
        
        return $query->allowedSorts($sorts)->defaultSort('-created_at');
    }

    /**
     * Apply pagination to a query builder.
     */
    public static function paginate(QueryBuilder $query, ?int $perPage = null)
    {
        $perPage = $perPage ?? request()->input('per_page', 15);
        
        if (request()->input('paginate', true) === 'false') {
            return $query->get();
        }

        return $query->paginate($perPage);
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Services")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "QueryBuilderService.php"), []byte(content), 0644)
}

func (s *SpatieQueryBuilder) CreateExample() error {
	content := `<?php

namespace App\Domain\Users\QueryBuilders;

use App\Models\User;
use Spatie\QueryBuilder\QueryBuilder;
use Spatie\QueryBuilder\AllowedFilter;

class UserQueryBuilder extends QueryBuilder
{
    public function __construct()
    {
        parent::__construct(User::query());

        $this->allowedFilters([
            AllowedFilter::partial('name'),
            AllowedFilter::exact('email'),
            AllowedFilter::exact('id'),
        ])
        ->allowedSorts(['name', 'email', 'created_at'])
        ->allowedIncludes(['posts', 'roles'])
        ->defaultSort('-created_at');
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Domain/Users/QueryBuilders")
	os.MkdirAll(dir, 0755)
	if s.DryRun {
		fmt.Printf("[Dry Run] Would create example QueryBuilder: %s\n", filepath.Join(dir, "UserQueryBuilder.php"))
		return nil
	}

	return os.WriteFile(filepath.Join(dir, "UserQueryBuilder.php"), []byte(content), 0644)
}
