// Package example provides an example service feature that wires HTTP client and database together.
package example

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Example{})
}

// Example adds a skeleton service that demonstrates wiring an HTTP client and database together.
type Example struct{}

func (e *Example) Name() string        { return "example" }
func (e *Example) Description() string { return "Example service wiring HTTP client + database" }
func (e *Example) Dependencies() []string { return []string{"httpclient", "postgres"} }
func (e *Example) Conflicts() []string    { return nil }

func (e *Example) Install(ctx *scaffold.Context) error {
	ctx.Inject("main_imports", `"`+ctx.ModulePath+`/pkg/example"`)
	ctx.Inject("main_init", `svc, err := example.New(database, nil)
	if err != nil {
		return fmt.Errorf("create example: %w", err)
	}

	resp, err := svc.Example(ctx)
	if err != nil {
		return fmt.Errorf("example: %w", err)
	}

	slog.InfoContext(ctx, "example response", "response", resp)`)

	return ctx.WriteTemplateDir(templates, "templates", "")
}
