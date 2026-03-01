// Package docker provides Docker Compose support.
package docker

import (
	"github.com/builderjarvis/create-go-app/scaffold"
)

func init() {
	scaffold.Register(&Docker{})
}

// Docker adds Dockerfile and compose.yaml to the project.
// The actual files are in shared templates since they depend on feature composition.
type Docker struct{}

func (d *Docker) Name() string           { return "docker" }
func (d *Docker) Description() string    { return "Dockerfile + compose with watch" }
func (d *Docker) Dependencies() []string { return nil }
func (d *Docker) Conflicts() []string    { return nil }

func (d *Docker) Install(ctx *scaffold.Context) error {
	// Docker files are rendered as shared templates (they depend on which other features are active).
	return nil
}
