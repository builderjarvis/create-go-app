// Package cli provides the interactive command-line interface for create-go-app.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/builderjarvis/create-go-app/internal/scaffold"
)

// Run is the main entry point for the CLI.
func Run(args []string) error {
	reader := bufio.NewReader(os.Stdin)

	// Project name from args or prompt.
	var projectName string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		projectName = args[0]
	} else {
		projectName = prompt(reader, "Project name", "my-app")
	}

	// Module path.
	defaultModule := "github.com/user/" + projectName
	modulePath := prompt(reader, "Go module path", defaultModule)

	// Go version.
	goVersion := runtime.Version()
	goVersion = strings.TrimPrefix(goVersion, "go")
	// Trim patch for go.mod (e.g., "1.25.0" stays, "1.25rc1" stays).

	fmt.Println()
	fmt.Println("Optional features:")
	fmt.Println()

	// Show available features and let user select.
	allFeatures := scaffold.All()
	selected := make([]string, 0)

	for _, f := range allFeatures {
		answer := prompt(reader, fmt.Sprintf("  Include %s? (%s)", f.Name(), f.Description()), "y/N")
		if isYes(answer) {
			selected = append(selected, f.Name())
		}
	}

	fmt.Println()

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	cfg := scaffold.Config{
		ProjectName: projectName,
		ModulePath:  modulePath,
		GoVersion:   goVersion,
		Features:    selected,
		OutputDir:   cwd,
	}

	fmt.Printf("Creating %s in %s...\n", projectName, filepath.Join(cwd, projectName))
	fmt.Println()

	if err := scaffold.Generate(cfg); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println("✓ Project created successfully!")
	fmt.Println()
	fmt.Printf("  cd %s\n", projectName)
	fmt.Println("  cp .env.example .env")
	fmt.Println("  make dev")
	fmt.Println()

	return nil
}

// prompt shows a prompt and returns user input (or default if empty).
func prompt(reader *bufio.Reader, label string, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultVal
	}
	return input
}

func isYes(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == "y" || s == "yes"
}
