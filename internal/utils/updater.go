package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
}

func CheckForUpdate(currentVersion string) {
	fmt.Printf("🔍 Checking for updates (current: %s)...\n", currentVersion)

	resp, err := http.Get("https://api.github.com/repos/codewithme224/laravelboot/releases/latest")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return
	}

	if release.TagName != "" && release.TagName != currentVersion {
		fmt.Printf("\n✨ A new version is available: %s\n", release.TagName)
		fmt.Println("👉 Run 'laravelboot update' to upgrade now!")
	}
}

func SelfUpdate() error {
	fmt.Println("🚀 Starting self-update...")

	if runtime.GOOS == "windows" {
		return fmt.Errorf("self-update is not supported on Windows. Please download the latest release manually")
	}

	// We use the install script we already created
	cmd := exec.Command("bash", "-c", "curl -sL https://raw.githubusercontent.com/codewithme224/laravelboot/main/install.sh | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run update script: %v", err)
	}

	fmt.Println("\n✅ LaravelBoot has been updated successfully!")
	return nil
}
