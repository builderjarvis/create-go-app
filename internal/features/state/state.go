// Package state provides a file-backed JSON state persistence feature.
package state

import (
	"embed"

	"github.com/builderjarvis/create-go-app/internal/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&State{})
}

// State adds a generic file-backed JSON state package with file locking.
type State struct{}

func (s *State) Name() string        { return "state" }
func (s *State) Description() string { return "File-backed JSON state persistence with file locking" }
func (s *State) Dependencies() []string { return nil }
func (s *State) Conflicts() []string    { return nil }

func (s *State) Install(ctx *scaffold.Context) error {
	return ctx.WriteTemplateDir(templates, "templates", "")
}
