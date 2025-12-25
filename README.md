# 🚀 LaravelBoot

**LaravelBoot** is a professional, opinionated Go-based CLI tool designed to scaffold production-grade, API-first Laravel applications. It transforms the standard Laravel skeleton into a hardened, enterprise-ready API platform in seconds.

## 🌟 Key Features

- **API-First Defaults**: Pure JSON responses, standard pagination, and built-in response macros.
- **Domain-Driven Architecture**: Clean separation of concerns with a domain-based folder structure.
- **Enterprise Stack**: Integration with Spatie packages, Scramble (OpenAPI), Larastan, and Pest.
- **Hardened Infrastructure**: Production-ready Dockerfiles, Rate Limiting, Health checks, and Security middleware.
- **Config-Driven**: Reproducible setups via `.laravelboot.yaml`.
- **Interactive Mode**: Quick setup with an intuitive greeting and question flow.

---

## 🛠 Installation

### Quick Install (macOS/Linux)

```bash
curl -sL https://raw.githubusercontent.com/codewithme224/laravelboot/main/install.sh | bash
```

### Manual Install

1. Download the latest binary for your OS from the [Releases Page](https://github.com/codewithme224/laravelboot/releases).
2. Extract the archive and move the `laravelboot` binary to your path:
   ```bash
   chmod +x laravelboot
   sudo mv laravelboot /usr/local/bin/
   ```

### From Source

1. Clone the repository:
   ```bash
   git clone https://github.com/codewithme224/laravelboot.git
   cd laravelboot
   ```
2. Build the binary:
   ```bash
   go build -o laravelboot ./cmd/laravelboot/main.go
   ```
3. (Optional) Move to your bin:
   ```bash
   sudo mv laravelboot /usr/local/bin/
   ```

---

## � Updating

To keep your CLI tool up to date with the latest features and security fixes:

### Modern Versions (v1.0.0+)

Simply run:

```bash
laravelboot update
```

### Older Versions (< v1.0.0)

If you are on an older version that doesn't support the `update` command, run the install script again to perform a manual upgrade:

```bash
curl -sL https://raw.githubusercontent.com/codewithme224/laravelboot/main/install.sh | bash
```

---

## �🚀 Usage

### 1. Interactive Initialize

Generate your project configuration interactively:

```bash
laravelboot init
```

### 2. Create a New Project

Create a complete API project based on your config or defaults:

```bash
laravelboot new my-api-project
```

### 3. Utility Commands

```bash
laravelboot version         # Show current version
laravelboot update          # Self-update to the latest version
```

### 4. Add Features Incrementally

You can add specific stacks to an existing project:

#### Authentication & DB

```bash
laravelboot add auth        # Sanctum + Base Auth Controller
```

#### Platform Features

```bash
laravelboot add roles         # Spatie Permissions
laravelboot add media         # Spatie MediaLibrary + SpatieMediaService
laravelboot add search        # Laravel Scout + Typesense
laravelboot add activity      # Spatie ActivityLog
laravelboot add reporting     # Excel (Maatwebsite) + PDF (dompdf)
laravelboot add traits        # Common API traits (Api, HandlesPagination, Auditable)
laravelboot add middleware    # DBTransaction + ForceJson middleware
laravelboot add exports       # Base Export/Import classes for Excel
laravelboot add jobs          # Base Job class with queue support
laravelboot add rules         # Custom validation rules (Base64Image, PhoneNumber, etc.)
laravelboot add responses     # API response helpers + Exception handler
laravelboot add notifications # Notifications system with services
laravelboot add scheduler     # Console commands + scheduling
laravelboot add cache         # Caching layer with Redis + Cacheable trait
laravelboot add versioning    # API versioning (v1, v2 structure)
laravelboot add softdeletes   # Soft deletes + Trash management
laravelboot add storage       # File storage service + controller
laravelboot add events        # Events & Listeners scaffolding
laravelboot add logging       # Request logging + Slack notifications
laravelboot add platform      # All of the above
```

#### Infrastructure & Security

```bash
laravelboot add docker      # Dev & Prod Dockerfiles + Compose
laravelboot add security    # Force JSON middleware + Env validation
laravelboot add health      # Health & Readiness endpoints
laravelboot add rate-limit  # API Throttling
laravelboot add infra       # All of the above
```

#### Enterprise & Quality

```bash
laravelboot add quality     # Pint + PHPStan + Pest
laravelboot add docs-pro    # Automated Swagger (Scramble)
laravelboot add monitoring  # Laravel Pulse for monitoring
laravelboot add ci          # GitHub Actions + GitLab CI Workflows
laravelboot add tenancy     # Multi-tenancy (stancl/tenancy)
laravelboot add helpers     # Global helpers.php with auto-registration
laravelboot add enterprise  # All of the above
```

#### The "Giga" Stack

```bash
laravelboot add all         # INSTALL EVERY SINGLE FEATURE (Phase 2-5)
```

---

## ⚙️ Configuration (`.laravelboot.yaml`)

Teams can standardize their backend architecture using the configuration file:

```yaml
project_name: myapp
database: postgres # mysql, postgres, sqlite
auth: sanctum # sanctum, passport
features:
  - roles
  - search
  - activity-log
  - reporting
  - traits
  - middleware
  - exports
  - jobs
infra:
  - docker
  - security
  - health
enterprise:
  - quality
  - ci
  - docs-pro
  - tenancy
  - helpers
architecture: domain-based
```

---

## 🧪 Testing with Dry Run

Simulate any command without touching the filesystem:

```bash
laravelboot new myapp --dry-run
laravelboot add enterprise --dry-run
```

## 🏗 Architecture Philosophy

- **Controllers**: Thin and focused on request/response.
- **Logic**: Encapsulated in **Actions** and **Services**.
- **Data**: Strictly typed using **DTOs** (Spatie Laravel Data).
- **Queries**: Enforced through **Query Objects** (Spatie Query Builder).
- **Strictly JSON**: No sessions, no blades, no noise.

---

Built with ❤️ for Laravel Platform Engineers.
