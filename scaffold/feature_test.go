package scaffold

import (
	"testing"
)

// stubFeature is a minimal Feature for testing.
type stubFeature struct {
	name      string
	deps      []string
	conflicts []string
}

func (f *stubFeature) Name() string           { return f.name }
func (f *stubFeature) Description() string    { return f.name }
func (f *stubFeature) Dependencies() []string { return f.deps }
func (f *stubFeature) Conflicts() []string    { return f.conflicts }
func (f *stubFeature) Install(*Context) error { return nil }

func TestResolve_ExpandsDeps(t *testing.T) {
	orig := registry
	defer func() { registry = orig }()
	registry = map[string]Feature{}

	Register(&stubFeature{name: "a", deps: []string{"b"}})
	Register(&stubFeature{name: "b"})
	Register(&stubFeature{name: "c"})

	features, err := Resolve([]string{"a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	names := make(map[string]bool)
	for _, f := range features {
		names[f.Name()] = true
	}

	if !names["a"] || !names["b"] {
		t.Errorf("expected a and b, got %v", names)
	}
	if names["c"] {
		t.Error("c should not be included")
	}
}

func TestResolve_TopoOrder(t *testing.T) {
	orig := registry
	defer func() { registry = orig }()
	registry = map[string]Feature{}

	Register(&stubFeature{name: "a", deps: []string{"b"}})
	Register(&stubFeature{name: "b"})

	features, err := Resolve([]string{"a"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(features) != 2 {
		t.Fatalf("expected 2 features, got %d", len(features))
	}
	if features[0].Name() != "b" || features[1].Name() != "a" {
		t.Errorf("expected [b, a], got [%s, %s]", features[0].Name(), features[1].Name())
	}
}

func TestResolve_Conflict(t *testing.T) {
	orig := registry
	defer func() { registry = orig }()
	registry = map[string]Feature{}

	Register(&stubFeature{name: "a", conflicts: []string{"b"}})
	Register(&stubFeature{name: "b"})

	_, err := Resolve([]string{"a", "b"})
	if err == nil {
		t.Fatal("expected conflict error")
	}
}

func TestResolve_UnknownFeature(t *testing.T) {
	orig := registry
	defer func() { registry = orig }()
	registry = map[string]Feature{}

	_, err := Resolve([]string{"nonexistent"})
	if err == nil {
		t.Fatal("expected error for unknown feature")
	}
}

func TestRegisterAndAll(t *testing.T) {
	orig := registry
	defer func() { registry = orig }()
	registry = map[string]Feature{}

	Register(&stubFeature{name: "z"})
	Register(&stubFeature{name: "a"})

	all := All()
	if len(all) != 2 {
		t.Fatalf("expected 2 features, got %d", len(all))
	}
	if all[0].Name() != "a" || all[1].Name() != "z" {
		t.Errorf("expected sorted [a, z], got [%s, %s]", all[0].Name(), all[1].Name())
	}
}
