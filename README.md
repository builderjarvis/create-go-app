# create-go-app

A scaffolder for Go projects with composable optional features. Inspired by [create-t3-app](https://github.com/t3-oss/create-t3-app).

## Usage

```bash
go run github.com/builderjarvis/create-go-app@latest my-app
```

Or interactively:

```bash
go run github.com/builderjarvis/create-go-app@latest
```

## Features

Each feature is self-contained. Adding a new feature = adding a new package, zero changes to existing code.

| Feature | Description |
|---------|-------------|
| **postgres** | PostgreSQL with pgx, goose migrations, and sqlc |
| **docker** | Dockerfile and Docker Compose |
| **httpclient** | HTTP client with TLS fingerprinting (fhttp + mimic) |
| **worker** | Bounded-concurrency worker pool |
| **state** | File-backed JSON state persistence with file locking |
| **retry** | Retry with exponential backoff and jitter |

Features declare dependencies and conflicts. The scaffolder automatically resolves the dependency graph and installs features in the correct order.

## Generated Project Structure

```
my-app/
├── cmd/app/main.go          # Entry point
├── pkg/
│   ├── env/                  # Environment config (.env loading + typed Config)
│   ├── log/                  # Structured logging (slog + tint)
│   ├── db/                   # PostgreSQL (if selected)
│   ├── client/               # HTTP client (if selected)
│   ├── worker/               # Worker pool (if selected)
│   ├── state/                # State persistence (if selected)
│   └── retry/                # Retry with backoff (if selected)
├── Makefile
├── .env.example
├── Dockerfile                # If docker selected
├── compose.yaml              # If docker selected
└── .gitignore
```

## Architecture

- **Feature interface**: Each feature implements `scaffold.Feature` with `Name()`, `Description()`, `Dependencies()`, `Conflicts()`, and `Install()`
- **Injection points**: Features inject content into shared templates (imports, config fields, env vars, compose services, etc.)
- **Two-pass rendering**: Base templates are written first, then features install their packages, then shared templates render with all injections
- **Dependency resolution**: Topological sort with conflict detection
- **`embed.FS` per feature**: Templates are co-located with feature code

### Adding a New Feature

1. Create `internal/features/myfeature/myfeature.go` implementing `scaffold.Feature`
2. Create `internal/features/myfeature/templates/` with embedded template files
3. Add `_ "github.com/builderjarvis/create-go-app/internal/features/myfeature"` to `cmd/create-go-app/main.go`
4. Done. The registry picks it up automatically.
