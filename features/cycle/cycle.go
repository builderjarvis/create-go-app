// Package cycle provides a thread-safe round-robin file rotator feature.
package cycle

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Cycle{})
}

// Cycle adds a thread-safe round-robin file rotator package.
type Cycle struct{}

func (c *Cycle) Name() string        { return "cycle" }
func (c *Cycle) Description() string { return "Thread-safe round-robin file rotator" }
func (c *Cycle) Dependencies() []string { return nil }
func (c *Cycle) Conflicts() []string    { return nil }

func (c *Cycle) Install(ctx *scaffold.Context) error {
	return ctx.WriteTemplateDir(templates, "templates", "")
}
