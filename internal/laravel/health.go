package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type HealthSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewHealthSetup(projectPath string, dryRun bool) *HealthSetup {
	return &HealthSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (h *HealthSetup) Setup() error {
	if err := h.createHealthController(); err != nil {
		return err
	}
	if err := h.registerRoute(); err != nil {
		return err
	}
	return nil
}

func (h *HealthSetup) createHealthController() error {
	content := `<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use Illuminate\Http\JsonResponse;
use Illuminate\Support\Facades\DB;

class HealthController extends Controller
{
    public function check(): JsonResponse
    {
        try {
            DB::connection()->getPdo();
            return response()->json([
                'status' => 'ok',
                'database' => 'connected',
                'timestamp' => now()->toIso8601String(),
            ]);
        } catch (\Exception $e) {
            return response()->json([
                'status' => 'error',
                'database' => 'disconnected',
                'message' => $e->getMessage(),
            ], 503);
        }
    }
}
`
	path := filepath.Join(h.ProjectPath, "app/Http/Controllers/Api/HealthController.php")
	if h.DryRun {
		fmt.Printf("[Dry Run] Would create HealthController: %s\n", path)
		return nil
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", dir, err)
	}

	return os.WriteFile(path, []byte(content), 0644)
}

func (h *HealthSetup) registerRoute() error {
	path := filepath.Join(h.ProjectPath, "routes/api.php")
	if h.DryRun {
		fmt.Printf("[Dry Run] Would add health route to %s\n", path)
		return nil
	}

	// Check if api.php exists (Laravel 11+ might not have it by default)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Println("📍 routes/api.php not found. Running php artisan install:api...")
		cmd := exec.Command("php", "artisan", "install:api", "--no-interaction")
		cmd.Dir = h.ProjectPath
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to run install:api: %v\nOutput: %s", err, string(output))
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	route := "\nRoute::get('/health', [\\App\\Http\\Controllers\\Api\\HealthController::class, 'check']);\n"
	if !strings.Contains(string(content), "/health") {
		newContent := string(content) + route
		return os.WriteFile(path, []byte(newContent), 0644)
	}
	return nil
}
