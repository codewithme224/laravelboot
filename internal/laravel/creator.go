package laravel

import (
	"fmt"
	"laravelboot/internal/config"
	"laravelboot/internal/docs"
	"laravelboot/internal/plugins"
	"laravelboot/internal/presets"
	"os"
)

type Creator struct {
	Name   string
	DryRun bool
	Config *config.Config
}

func NewCreator(name string, preset string, dryRun bool) *Creator {
	// Try to load config if it exists
	conf, _ := config.LoadConfig(".laravelboot.yaml")
	if conf == nil {
		if preset != "" {
			conf = presets.GetPreset(preset)
		} else {
			conf = config.DefaultConfig()
		}
	}
	conf.ProjectName = name

	return &Creator{
		Name:   name,
		DryRun: dryRun,
		Config: conf,
	}
}

func (c *Creator) Create() error {
	fmt.Printf("🚀 Creating new Laravel API project: %s (Config-Driven)\n", c.Name)

	installer := NewInstaller(c.DryRun)
	if err := installer.CheckDependencies(); err != nil {
		return err
	}

	if err := installer.CreateProject(c.Name); err != nil {
		return err
	}

	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}
	projectPath = fmt.Sprintf("%s/%s", projectPath, c.Name)

	// Base Architecture
	arch := NewArchitecture(projectPath, c.DryRun)
	if err := arch.SetupFolders(); err != nil {
		return err
	}

	// API Setup
	api := NewApiSetup(projectPath, c.DryRun)
	if err := api.Configure(); err != nil {
		return err
	}

	// Spatie Query Builder (Core in Phase 1)
	spatie := NewSpatieQueryBuilder(projectPath, c.DryRun)
	if err := spatie.Install(); err != nil {
		return err
	}
	if err := spatie.CreateExample(); err != nil {
		return err
	}

	// Apply Auth from config
	if c.Config.Auth != "" {
		authMgr := NewAuthManager(projectPath, c.DryRun)
		if err := authMgr.AddAuth(); err != nil {
			fmt.Printf("⚠️ Warning: auth setup failed: %v\n", err)
		}
	}

	// Apply features from config
	platform := NewPlatformManager(projectPath, c.DryRun)
	for _, feature := range c.Config.Features {
		fmt.Printf("📦 Adding feature: %s\n", feature)
		if err := platform.RunStep(feature); err != nil {
			fmt.Printf("⚠️ Warning: feature %s failed: %v\n", feature, err)
		}
	}

	// Apply infra from config
	infra := NewInfraManager(projectPath, c.DryRun)
	for _, step := range c.Config.Infra {
		fmt.Printf("🛡️ Adding infra: %s\n", step)
		if err := infra.RunStep(step); err != nil {
			fmt.Printf("⚠️ Warning: infra %s failed: %v\n", step, err)
		}
	}

	// Apply enterprise from config
	enterprise := NewEnterpriseManager(projectPath, c.DryRun)
	for _, step := range c.Config.Enterprise {
		fmt.Printf("👑 Adding enterprise feature: %s\n", step)
		if err := enterprise.RunStep(step); err != nil {
			fmt.Printf("⚠️ Warning: enterprise feature %s failed: %v\n", step, err)
		}
	}

	// Run Plugins
	pluginMgr := plugins.NewPluginManager()
	pluginMgr.Register(&plugins.LogPlugin{})
	if err := pluginMgr.RunAll(c.Config, projectPath); err != nil {
		return err
	}

	// Generate Docs
	if err := docs.Generate(projectPath, c.Config, c.DryRun); err != nil {
		return err
	}

	fmt.Printf("\n✨ Project '%s' created successfully!\n", c.Name)
	return nil
}
