package scaffold_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/builderjarvis/create-go-app/scaffold"

	// Register all features.
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

func TestGenerate_AllFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()

	cfg := scaffold.Config{
		ProjectName: "test-all",
		ModulePath:  "github.com/test/test-all",
		GoVersion:   goVersion(t),
		Features:    scaffold.AllNames(),
		OutputDir:   dir,
	}

	if err := scaffold.Generate(cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	projectDir := filepath.Join(dir, "test-all")
	assertGoBuild(t, projectDir)
}

func TestGenerate_NoFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()

	cfg := scaffold.Config{
		ProjectName: "test-bare",
		ModulePath:  "github.com/test/test-bare",
		GoVersion:   goVersion(t),
		Features:    nil,
		OutputDir:   dir,
	}

	if err := scaffold.Generate(cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	projectDir := filepath.Join(dir, "test-bare")
	assertGoBuild(t, projectDir)
}

func goVersion(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		t.Fatalf("go env GOVERSION: %v", err)
	}
	v := string(out)
	// Strip "go" prefix and newline.
	v = v[2:]
	if v[len(v)-1] == '\n' {
		v = v[:len(v)-1]
	}
	return v
}

func assertGoBuild(t *testing.T, dir string) {
	t.Helper()

	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("project dir does not exist: %v", err)
	}

	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go build failed:\n%s\n%v", out, err)
	}
}
