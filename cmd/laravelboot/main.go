package main

import (
	"fmt"
	"laravelboot/internal/interactive"
	"laravelboot/internal/laravel"
	"laravelboot/internal/utils"
	"os"
)

const VERSION = "v1.0.0"

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

	if len(args) < 1 {
		printUsage()
		os.Exit(1)
	}

	command := args[0]
	target := ""
	if len(args) > 1 {
		target = args[1]
	}

	switch command {
	case "version":
		fmt.Printf("LaravelBoot %s\n", VERSION)

	case "update":
		if err := utils.SelfUpdate(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

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
		if !dryRun {
			go utils.CheckForUpdate(VERSION)
		}
		if target == "" {
			printUsage()
			os.Exit(1)
		}
		appName := target
		preset := ""

		// Check for preset flags or positional
		for _, arg := range os.Args[1:] {
			if arg == "--all" {
				preset = "all"
			} else if arg == "--enterprise" {
				preset = "enterprise"
			}
		}

		// Fallback to third arg if defined (e.g. laravelboot new myapp saas)
		if len(args) > 2 && preset == "" {
			preset = args[2]
		}

		creator := laravel.NewCreator(appName, preset, dryRun)
		if err := creator.Create(); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
			os.Exit(1)
		}

	case "add":
		if !dryRun {
			go utils.CheckForUpdate(VERSION)
		}
		if target == "" {
			printUsage()
			os.Exit(1)
		}
		cwd, _ := os.Getwd()
		switch target {
		case "auth":
			manager := laravel.NewAuthManager(cwd, dryRun)
			if err := manager.AddAuth(); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "roles", "media", "activity-log", "search", "platform":
			manager := laravel.NewPlatformManager(cwd, dryRun)
			if err := manager.RunStep(target); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "docker", "security", "rate-limit", "health", "infra":
			manager := laravel.NewInfraManager(cwd, dryRun)
			if err := manager.RunStep(target); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "quality", "pro-arch", "docs-pro", "ci", "monitoring", "enterprise":
			manager := laravel.NewEnterpriseManager(cwd, dryRun)
			if err := manager.RunStep(target); err != nil {
				fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
				os.Exit(1)
			}
		case "all":
			manager := laravel.NewFullStackManager(cwd, dryRun)
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
	fmt.Printf("LaravelBoot %s\n\n", VERSION)
	fmt.Println("Usage:")
	fmt.Println("  laravelboot init                    Initialize configuration")
	fmt.Println("  laravelboot new <project-name>      Create new project")
	fmt.Println("  laravelboot update                  Update CLI tool")
	fmt.Println("  laravelboot version                 Show version")
	fmt.Println("\nAdd Stacks:")
	fmt.Println("  laravelboot add auth                Passport/Sanctum Auth")
	fmt.Println("  laravelboot add platform            Roles, Media, Search, Activity")
	fmt.Println("  laravelboot add infra               Docker, Health, Security")
	fmt.Println("  laravelboot add enterprise          Pest, Scramle, CI, Monitoring")
	fmt.Println("  laravelboot add all                 Install EVERYTHING")
}
