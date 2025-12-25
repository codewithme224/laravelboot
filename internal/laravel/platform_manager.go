package laravel

import (
	"fmt"
)

type PlatformManager struct {
	ProjectPath string
	DryRun      bool
}

func NewPlatformManager(projectPath string, dryRun bool) *PlatformManager {
	return &PlatformManager{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *PlatformManager) RunStep(name string) error {
	switch name {
	case "roles":
		return NewRolesSetup(m.ProjectPath, m.DryRun).Setup()
	case "media":
		return NewMediaSetup(m.ProjectPath, m.DryRun).Setup()
	case "activity-log":
		return NewActivityLogSetup(m.ProjectPath, m.DryRun).Setup()
	case "search":
		return NewSearchSetup(m.ProjectPath, m.DryRun).Setup()
	case "reporting":
		return NewReportingSetup(m.ProjectPath, m.DryRun).Setup()
	case "traits":
		return NewTraitsSetup(m.ProjectPath, m.DryRun).Setup()
	case "middleware":
		return NewMiddlewareSetup(m.ProjectPath, m.DryRun).Setup()
	case "exports":
		return NewExportsSetup(m.ProjectPath, m.DryRun).Setup()
	case "jobs":
		return NewJobsSetup(m.ProjectPath, m.DryRun).Setup()
	case "rules":
		return NewRulesSetup(m.ProjectPath, m.DryRun).Setup()
	case "responses":
		return NewResponsesSetup(m.ProjectPath, m.DryRun).Setup()
	case "notifications":
		return NewNotificationsSetup(m.ProjectPath, m.DryRun).Setup()
	case "scheduler":
		return NewSchedulerSetup(m.ProjectPath, m.DryRun).Setup()
	case "cache":
		return NewCacheSetup(m.ProjectPath, m.DryRun).Setup()
	case "versioning":
		return NewVersioningSetup(m.ProjectPath, m.DryRun).Setup()
	case "softdeletes":
		return NewSoftDeletesSetup(m.ProjectPath, m.DryRun).Setup()
	case "storage":
		return NewStorageSetup(m.ProjectPath, m.DryRun).Setup()
	case "events":
		return NewEventsSetup(m.ProjectPath, m.DryRun).Setup()
	case "logging":
		return NewLoggingSetup(m.ProjectPath, m.DryRun).Setup()
	case "platform":
		fmt.Println("🚀 Installing complete platform stack...")
		if err := NewRolesSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewMediaSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewActivityLogSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewSearchSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewReportingSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewTraitsSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewMiddlewareSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewExportsSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewJobsSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewRulesSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewResponsesSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewNotificationsSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewSchedulerSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewCacheSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewVersioningSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewSoftDeletesSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewStorageSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewEventsSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		if err := NewLoggingSetup(m.ProjectPath, m.DryRun).Setup(); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown platform feature: %s", name)
	}
}
