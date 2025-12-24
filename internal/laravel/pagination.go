package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type PaginationSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewPaginationSetup(projectPath string, dryRun bool) *PaginationSetup {
	return &PaginationSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (p *PaginationSetup) Setup() error {
	if err := p.createApiResponseSupport(); err != nil {
		return err
	}
	if err := p.createQuerySupport(); err != nil {
		return err
	}
	return nil
}

func (p *PaginationSetup) createApiResponseSupport() error {
	content := `<?php

namespace App\Support\Api;

use Illuminate\Http\JsonResponse;
use Illuminate\Pagination\LengthAwarePaginator;

trait ApiResponse
{
    public function ok($data, string $message = 'Success'): JsonResponse
    {
        return response()->json([
            'success' => true,
            'message' => $message,
            'data' => $data,
        ]);
    }

    public function created($data, string $message = 'Resource created successfully'): JsonResponse
    {
        return response()->json([
            'success' => true,
            'message' => $message,
            'data' => $data,
        ], 201);
    }

    public function deleted(string $message = 'Resource deleted successfully'): JsonResponse
    {
        return response()->json([
            'success' => true,
            'message' => $message,
            'data' => null,
        ], 200);
    }

    public function paginate(LengthAwarePaginator $paginator, string $message = 'Success'): JsonResponse
    {
        return response()->json([
            'success' => true,
            'message' => $message,
            'data' => $paginator->items(),
            'meta' => [
                'current_page' => $paginator->currentPage(),
                'last_page' => $paginator->lastPage(),
                'per_page' => $paginator->perPage(),
                'total' => $paginator->total(),
            ],
        ]);
    }

    public function error(string $message = 'Error', int $code = 400, array $errors = []): JsonResponse
    {
        return response()->json([
            'success' => false,
            'message' => $message,
            'errors' => $errors,
        ], $code);
    }

    public function unauthorized(string $message = 'Unauthorized', array $errors = []): JsonResponse
    {
        return $this->error($message, 401, $errors);
    }

    public function unauthenticated(string $message = 'Unauthenticated', array $errors = []): JsonResponse
    {
        return $this->error($message, 401, $errors);
    }

    public function forbidden(string $message = 'Forbidden', array $errors = []): JsonResponse
    {
        return $this->error($message, 403, $errors);
    }
}
`
	dir := filepath.Join(p.ProjectPath, "app/Support/Api")
	if !p.DryRun {
		os.MkdirAll(dir, 0755)
	}
	path := filepath.Join(dir, "ApiResponse.php")
	if p.DryRun {
		fmt.Printf("[Dry Run] Would create support file: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (p *PaginationSetup) createQuerySupport() error {
	content := `<?php

namespace App\Support\Query;

use Spatie\QueryBuilder\QueryBuilder;
use Illuminate\Database\Eloquent\Builder;

trait AppliesQueryBuilder
{
    /**
     * @param Builder|string $subject
     * @param array $allowedFilters
     * @param array $allowedSorts
     * @return QueryBuilder
     */
    protected function buildQuery($subject, array $allowedFilters = [], array $allowedSorts = []): QueryBuilder
    {
        return QueryBuilder::for($subject)
            ->allowedFilters($allowedFilters)
            ->allowedSorts($allowedSorts);
    }
}
`
	dir := filepath.Join(p.ProjectPath, "app/Support/Query")
	if !p.DryRun {
		os.MkdirAll(dir, 0755)
	}
	path := filepath.Join(dir, "AppliesQueryBuilder.php")
	if p.DryRun {
		fmt.Printf("[Dry Run] Would create support file: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}
