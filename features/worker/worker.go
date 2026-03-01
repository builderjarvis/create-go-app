// Package worker provides a bounded-concurrency worker pool feature.
package worker

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&Worker{})
}

// Worker adds a generic bounded-concurrency worker pool package.
type Worker struct{}

func (w *Worker) Name() string        { return "worker" }
func (w *Worker) Description() string { return "Bounded-concurrency worker pool" }
func (w *Worker) Dependencies() []string { return nil }
func (w *Worker) Conflicts() []string    { return nil }

func (w *Worker) Install(ctx *scaffold.Context) error {
	ctx.AddPackage("golang.org/x/sync")
	return ctx.WriteTemplateDir(templates, "templates", "")
}
