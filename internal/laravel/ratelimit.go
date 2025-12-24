package laravel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type RateLimitSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewRateLimitSetup(projectPath string, dryRun bool) *RateLimitSetup {
	return &RateLimitSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (r *RateLimitSetup) Setup() error {
	path := filepath.Join(r.ProjectPath, "app/Providers/AppServiceProvider.php")
	if r.DryRun {
		fmt.Printf("[Dry Run] Would configure Rate Limiting in %s\n", path)
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	sContent := string(content)
	if !strings.Contains(sContent, "RateLimiter::for") {
		importLine := "use Illuminate\\Support\\Facades\\RateLimiter;\nuse Illuminate\\Http\\Request;\nuse Illuminate\\Cache\\RateLimiting\\Limit;"
		// Add imports
		sContent = strings.Replace(sContent, "use Illuminate\\Support\\ServiceProvider;", "use Illuminate\\Support\\ServiceProvider;\n"+importLine, 1)

		// Add rate limiter block in boot method
		rateLimitBlock := `
        RateLimiter::for('api', function (Request $request) {
            return Limit::perMinute(60)->by($request->user()?->id ?: $request->ip());
        });`

		sContent = strings.Replace(sContent, "public function boot(): void\n    {", "public function boot(): void\n    {"+rateLimitBlock, 1)
	}

	return os.WriteFile(path, []byte(sContent), 0644)
}
