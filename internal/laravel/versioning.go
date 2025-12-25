package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type VersioningSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewVersioningSetup(projectPath string, dryRun bool) *VersioningSetup {
	return &VersioningSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (v *VersioningSetup) Setup() error {
	if v.DryRun {
		fmt.Printf("[Dry Run] Would setup API versioning\n")
		return nil
	}

	fmt.Println("🔢 Setting up API versioning...")

	// Create versioned controller directories
	v1Dir := filepath.Join(v.ProjectPath, "app/Http/Controllers/Api/V1")
	v2Dir := filepath.Join(v.ProjectPath, "app/Http/Controllers/Api/V2")
	os.MkdirAll(v1Dir, 0755)
	os.MkdirAll(v2Dir, 0755)

	if err := v.createBaseApiController(); err != nil {
		return err
	}
	if err := v.createV1Controller(); err != nil {
		return err
	}
	if err := v.createV2Controller(); err != nil {
		return err
	}
	if err := v.createVersionedRoutes(); err != nil {
		return err
	}

	return nil
}

func (v *VersioningSetup) createBaseApiController() error {
	content := `<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Traits\ApiResponse;

abstract class BaseApiController extends Controller
{
    use ApiResponse;

    /**
     * Get the API version.
     */
    abstract protected function version(): string;

    /**
     * Get deprecation notice if applicable.
     */
    protected function deprecationNotice(): ?string
    {
        return null;
    }

    /**
     * Add version headers to response.
     */
    protected function withVersionHeaders($response)
    {
        $response->header('X-API-Version', $this->version());
        
        if ($notice = $this->deprecationNotice()) {
            $response->header('X-API-Deprecated', $notice);
        }

        return $response;
    }
}
`
	dir := filepath.Join(v.ProjectPath, "app/Http/Controllers/Api")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "BaseApiController.php"), []byte(content), 0644)
}

func (v *VersioningSetup) createV1Controller() error {
	content := `<?php

namespace App\Http\Controllers\Api\V1;

use App\Http\Controllers\Api\BaseApiController;

abstract class V1Controller extends BaseApiController
{
    protected function version(): string
    {
        return 'v1';
    }

    /**
     * V1 is deprecated, clients should migrate to V2.
     * Remove this method if V1 is still the current version.
     */
    // protected function deprecationNotice(): ?string
    // {
    //     return 'API v1 is deprecated. Please migrate to v2 by 2025-06-01.';
    // }
}
`
	path := filepath.Join(v.ProjectPath, "app/Http/Controllers/Api/V1/V1Controller.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (v *VersioningSetup) createV2Controller() error {
	content := `<?php

namespace App\Http\Controllers\Api\V2;

use App\Http\Controllers\Api\BaseApiController;

abstract class V2Controller extends BaseApiController
{
    protected function version(): string
    {
        return 'v2';
    }
}
`
	path := filepath.Join(v.ProjectPath, "app/Http/Controllers/Api/V2/V2Controller.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (v *VersioningSetup) createVersionedRoutes() error {
	content := `<?php

use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API V1 Routes
|--------------------------------------------------------------------------
|
| Routes for API version 1. These may be deprecated in favor of V2.
|
*/

Route::prefix('v1')->name('v1.')->group(function () {
    // Add your V1 routes here
    // Route::apiResource('users', \App\Http\Controllers\Api\V1\UserController::class);
});

/*
|--------------------------------------------------------------------------
| API V2 Routes
|--------------------------------------------------------------------------
|
| Routes for API version 2 (current version).
|
*/

Route::prefix('v2')->name('v2.')->group(function () {
    // Add your V2 routes here
    // Route::apiResource('users', \App\Http\Controllers\Api\V2\UserController::class);
});

/*
|--------------------------------------------------------------------------
| Latest API Routes (Alias to current version)
|--------------------------------------------------------------------------
|
| Routes without version prefix that point to the latest API version.
|
*/

// Route::apiResource('users', \App\Http\Controllers\Api\V2\UserController::class);
`
	path := filepath.Join(v.ProjectPath, "routes/api_versioned.php")
	return os.WriteFile(path, []byte(content), 0644)
}
