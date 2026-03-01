// Package main is the entry point for the create-go-app scaffolder.
package main

import (
	"fmt"
	"os"

	"github.com/builderjarvis/create-go-app/cli"

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
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
