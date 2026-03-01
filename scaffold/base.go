package scaffold

import (
	"embed"
)

//go:embed all:base
var baseTemplates embed.FS

// installBase writes the base project files that are always included.
// This includes all utility packages (env, log, retry, worker, state, cycle, ptr).
func installBase(ctx *Context) error {
	// Base logging packages.
	ctx.AddPackage("github.com/lmittmann/tint")
	ctx.AddPackage("github.com/mattn/go-colorable")
	ctx.AddPackage("github.com/mattn/go-isatty")

	// Worker pool needs x/sync.
	ctx.AddPackage("golang.org/x/sync")

	// State file locking on Windows needs x/sys.
	ctx.AddPackage("golang.org/x/sys")

	return ctx.WriteTemplateDir(baseTemplates, "base", "")
}
