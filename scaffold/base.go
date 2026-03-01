package scaffold

import (
	"embed"
)

//go:embed all:base
var baseTemplates embed.FS

// installBase writes the base project files that are always included.
func installBase(ctx *Context) error {
	// Always include logging packages.
	ctx.AddPackage("github.com/lmittmann/tint")
	ctx.AddPackage("github.com/mattn/go-colorable")
	ctx.AddPackage("github.com/mattn/go-isatty")

	return ctx.WriteTemplateDir(baseTemplates, "base", "")
}
