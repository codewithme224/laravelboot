package presets

import "laravelboot/internal/config"

func GetPreset(name string) *config.Config {
	switch name {
	case "saas":
		return &config.Config{
			Database:     "postgres",
			Auth:         "sanctum",
			Features:     []string{"roles", "media", "search", "activity-log"},
			Infra:        []string{"docker", "security", "rate-limit", "health"},
			Architecture: "domain-based",
		}
	case "fintech":
		return &config.Config{
			Database:     "postgres",
			Auth:         "passport",
			Features:     []string{"roles", "activity-log"},
			Infra:        []string{"docker", "security", "rate-limit", "health"},
			Architecture: "domain-based",
		}
	case "enterprise":
		return &config.Config{
			Database: "postgres",
			Auth:     "sanctum",
			Features: []string{
				"roles", "media", "search", "activity-log",
				"reporting", "traits", "middleware", "exports", "jobs",
				"rules", "responses", "notifications", "scheduler", "cache",
				"versioning", "softdeletes", "storage", "events", "logging",
			},
			Infra:        []string{"docker", "security", "rate-limit", "health"},
			Enterprise:   []string{"quality", "pro-arch", "docs-pro", "ci", "monitoring", "helpers"},
			Architecture: "domain-based",
		}
	case "all":
		return &config.Config{
			Database: "postgres",
			Auth:     "passport",
			Features: []string{
				"roles", "media", "search", "activity-log",
				"reporting", "traits", "middleware", "exports", "jobs",
				"rules", "responses", "notifications", "scheduler", "cache",
				"versioning", "softdeletes", "storage", "events", "logging",
			},
			Infra:        []string{"docker", "security", "rate-limit", "health"},
			Enterprise:   []string{"quality", "pro-arch", "docs-pro", "ci", "monitoring", "tenancy", "helpers"},
			Architecture: "domain-based",
		}
	default:
		return config.DefaultConfig()
	}
}
