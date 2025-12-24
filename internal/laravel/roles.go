package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type RolesSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewRolesSetup(projectPath string, dryRun bool) *RolesSetup {
	return &RolesSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (r *RolesSetup) Setup() error {
	if r.DryRun {
		fmt.Printf("[Dry Run] Would install spatie/laravel-permission\n")
		return nil
	}

	fmt.Println("🔑 Installing spatie/laravel-permission...")
	cmd := exec.Command("composer", "require", "spatie/laravel-permission")
	cmd.Dir = r.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install spatie/laravel-permission: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("📦 Publishing configuration...")
	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Spatie\\Permission\\PermissionServiceProvider")
	cmd.Dir = r.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to publish permission config: %v\nOutput: %s", err, string(output))
	}

	if err := r.modifyUserModel(); err != nil {
		return err
	}

	return nil
}

func (r *RolesSetup) modifyUserModel() error {
	path := filepath.Join(r.ProjectPath, "app/Models/User.php")
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	sContent := string(content)
	if !strings.Contains(sContent, "Spatie\\Permission\\Traits\\HasRoles") {
		sContent = strings.Replace(sContent, "use HasApiTokens, HasFactory, Notifiable;", "use HasApiTokens, HasFactory, Notifiable, HasRoles;", 1)
		sContent = strings.Replace(sContent, "use Laravel\\Sanctum\\HasApiTokens;", "use Laravel\\Sanctum\\HasApiTokens;\nuse Spatie\\Permission\\Traits\\HasRoles;", 1)
	}

	return os.WriteFile(path, []byte(sContent), 0644)
}
