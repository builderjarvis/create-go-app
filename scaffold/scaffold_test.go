package scaffold_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/builderjarvis/create-go-app/scaffold"

	// Register all features.
	_ "github.com/builderjarvis/create-go-app/features/ci"
	_ "github.com/builderjarvis/create-go-app/features/docker"
	_ "github.com/builderjarvis/create-go-app/features/httpclient"
	_ "github.com/builderjarvis/create-go-app/features/postgres"
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

func TestGenerate_PostgresOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()

	cfg := scaffold.Config{
		ProjectName: "test-pg",
		ModulePath:  "github.com/test/test-pg",
		GoVersion:   goVersion(t),
		Features:    []string{"postgres"},
		OutputDir:   dir,
	}

	if err := scaffold.Generate(cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	projectDir := filepath.Join(dir, "test-pg")
	assertGoBuild(t, projectDir)
}

func TestGenerate_DockerOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()

	cfg := scaffold.Config{
		ProjectName: "test-docker",
		ModulePath:  "github.com/test/test-docker",
		GoVersion:   goVersion(t),
		Features:    []string{"docker"},
		OutputDir:   dir,
	}

	if err := scaffold.Generate(cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	projectDir := filepath.Join(dir, "test-docker")
	assertGoBuild(t, projectDir)

	// Verify compose.yaml exists but has no postgres service
	compose, err := os.ReadFile(filepath.Join(projectDir, "compose.yaml"))
	if err != nil {
		t.Fatalf("reading compose.yaml: %v", err)
	}
	if strings.Contains(string(compose), "postgres") {
		t.Error("compose.yaml should not contain postgres service when postgres feature is not selected")
	}
}

func TestGenerate_HTTPClientOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()

	cfg := scaffold.Config{
		ProjectName: "test-http",
		ModulePath:  "github.com/test/test-http",
		GoVersion:   goVersion(t),
		Features:    []string{"httpclient"},
		OutputDir:   dir,
	}

	if err := scaffold.Generate(cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	projectDir := filepath.Join(dir, "test-http")
	assertGoBuild(t, projectDir)
}

func TestGenerate_CIOnly(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()

	cfg := scaffold.Config{
		ProjectName: "test-ci",
		ModulePath:  "github.com/test/test-ci",
		GoVersion:   goVersion(t),
		Features:    []string{"ci"},
		OutputDir:   dir,
	}

	if err := scaffold.Generate(cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	projectDir := filepath.Join(dir, "test-ci")
	assertGoBuild(t, projectDir)

	// Verify CI workflow exists
	if _, err := os.Stat(filepath.Join(projectDir, ".github", "workflows", "ci.yaml")); err != nil {
		t.Error("expected .github/workflows/ci.yaml to exist")
	}
}

func TestGenerate_PostgresDocker(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	dir := t.TempDir()

	cfg := scaffold.Config{
		ProjectName: "test-pgdocker",
		ModulePath:  "github.com/test/test-pgdocker",
		GoVersion:   goVersion(t),
		Features:    []string{"postgres", "docker"},
		OutputDir:   dir,
	}

	if err := scaffold.Generate(cfg); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	projectDir := filepath.Join(dir, "test-pgdocker")
	assertGoBuild(t, projectDir)

	// Verify compose.yaml has postgres service
	compose, err := os.ReadFile(filepath.Join(projectDir, "compose.yaml"))
	if err != nil {
		t.Fatalf("reading compose.yaml: %v", err)
	}
	if !strings.Contains(string(compose), "postgres") {
		t.Error("compose.yaml should contain postgres service when both features are selected")
	}
}

func goVersion(t *testing.T) string {
	t.Helper()
	out, err := exec.Command("go", "env", "GOVERSION").Output()
	if err != nil {
		t.Fatalf("go env GOVERSION: %v", err)
	}
	v := string(out)
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
