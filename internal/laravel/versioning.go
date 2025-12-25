package laravel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
		fmt.Printf("[Dry Run] Would setup API versioning app structure\n")
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
	if err := v.updateBootstrapApp(); err != nil {
		return err
	}
	if err := v.updateApiRoutes(); err != nil {
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

func (v *VersioningSetup) updateBootstrapApp() error {
	path := filepath.Join(v.ProjectPath, "bootstrap/app.php")
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	text := string(content)
	if !strings.Contains(text, "apiPrefix:") {
		// Look for the api route definition
		target := "api: __DIR__ . '/../routes/api.php',"
		replacement := "api: __DIR__ . '/../routes/api.php',\n        apiPrefix: 'api/v1',"

		newText := strings.Replace(text, target, replacement, 1)
		if newText == text {
			// Try without comma just in case
			target = "api: __DIR__ . '/../routes/api.php'"
			replacement = "api: __DIR__ . '/../routes/api.php',\n        apiPrefix: 'api/v1'"
			newText = strings.Replace(text, target, replacement, 1)
		}

		if newText != text {
			return os.WriteFile(path, []byte(newText), 0644)
		}
		fmt.Println("⚠️ Could not automatically inject apiPrefix in bootstrap/app.php")
	}

	return nil
}

func (v *VersioningSetup) updateApiRoutes() error {
	content := `<?php

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider and all of them will
| be assigned to the "api" middleware group. Make something great!
|
*/

// V1 Routes (Pre-fixed with api/v1 in bootstrap/app.php)
Route::middleware('auth:sanctum')->group(function () {
    Route::get('/user', function (Request $request) {
        return $request->user();
    });
    
    // Add your V1 routes here
    // Route::apiResource('users', \App\Http\Controllers\Api\V1\UserController::class);
});

// Health check route
Route::get('/health', function () {
    return response()->json(['status' => 'ok']);
});
`
	path := filepath.Join(v.ProjectPath, "routes/api.php")
	return os.WriteFile(path, []byte(content), 0644)
}
