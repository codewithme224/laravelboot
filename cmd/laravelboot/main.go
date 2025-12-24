package main

import (
	"fmt"
	"laravelboot/internal/interactive"
	"laravelboot/internal/laravel"
	"os"
)

func main() {
	var dryRun bool
	var args []string

	for _, arg := range os.Args[1:] {
		if arg == "--dry-run" {
			dryRun = true
		} else {
			args = append(args, arg)
		}
	}

	if len(args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]
	target := ""
	if len(args) > 1 {
		target = args[1]
	}

	switch command {
	case "init":
		conf, err := interactive.RunInit()
		if err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}
		if err := conf.Save(".laravelboot.yaml"); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error saving config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✨ Configuration saved to .laravelboot.yaml")

	case "new":
		if target == "" {
			printUsage()
			os.Exit(1)
		}
		appName := target
		creator := laravel.NewCreator(appName, dryRun)
		if err := creator.Create(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "add":
		if target == "" {
			printUsage()
			os.Exit(1)
		}
		switch target {
		case "auth":
			manager := laravel.NewAuthManager(dryRun)
			if err := manager.AddAuth(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "roles", "media", "activity-log", "search", "platform":
			manager := laravel.NewPlatformManager(dryRun)
			if err := manager.RunStep(target); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "docker", "security", "rate-limit", "health", "infra":
			manager := laravel.NewInfraManager(dryRun)
			if err := manager.RunStep(target); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "quality", "pro-arch", "docs-pro", "ci", "monitoring", "enterprise":
			manager := laravel.NewEnterpriseManager(dryRun)
			if err := manager.RunStep(target); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "all":
			manager := laravel.NewFullStackManager(dryRun)
			if err := manager.AddAll(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		default:
			printUsage()
			os.Exit(1)
		}
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  laravelboot new <project-name> [--dry-run]")
	fmt.Println("  laravelboot add auth [--dry-run]")
	fmt.Println("  laravelboot add roles [--dry-run]")
	fmt.Println("  laravelboot add media [--dry-run]")
	fmt.Println("  laravelboot add activity-log [--dry-run]")
	fmt.Println("  laravelboot add search [--dry-run]")
	fmt.Println("  laravelboot add platform [--dry-run]")
	fmt.Println("  laravelboot add docker [--dry-run]")
	fmt.Println("  laravelboot add security [--dry-run]")
	fmt.Println("  laravelboot add rate-limit [--dry-run]")
	fmt.Println("  laravelboot add health [--dry-run]")
	fmt.Println("  laravelboot add infra [--dry-run]")
	fmt.Println("  laravelboot add quality [--dry-run]")
	fmt.Println("  laravelboot add pro-arch [--dry-run]")
	fmt.Println("  laravelboot add docs-pro [--dry-run]")
	fmt.Println("  laravelboot add ci [--dry-run]")
	fmt.Println("  laravelboot add monitoring [--dry-run]")
	fmt.Println("  laravelboot add enterprise [--dry-run]")
	fmt.Println("  laravelboot add all [--dry-run]")
}
