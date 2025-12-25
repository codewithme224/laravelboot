package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type MiddlewareSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewMiddlewareSetup(projectPath string, dryRun bool) *MiddlewareSetup {
	return &MiddlewareSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *MiddlewareSetup) Setup() error {
	if m.DryRun {
		fmt.Printf("[Dry Run] Would create common middleware\n")
		return nil
	}

	fmt.Println("🛡️ Creating common middleware...")

	if err := m.createDBTransactionMiddleware(); err != nil {
		return err
	}
	if err := m.createForceJsonMiddleware(); err != nil {
		return err
	}

	return nil
}

func (m *MiddlewareSetup) createDBTransactionMiddleware() error {
	content := `<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Symfony\Component\HttpFoundation\Response;

class DBTransaction
{
    /**
     * Handle an incoming request.
     * Wraps the request in a database transaction.
     */
    public function handle(Request $request, Closure $next): Response
    {
        DB::beginTransaction();

        $response = $next($request);

        if ($response->getStatusCode() > 399) {
            DB::rollBack();
        } else {
            DB::commit();
        }

        return $response;
    }
}
`
	dir := filepath.Join(m.ProjectPath, "app/Http/Middleware")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "DBTransaction.php"), []byte(content), 0644)
}

func (m *MiddlewareSetup) createForceJsonMiddleware() error {
	content := `<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class ForceJson
{
    /**
     * Force JSON responses for API requests.
     */
    public function handle(Request $request, Closure $next): Response
    {
        $request->headers->set('Accept', 'application/json');
        
        return $next($request);
    }
}
`
	dir := filepath.Join(m.ProjectPath, "app/Http/Middleware")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "ForceJson.php"), []byte(content), 0644)
}
