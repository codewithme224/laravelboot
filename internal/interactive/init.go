package interactive

import (
	"bufio"
	"fmt"
	"laravelboot/internal/config"
	"os"
	"strings"
)

func RunInit() (*config.Config, error) {
	reader := bufio.NewReader(os.Stdin)
	conf := config.DefaultConfig()

	fmt.Println("🚀 Welcome to LaravelBoot Interactive Setup")
	fmt.Println("-------------------------------------------")

	fmt.Printf("Project Name [%s]: ", conf.ProjectName)
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	if name != "" {
		conf.ProjectName = name
	}

	fmt.Printf("Database (mysql, postgres, sqlite) [%s]: ", conf.Database)
	db, _ := reader.ReadString('\n')
	db = strings.TrimSpace(db)
	if db != "" {
		conf.Database = db
	}

	fmt.Printf("Auth (sanctum, passport) [%s]: ", conf.Auth)
	auth, _ := reader.ReadString('\n')
	auth = strings.TrimSpace(auth)
	if auth != "" {
		conf.Auth = auth
	}

	fmt.Printf("Enable Enterprise Stack (Quality, CI, Monitoring, Pro-Arch) (y/n) [n]: ")
	ent, _ := reader.ReadString('\n')
	ent = strings.TrimSpace(ent)
	if strings.ToLower(ent) == "y" {
		conf.Enterprise = []string{"quality", "pro-arch", "docs-pro", "ci", "monitoring"}
	}

	fmt.Println("\n✅ Configuration generated!")
	return conf, nil
}
