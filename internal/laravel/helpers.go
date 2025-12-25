package laravel

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type HelpersSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewHelpersSetup(projectPath string, dryRun bool) *HelpersSetup {
	return &HelpersSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (h *HelpersSetup) Setup() error {
	if h.DryRun {
		fmt.Printf("[Dry Run] Would create app/helpers.php and register it in composer.json\n")
		return nil
	}

	fmt.Println("🤝 Setting up global helpers...")

	// 1. Create app/helpers.php
	content := `<?php

if (!function_exists('user')) {
    /**
     * Get the authenticated user.
     */
    function user(): ?\App\Models\User
    {
        return auth()->user();
    }
}

if (!function_exists('api_response')) {
    /**
     * Standardized API response.
     */
    function api_response(string $message, mixed $data = null, int $status = 200)
    {
        return response()->json([
            'message' => $message,
            'data' => $data,
        ], $status);
    }
}
`
	helperPath := filepath.Join(h.ProjectPath, "app/helpers.php")
	if err := os.WriteFile(helperPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to create app/helpers.php: %v", err)
	}

	// 2. Register in composer.json
	composerPath := filepath.Join(h.ProjectPath, "composer.json")
	data, err := os.ReadFile(composerPath)
	if err != nil {
		return fmt.Errorf("failed to read composer.json: %v", err)
	}

	var composer map[string]interface{}
	if err := json.Unmarshal(data, &composer); err != nil {
		return fmt.Errorf("failed to parse composer.json: %v", err)
	}

	autoload, ok := composer["autoload"].(map[string]interface{})
	if !ok {
		autoload = make(map[string]interface{})
		composer["autoload"] = autoload
	}

	files, ok := autoload["files"].([]interface{})
	if !ok {
		files = []interface{}{}
	}

	// Check if already registered
	exists := false
	for _, f := range files {
		if f == "app/helpers.php" {
			exists = true
			break
		}
	}

	if !exists {
		files = append(files, "app/helpers.php")
		autoload["files"] = files

		newData, err := json.MarshalIndent(composer, "", "    ")
		if err != nil {
			return fmt.Errorf("failed to marshal composer.json: %v", err)
		}

		if err := os.WriteFile(composerPath, newData, 0644); err != nil {
			return fmt.Errorf("failed to update composer.json: %v", err)
		}
	}

	return nil
}
