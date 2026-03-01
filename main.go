// Package main is the entry point for the create-go-app scaffolder.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/builderjarvis/create-go-app/scaffold"

	// Register all features via blank imports.
	_ "github.com/builderjarvis/create-go-app/features/ci"
	_ "github.com/builderjarvis/create-go-app/features/cycle"
	_ "github.com/builderjarvis/create-go-app/features/docker"
	_ "github.com/builderjarvis/create-go-app/features/example"
	_ "github.com/builderjarvis/create-go-app/features/httpclient"
	_ "github.com/builderjarvis/create-go-app/features/postgres"
	_ "github.com/builderjarvis/create-go-app/features/ptr"
	_ "github.com/builderjarvis/create-go-app/features/retry"
	_ "github.com/builderjarvis/create-go-app/features/state"
	_ "github.com/builderjarvis/create-go-app/features/worker"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
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
	goVersion := strings.TrimPrefix(runtime.Version(), "go")

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
