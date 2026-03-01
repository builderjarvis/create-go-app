// Package scaffold provides the core scaffolding engine for project generation.
package scaffold

import (
	"fmt"
	"sort"
)

// Feature is the interface every optional feature implements.
type Feature interface {
	// Name returns the unique identifier (e.g., "postgres", "docker").
	Name() string

	// Description returns a human-readable description for prompts.
	Description() string

	// Dependencies returns names of features this feature requires.
	Dependencies() []string

	// Conflicts returns names of features that cannot coexist.
	Conflicts() []string

	// Install writes files and injects content into the project.
	Install(ctx *Context) error
}

// registry holds all registered features.
var registry = map[string]Feature{}

// Register adds a feature to the global registry. Called from feature init() functions.
func Register(f Feature) {
	registry[f.Name()] = f
}

// Get returns a feature by name.
func Get(name string) (Feature, bool) {
	f, ok := registry[name]
	return f, ok
}

// All returns all registered features sorted by name.
func All() []Feature {
	features := make([]Feature, 0, len(registry))
	for _, f := range registry {
		features = append(features, f)
	}
	sort.Slice(features, func(i, j int) bool {
		return features[i].Name() < features[j].Name()
	})
	return features
}

// AllNames returns all registered feature names sorted.
func AllNames() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// Resolve takes user-selected feature names, expands dependencies,
// checks conflicts, and returns features in topological order.
func Resolve(selected []string) ([]Feature, error) {
	// Expand dependencies via BFS.
	all := map[string]bool{}
	queue := append([]string{}, selected...)

	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]

		if all[name] {
			continue
		}

		f, ok := Get(name)
		if !ok {
			return nil, fmt.Errorf("unknown feature: %s", name)
		}

		all[name] = true

		for _, dep := range f.Dependencies() {
			if !all[dep] {
				queue = append(queue, dep)
			}
		}
	}

	// Check conflicts.
	for name := range all {
		f, _ := Get(name)
		for _, conflict := range f.Conflicts() {
			if all[conflict] {
				return nil, fmt.Errorf("feature %q conflicts with %q", name, conflict)
			}
		}
	}

	// Topological sort (Kahn's algorithm).
	return topoSort(all)
}

// topoSort returns features in dependency order using Kahn's algorithm.
func topoSort(names map[string]bool) ([]Feature, error) {
	// Build adjacency: feature -> features that depend on it.
	inDegree := map[string]int{}
	dependents := map[string][]string{}

	for name := range names {
		inDegree[name] = 0
	}

	for name := range names {
		f, _ := Get(name)
		for _, dep := range f.Dependencies() {
			if names[dep] {
				inDegree[name]++
				dependents[dep] = append(dependents[dep], name)
			}
		}
	}

	// Seed queue with zero in-degree nodes.
	var queue []string
	for name, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, name)
		}
	}
	sort.Strings(queue) // deterministic order

	var result []Feature
	for len(queue) > 0 {
		name := queue[0]
		queue = queue[1:]

		f, _ := Get(name)
		result = append(result, f)

		for _, dep := range dependents[name] {
			inDegree[dep]--
			if inDegree[dep] == 0 {
				queue = append(queue, dep)
			}
		}
		sort.Strings(queue)
	}

	if len(result) != len(names) {
		return nil, fmt.Errorf("circular dependency detected")
	}

	return result, nil
}
