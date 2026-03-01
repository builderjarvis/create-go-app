package scaffold

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Config holds the user's selections for project generation.
type Config struct {
	ProjectName string
	ModulePath  string
	GoVersion   string
	Features    []string // selected feature names
	OutputDir   string   // parent directory (project created as subdir)
}

// Generate creates a new project from the given configuration.
func Generate(cfg Config) error {
	// Resolve features (expand deps, check conflicts, topo sort).
	features, err := Resolve(cfg.Features)
	if err != nil {
		return fmt.Errorf("resolving features: %w", err)
	}

	projectDir := filepath.Join(cfg.OutputDir, cfg.ProjectName)

	// Ensure project directory doesn't already exist.
	if _, err := os.Stat(projectDir); err == nil {
		return fmt.Errorf("directory %s already exists", projectDir)
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("creating project directory: %w", err)
	}

	// Build context.
	ctx := &Context{
		ProjectName: cfg.ProjectName,
		ModulePath:  cfg.ModulePath,
		GoVersion:   cfg.GoVersion,
		ProjectDir:  projectDir,
		Active:      make(map[string]Feature),
	}

	for _, f := range features {
		ctx.Active[f.Name()] = f
	}

	// Install base project first (all utility packages).
	if err := installBase(ctx); err != nil {
		return fmt.Errorf("installing base: %w", err)
	}

	// Install each feature in dependency order.
	for _, f := range features {
		if err := f.Install(ctx); err != nil {
			return fmt.Errorf("installing feature %s: %w", f.Name(), err)
		}
	}

	// Write shared files that depend on injections from all features.
	if err := writeSharedFiles(ctx); err != nil {
		return fmt.Errorf("writing shared files: %w", err)
	}

	// Initialize go module and tidy.
	if err := initGoModule(ctx); err != nil {
		return fmt.Errorf("initializing go module: %w", err)
	}

	// Initialize git repo.
	if err := initGit(projectDir); err != nil {
		return fmt.Errorf("initializing git: %w", err)
	}

	return nil
}

// initGoModule runs go mod init and go mod tidy.
func initGoModule(ctx *Context) error {
	cmds := [][]string{
		{"go", "mod", "init", ctx.ModulePath},
		{"go", "mod", "edit", "-go=" + ctx.GoVersion},
	}

	// Add all collected packages.
	for _, pkg := range ctx.Packages {
		cmds = append(cmds, []string{"go", "get", pkg})
	}

	cmds = append(cmds, []string{"go", "mod", "tidy"})

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = ctx.ProjectDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("running %s: %w", strings.Join(args, " "), err)
		}
	}

	return nil
}

// initGit initializes a git repo and makes an initial commit.
func initGit(dir string) error {
	cmds := [][]string{
		{"git", "init"},
		{"git", "add", "-A"},
		{"git", "commit", "-m", "Initial commit from create-go-app"},
	}

	for _, args := range cmds {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("running %s: %w", strings.Join(args, " "), err)
		}
	}

	return nil
}
