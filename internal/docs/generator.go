package docs

import (
	"fmt"
	"laravelboot/internal/config"
	"os"
	"path/filepath"
)

func Generate(projectPath string, conf *config.Config, dryRun bool) error {
	content := fmt.Sprintf(`# %s - API Documentation

## Overview
This project was generated using LaravelBoot.

## Architecture
- **Structure**: %s
- **Database**: %s
- **Auth**: %s

## Enabled Features
%v

## Infrastructure
%v

## API Endpoints
- POST /api/login
- POST /api/logout
- GET /api/me
- GET /api/health
`, conf.ProjectName, conf.Architecture, conf.Database, conf.Auth, conf.Features, conf.Infra)

	path := filepath.Join(projectPath, "README-API.md")
	if dryRun {
		fmt.Printf("[Dry Run] Would generate documentation: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}
