// Package ptr provides the pointer helpers feature.
package ptr

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Ptr{})
}

// Ptr adds generic pointer helper utilities.
type Ptr struct{}

func (p *Ptr) Name() string        { return "ptr" }
func (p *Ptr) Description() string { return "Generic pointer helpers (To, ValueOrEmpty)" }
func (p *Ptr) Dependencies() []string { return nil }
func (p *Ptr) Conflicts() []string    { return nil }

func (p *Ptr) Install(ctx *scaffold.Context) error {
	return ctx.WriteTemplateDir(templates, "templates", "")
}
