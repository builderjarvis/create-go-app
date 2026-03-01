// Package httpclient provides an HTTP client with TLS fingerprinting via fhttp + mimic.
package httpclient

import (
	"embed"

	"github.com/builderjarvis/create-go-app/scaffold"
)

//go:embed all:templates
var templates embed.FS

func init() {
	scaffold.Register(&HTTPClient{})
}

// HTTPClient adds an HTTP client package with TLS fingerprinting support.
type HTTPClient struct{}

func (h *HTTPClient) Name() string        { return "httpclient" }
func (h *HTTPClient) Description() string { return "HTTP client with TLS fingerprinting (fhttp + mimic)" }
func (h *HTTPClient) Dependencies() []string { return nil }
func (h *HTTPClient) Conflicts() []string    { return nil }

func (h *HTTPClient) Install(ctx *scaffold.Context) error {
	ctx.AddPackage("github.com/saucesteals/fhttp")
	ctx.AddPackage("github.com/saucesteals/mimic")
	ctx.AddPackage("github.com/refraction-networking/utls")
	ctx.AddPackage("github.com/andybalholm/brotli")
	ctx.AddPackage("github.com/klauspost/compress")

	return ctx.WriteTemplateDir(templates, "templates", "")
}
