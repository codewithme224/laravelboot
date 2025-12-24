package laravel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ApiSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewApiSetup(projectPath string, dryRun bool) *ApiSetup {
	return &ApiSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *ApiSetup) Configure() error {
	if err := s.createResponseServiceProvider(); err != nil {
		return err
	}
	if err := s.registerServiceProvider(); err != nil {
		return err
	}
	if err := s.forceJsonResponse(); err != nil {
		return err
	}
	return nil
}

func (s *ApiSetup) createResponseServiceProvider() error {
	content := `<?php

namespace App\Providers;

use Illuminate\Support\ServiceProvider;
use Illuminate\Support\Facades\Response;

class ApiResponseServiceProvider extends ServiceProvider
{
    public function boot(): void
    {
        Response::macro('success', function ($data, $message = 'Success', $code = 200) {
            return Response::json([
                'success' => true,
                'message' => $message,
                'data' => $data,
            ], $code);
        });

        Response::macro('created', function ($data, $message = 'Resource created successfully') {
            return Response::json([
                'success' => true,
                'message' => $message,
                'data' => $data,
            ], 201);
        });

        Response::macro('deleted', function ($message = 'Resource deleted successfully') {
            return Response::json([
                'success' => true,
                'message' => $message,
                'data' => null,
            ], 200);
        });

        Response::macro('error', function ($message = 'Error', $code = 400, $errors = []) {
            return Response::json([
                'success' => false,
                'message' => $message,
                'errors' => $errors,
            ], $code);
        });

        Response::macro('unauthorized', function ($message = 'Unauthorized', $errors = []) {
            return Response::error($message, 401, $errors);
        });

        Response::macro('unauthenticated', function ($message = 'Unauthenticated', $errors = []) {
            return Response::error($message, 401, $errors);
        });

        Response::macro('forbidden', function ($message = 'Forbidden', $errors = []) {
            return Response::error($message, 403, $errors);
        });
    }
}
`
	path := filepath.Join(s.ProjectPath, "app/Providers/ApiResponseServiceProvider.php")
	if s.DryRun {
		fmt.Printf("[Dry Run] Would create file: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *ApiSetup) registerServiceProvider() error {
	path := filepath.Join(s.ProjectPath, "bootstrap/providers.php")
	if s.DryRun {
		fmt.Printf("[Dry Run] Would register ApiResponseServiceProvider in %s\n", path)
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	newContent := strings.Replace(
		string(content),
		"];",
		"    App\\Providers\\ApiResponseServiceProvider::class,\n];",
		1,
	)

	return os.WriteFile(path, []byte(newContent), 0644)
}

func (s *ApiSetup) forceJsonResponse() error {
	path := filepath.Join(s.ProjectPath, "bootstrap/app.php")
	if s.DryRun {
		fmt.Printf("[Dry Run] Would check %s to ensure it exists\n", path)
		return nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("bootstrap/app.php not found at %s", path)
	}

	return nil
}
