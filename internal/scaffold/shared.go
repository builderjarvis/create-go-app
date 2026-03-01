package scaffold

import (
	"embed"
)

//go:embed all:shared
var sharedTemplates embed.FS

// writeSharedFiles renders templates that depend on injections from all features.
// These are written last, after all features have injected their content.
func writeSharedFiles(ctx *Context) error {
	return ctx.WriteTemplateDir(sharedTemplates, "shared", "")
}
