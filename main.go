// Package main is the entry point for the create-go-app scaffolder.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"

	"github.com/builderjarvis/create-go-app/scaffold"

	// Register all features via blank imports.
	_ "github.com/builderjarvis/create-go-app/features/ci"
	_ "github.com/builderjarvis/create-go-app/features/docker"
	_ "github.com/builderjarvis/create-go-app/features/httpclient"
	_ "github.com/builderjarvis/create-go-app/features/postgres"
)

const banner = `
                        _                                           
  ___ _ __ ___  __ _| |_ ___        __ _  ___         __ _ _ __  _ __  
 / __| '__/ _ \/ _` + "`" + ` | __/ _ \_____ / _` + "`" + ` |/ _ \ _____ / _` + "`" + ` | '_ \| '_ \ 
| (__| | |  __/ (_| | ||  __/_____| (_| | (_) |_____| (_| | |_) | |_) |
 \___|_|  \___|\__,_|\__\___|      \__, |\___/       \__,_| .__/| .__/ 
                                   |___/                   |_|   |_|    
`

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("63")).
		Render(banner)
	fmt.Println(title)

	// Defaults.
	projectName := "my-app"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		projectName = args[0]
	}

	modulePath := "github.com/user/" + projectName
	goVersion := strings.TrimPrefix(runtime.Version(), "go")

	// Build feature options from registry.
	allFeatures := scaffold.All()
	featureOptions := make([]huh.Option[string], len(allFeatures))
	for i, f := range allFeatures {
		featureOptions[i] = huh.NewOption(
			fmt.Sprintf("%s (%s)", f.Name(), f.Description()),
			f.Name(),
		)
	}

	var selectedFeatures []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("What will your project be called?").
				Value(&projectName).
				Placeholder("my-app"),

			huh.NewInput().
				Title("What is your Go module path?").
				Value(&modulePath).
				Placeholder("github.com/user/my-app"),

			huh.NewMultiSelect[string]().
				Title("Which features would you like to include?").
				Options(featureOptions...).
				Value(&selectedFeatures),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return fmt.Errorf("form: %w", err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("getting working directory: %w", err)
	}

	cfg := scaffold.Config{
		ProjectName: projectName,
		ModulePath:  modulePath,
		GoVersion:   goVersion,
		Features:    selectedFeatures,
		OutputDir:   cwd,
	}

	// Summary.
	fmt.Println()
	summaryStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("63"))
	fmt.Println(summaryStyle.Render("  Creating project:"))
	fmt.Printf("    Name:     %s\n", projectName)
	fmt.Printf("    Module:   %s\n", modulePath)
	fmt.Printf("    Path:     %s\n", filepath.Join(cwd, projectName))
	if len(selectedFeatures) > 0 {
		fmt.Printf("    Features: %s\n", strings.Join(selectedFeatures, ", "))
	} else {
		fmt.Printf("    Features: (base only)\n")
	}
	fmt.Println()

	if err := scaffold.Generate(cfg); err != nil {
		return err
	}

	fmt.Println()
	success := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	fmt.Println(success.Render("  ✓ Project created successfully!"))
	fmt.Println()
	fmt.Printf("    cd %s\n", projectName)
	fmt.Println("    cp .env.example .env")
	fmt.Println("    make dev")
	fmt.Println()

	return nil
}
