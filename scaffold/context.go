package scaffold

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// ConfigField describes a single environment variable that gets rendered into
// the generated Config struct and .env.example file.
type ConfigField struct {
	Name     string // Go field name, e.g. "DatabaseURL"
	Key      string // Env var name, e.g. "DATABASE_URL"
	Value    string // Default value for .env.example, e.g. "postgres://localhost:5432/app"
	Type     string // Go type, e.g. "string"
	Required bool
}

// Context is the shared state passed to every feature's Install method.
type Context struct {
	// ProjectName is the directory/app name (e.g., "my-app").
	ProjectName string

	// ModulePath is the Go module path (e.g., "github.com/user/my-app").
	ModulePath string

	// GoVersion is the Go version for go.mod (e.g., "1.25.0").
	GoVersion string

	// ProjectDir is the absolute path to the generated project.
	ProjectDir string

	// Active holds all resolved features by name.
	Active map[string]Feature

	// Packages collects package paths for `go get`.
	Packages []string

	// ConfigFields holds the ordered list of env config fields contributed by
	// base and each feature. Templates range over this to render Config struct
	// and .env.example.
	ConfigFields []ConfigField
}

// Has returns true if a feature is active.
func (c *Context) Has(name string) bool {
	_, ok := c.Active[name]
	return ok
}

// AddPackage adds a Go module dependency.
func (c *Context) AddPackage(pkg string) {
	c.Packages = append(c.Packages, pkg)
}

// AddConfig appends a ConfigField to the context.
func (c *Context) AddConfig(field ConfigField) {
	c.ConfigFields = append(c.ConfigFields, field)
}

// WriteFile writes content to a file relative to ProjectDir, creating directories as needed.
func (c *Context) WriteFile(relPath string, content []byte) error {
	abs := filepath.Join(c.ProjectDir, relPath)
	if err := os.MkdirAll(filepath.Dir(abs), 0755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(abs), err)
	}
	return os.WriteFile(abs, content, 0644)
}

// WriteTemplate renders a single template file from an embed.FS and writes it to relPath.
func (c *Context) WriteTemplate(fsys fs.FS, tmplPath string, relPath string) error {
	data, err := fs.ReadFile(fsys, tmplPath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", tmplPath, err)
	}

	t, err := template.New(filepath.Base(tmplPath)).
		Funcs(c.templateFuncs()).
		Parse(string(data))
	if err != nil {
		return fmt.Errorf("parse template %s: %w", tmplPath, err)
	}

	var buf strings.Builder
	if err := t.Execute(&buf, c.templateData()); err != nil {
		return fmt.Errorf("execute template %s: %w", tmplPath, err)
	}

	content := strings.TrimSpace(buf.String())
	if content == "" {
		return nil
	}

	return c.WriteFile(relPath, []byte(buf.String()))
}

// WriteTemplateDir walks an embed.FS directory and writes all files, rendering .tmpl files
// as templates and stripping the .tmpl extension.
func (c *Context) WriteTemplateDir(fsys fs.FS, dir string, outPrefix string) error {
	return fs.WalkDir(fsys, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		rel := strings.TrimPrefix(path, dir+"/")
		outPath := filepath.Join(outPrefix, rel)

		if strings.HasSuffix(path, ".tmpl") {
			outPath = strings.TrimSuffix(outPath, ".tmpl")
			return c.WriteTemplate(fsys, path, outPath)
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}
		return c.WriteFile(outPath, data)
	})
}

// templateData returns the data map passed to all templates.
func (c *Context) templateData() map[string]any {
	return map[string]any{
		"ProjectName":  c.ProjectName,
		"ModulePath":   c.ModulePath,
		"GoVersion":    c.GoVersion,
		"Has":          c.Has,
		"Packages":     c.Packages,
		"ConfigFields": c.ConfigFields,
	}
}

// templateFuncs returns custom template functions.
func (c *Context) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"has":  c.Has,
		"join": strings.Join,
	}
}
