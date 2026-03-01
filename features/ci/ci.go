// Package ci provides GitHub Actions CI workflow support.
package ci

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&CI{})
}

// CI adds a GitHub Actions CI workflow.
type CI struct{}

func (c *CI) Name() string           { return "ci" }
func (c *CI) Description() string    { return "GitHub Actions CI" }
func (c *CI) Dependencies() []string { return nil }
func (c *CI) Conflicts() []string    { return nil }

func (c *CI) Install(ctx *scaffold.Context) error {
	return ctx.WriteTemplateDir(templates, "templates", "")
}
