// Package docker provides Docker Compose support.
package docker

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Docker{})
}

// Docker adds Dockerfile and compose.yaml to the project.
type Docker struct{}

func (d *Docker) Name() string           { return "docker" }
func (d *Docker) Description() string    { return "Dockerfile + compose with watch" }
func (d *Docker) Dependencies() []string { return nil }
func (d *Docker) Conflicts() []string    { return nil }

func (d *Docker) Install(ctx *scaffold.Context) error {
	return ctx.WriteTemplateDir(templates, "templates", "")
}
