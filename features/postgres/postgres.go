// Package postgres provides PostgreSQL database support with pgx, goose migrations, and sqlc.
package postgres

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Postgres{})
}

// Postgres adds PostgreSQL database support with pgx pool, goose migrations,
// and sqlc code generation.
type Postgres struct{}

func (p *Postgres) Name() string           { return "postgres" }
func (p *Postgres) Description() string    { return "pgx + sqlc + goose migrations" }
func (p *Postgres) Dependencies() []string { return nil }
func (p *Postgres) Conflicts() []string    { return nil }

func (p *Postgres) Install(ctx *scaffold.Context) error {
	ctx.AddPackage("github.com/jackc/pgx/v5")

	return ctx.WriteTemplateDir(templates, "templates", "")
}
