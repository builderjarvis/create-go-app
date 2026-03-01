# create-go-app

Scaffold production-ready Go projects with an interactive CLI and composable optional features.

## Quick Start

```bash
go run github.com/builderjarvis/create-go-app@latest
```

Or pass the project name directly:

```bash
go run github.com/builderjarvis/create-go-app@latest my-app
```

The CLI prompts for your project name, module path, and feature selection — then generates a fully wired project with a `go.mod`, initial git commit, and working `make dev`.

## Base (Always Included)

Every project gets these packages out of the box:

| Package | Description |
|---------|-------------|
| `pkg/env` | Typed config from environment variables (`.env` support) |
| `pkg/log` | Structured logging via `log/slog` + `tint` |
| `pkg/retry` | Retry with exponential backoff and jitter |
| `pkg/worker` | Bounded-concurrency worker pool (`x/sync`) |
| `pkg/state` | File-backed JSON state with cross-platform file locking |
| `pkg/cycle` | Lifecycle manager for graceful startup/shutdown |
| `pkg/ptr` | Generic pointer helpers |

## Optional Features

Select any combination during project creation:

| Feature | What it adds |
|---------|-------------|
| **postgres** | `pkg/db` — pgx/v5 connection pool, goose migrations, sqlc config |
| **docker** | `Dockerfile` + `compose.yaml` with hot-reload via `docker compose watch` |
| **httpclient** | `pkg/client` — HTTP client with TLS fingerprinting (fhttp + mimic + uTLS) |
| **ci** | `.github/workflows/ci.yaml` — GitHub Actions CI pipeline |

## Feature Composition

Features are aware of each other. Selecting multiple features produces a correctly integrated result:

- **postgres + docker** → `compose.yaml` includes a `postgres` service (image: `postgres:17-alpine`) with health checks
- **postgres** alone → `pkg/db` wired into `cmd/app/main.go` with connection setup and teardown
- **httpclient** alone → `pkg/client` initialized in `main.go` with a ready-to-use client

Dependencies are declared per-feature and resolved automatically. No manual wiring needed.

## Generated Structure

```
my-app/
├── cmd/app/main.go       # entry point, wired to selected features
├── pkg/
│   ├── env/              # config (always)
│   ├── log/              # logging (always)
│   ├── retry/            # retry (always)
│   ├── worker/           # worker pool (always)
│   ├── state/            # state persistence (always)
│   ├── cycle/            # lifecycle (always)
│   ├── ptr/              # pointer helpers (always)
│   ├── db/               # postgres (if selected)
│   └── client/           # httpclient (if selected)
├── Makefile
├── .env.example
├── .gitignore
├── Dockerfile            # docker (if selected)
├── compose.yaml          # docker (if selected)
└── .github/workflows/
    └── ci.yaml           # ci (if selected)
```

## License

MIT
