package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type CicdSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewCicdSetup(projectPath string, dryRun bool) *CicdSetup {
	return &CicdSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (c *CicdSetup) Setup() error {
	content := `name: CI

on:
  push:
    branches: [ main, master ]
  pull_request:
    branches: [ main, master ]

jobs:
  laravel-tests:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    - name: Setup PHP
      uses: shivammathur/setup-php@v2
      with:
        php-version: '8.3'
        extensions: mbstring, dom, curl, libxml, mysql, pdo_mysql
        coverage: xdebug
    - name: Install Dependencies
      run: composer install -q --no-ansi --no-interaction --no-scripts --no-progress --prefer-dist
    - name: Copy .env
      run: php -r "file_exists('.env') || copy('.env.example', '.env');"
    - name: Generate key
      run: php artisan key:generate
    - name: Directory Permissions
      run: chmod -R 777 storage bootstrap/cache
    - name: Run Tests
      run: php artisan test
    - name: Check Code Style (Pint)
      run: ./vendor/bin/pint --test
    - name: Static Analysis (Larastan)
      run: ./vendor/bin/phpstan analyse
`
	dir := filepath.Join(c.ProjectPath, ".github/workflows")
	if c.DryRun {
		fmt.Printf("[Dry Run] Would create GitHub Actions workflow: %s\n", filepath.Join(dir, "ci.yml"))
		return nil
	}

	os.MkdirAll(dir, 0755)
	path := filepath.Join(dir, "ci.yml")
	return os.WriteFile(path, []byte(content), 0644)
}
