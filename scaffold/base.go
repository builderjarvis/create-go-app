package scaffold

import (
	"embed"
)

//go:embed all:base
var baseTemplates embed.FS

// installBase registers the base config fields and packages that are always
// included. It does NOT write any template files yet — that happens in
// writeBaseFiles after all features have had a chance to inject their own
// config fields and packages.
func installBase(ctx *Context) {
	// Base logging packages.
	ctx.AddPackage("github.com/lmittmann/tint")
	ctx.AddPackage("github.com/mattn/go-colorable")
	ctx.AddPackage("github.com/mattn/go-isatty")

	// Worker pool needs x/sync.
	ctx.AddPackage("golang.org/x/sync")

	// State file locking on Windows needs x/sys.
	ctx.AddPackage("golang.org/x/sys")

	// Register the base config field.
	ctx.AddConfig(ConfigField{
		Name:  "LogLevel",
		Key:   "LOG_LEVEL",
		Value: "info",
		Type:  "string",
	})
}

// writeBaseFiles renders all base templates into the project directory.
// Called after all features have been installed so that ConfigFields and
// Packages are fully populated before the templates are executed.
func writeBaseFiles(ctx *Context) error {
	return ctx.WriteTemplateDir(baseTemplates, "base", "")
}
