package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type TraitsSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewTraitsSetup(projectPath string, dryRun bool) *TraitsSetup {
	return &TraitsSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (t *TraitsSetup) Setup() error {
	if t.DryRun {
		fmt.Printf("[Dry Run] Would create common API traits\n")
		return nil
	}

	fmt.Println("🧬 Creating common traits...")

	if err := t.createApiTrait(); err != nil {
		return err
	}
	if err := t.createHandlesPaginationTrait(); err != nil {
		return err
	}
	if err := t.createAuditableTrait(); err != nil {
		return err
	}

	return nil
}

func (t *TraitsSetup) createApiTrait() error {
	content := `<?php

namespace App\Traits;

use Illuminate\Support\Facades\Http;
use Illuminate\Http\Client\Response;

trait Api
{
    protected function get(string $url, array $query = [], array $headers = []): Response
    {
        return Http::withHeaders($this->mergeHeaders($headers))->get($url, $query);
    }

    protected function post(string $url, array $data = [], array $headers = []): Response
    {
        return Http::withHeaders($this->mergeHeaders($headers))->post($url, $data);
    }

    protected function put(string $url, array $data = [], array $headers = []): Response
    {
        return Http::withHeaders($this->mergeHeaders($headers))->put($url, $data);
    }

    protected function delete(string $url, array $data = [], array $headers = []): Response
    {
        return Http::withHeaders($this->mergeHeaders($headers))->delete($url, $data);
    }

    protected function mergeHeaders(array $headers = []): array
    {
        return array_merge([
            'Accept' => 'application/json',
            'Authorization' => request()->header('Authorization'),
        ], $headers);
    }
}
`
	dir := filepath.Join(t.ProjectPath, "app/Traits")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "Api.php"), []byte(content), 0644)
}

func (t *TraitsSetup) createHandlesPaginationTrait() error {
	content := `<?php

namespace App\Traits;

use Illuminate\Database\Eloquent\Builder;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\ResourceCollection;

trait HandlesPagination
{
    /**
     * Handle pagination for queries
     */
    public function handlePagination(Builder $query, Request $request, string $collectionClass): ResourceCollection
    {
        $searchQuery = $request->query('search');

        if ($searchQuery && method_exists($query->getModel(), 'scopeSearch')) {
            $query->search($searchQuery);
        }

        if ($request->query('paginate', true) !== 'false') {
            $data = $query->paginate($request->query('per_page', 15));
        } else {
            $data = $query->get();
        }

        return new $collectionClass($data);
    }

    /**
     * Paginate a collection
     */
    public function paginateCollection(\Illuminate\Support\Collection $collection, Request $request, string $collectionClass): ResourceCollection
    {
        $page = (int) $request->query('page', 1);
        $perPage = (int) $request->query('per_page', 15);

        $items = $collection->slice(($page - 1) * $perPage, $perPage)->values();

        $paginated = new \Illuminate\Pagination\LengthAwarePaginator(
            $items,
            $collection->count(),
            $perPage,
            $page,
            ['path' => $request->url(), 'query' => $request->query()]
        );

        return new $collectionClass($paginated);
    }
}
`
	dir := filepath.Join(t.ProjectPath, "app/Traits")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "HandlesPagination.php"), []byte(content), 0644)
}

func (t *TraitsSetup) createAuditableTrait() error {
	content := `<?php

namespace App\Traits;

trait Auditable
{
    public static function bootAuditable()
    {
        static::created(function ($model) {
            activity()
                ->performedOn($model)
                ->causedBy(auth()->user())
                ->withProperties(['attributes' => $model->getAttributes()])
                ->log('created');
        });

        static::updated(function ($model) {
            activity()
                ->performedOn($model)
                ->causedBy(auth()->user())
                ->withProperties([
                    'old' => $model->getOriginal(),
                    'new' => $model->getAttributes(),
                ])
                ->log('updated');
        });

        static::deleted(function ($model) {
            activity()
                ->performedOn($model)
                ->causedBy(auth()->user())
                ->log('deleted');
        });
    }
}
`
	dir := filepath.Join(t.ProjectPath, "app/Traits")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "Auditable.php"), []byte(content), 0644)
}
