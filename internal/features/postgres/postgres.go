// Package postgres provides PostgreSQL database support with pgx, goose migrations, and sqlc.
package postgres

import (
	"embed"

	"github.com/builderjarvis/create-go-app/internal/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Postgres{})
}

// Postgres adds PostgreSQL database support with pgx pool, goose migrations,
// and sqlc code generation.
type Postgres struct{}

func (p *Postgres) Name() string        { return "postgres" }
func (p *Postgres) Description() string { return "PostgreSQL with pgx, goose migrations, and sqlc" }
func (p *Postgres) Dependencies() []string { return nil }
func (p *Postgres) Conflicts() []string    { return nil }

func (p *Postgres) Install(ctx *scaffold.Context) error {
	ctx.AddPackage("github.com/jackc/pgx/v5")

	ctx.Inject("config_fields", `DatabaseURL string `+"`"+`env:"DATABASE_URL,required"`+"`")
	ctx.Inject("env_vars", "\n# PostgreSQL Database URL\nDATABASE_URL=postgres://postgres:postgres@localhost:5432/app?sslmode=disable")
	ctx.Inject("main_imports", `"fmt"`)
	ctx.Inject("main_imports", `"`+ctx.ModulePath+`/pkg/db"`)
	ctx.Inject("main_init", `database, err := db.New(ctx, config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("create database: %w", err)
	}
	defer database.Close()

	_ = database`)

	ctx.Inject("compose_services", `  postgres:
    profiles:
      - postgres
    image: postgres:17-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: app
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres -d app"]
      interval: 5s
      timeout: 5s
      retries: 5`)

	return ctx.WriteTemplateDir(templates, "templates", "")
}
