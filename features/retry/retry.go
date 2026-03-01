// Package retry provides the retry with exponential backoff feature.
package retry

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Retry{})
}

// Retry adds a generic retry-with-backoff package.
type Retry struct{}

func (r *Retry) Name() string        { return "retry" }
func (r *Retry) Description() string { return "Retry with exponential backoff and jitter" }
func (r *Retry) Dependencies() []string { return nil }
func (r *Retry) Conflicts() []string    { return nil }

func (r *Retry) Install(ctx *scaffold.Context) error {
	return ctx.WriteTemplateDir(templates, "templates", "")
}
