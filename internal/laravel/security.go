package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type SecuritySetup struct {
	ProjectPath string
	DryRun      bool
}

func NewSecuritySetup(projectPath string, dryRun bool) *SecuritySetup {
	return &SecuritySetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *SecuritySetup) Setup() error {
	if err := s.createForceJsonResponseMiddleware(); err != nil {
		return err
	}
	if err := s.createEnvValidator(); err != nil {
		return err
	}
	return nil
}

func (s *SecuritySetup) createForceJsonResponseMiddleware() error {
	content := `<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Symfony\Component\HttpFoundation\Response;

class ForceJsonResponse
{
    public function handle(Request $request, Closure $next): Response
    {
        $request->headers->set('Accept', 'application/json');

        return $next($request);
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Http/Middleware")
	if !s.DryRun {
		os.MkdirAll(dir, 0755)
	}
	path := filepath.Join(dir, "ForceJsonResponse.php")
	if s.DryRun {
		fmt.Printf("[Dry Run] Would create middleware: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *SecuritySetup) createEnvValidator() error {
	content := `<?php

namespace App\Support\Env;

use Illuminate\Support\Facades\App;
use RuntimeException;

class EnvValidator
{
    public static function validate(): void
    {
        $requiredEnv = [
            'APP_KEY',
            'DB_HOST',
            'DB_USERNAME',
            'DB_PASSWORD',
        ];

        foreach ($requiredEnv as $env) {
            if (empty(env($env))) {
                throw new RuntimeException("Missing required environment variable: {$env}");
            }
        }
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Support/Env")
	if !s.DryRun {
		os.MkdirAll(dir, 0755)
	}
	path := filepath.Join(dir, "EnvValidator.php")
	if s.DryRun {
		fmt.Printf("[Dry Run] Would create environment validator: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}
