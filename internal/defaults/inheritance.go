package defaults

import (
	"fmt"

	"github.com/idelchi/godyl/internal/debug"
	"github.com/idelchi/godyl/pkg/dag"
)

// BuildGraph builds a directed acyclic graph (DAG) from the defaults.
func (d Defaults) BuildGraph() (*dag.DAG[string], error) {
	// 1. collect node IDs (= tool names)
	nodes := make([]string, 0, len(d))
	for name := range d {
		nodes = append(nodes, name)
	}

	// 2. parent-lookup function required by dag.Build
	parentFn := func(name string) []string {
		t := d[name]
		if t == nil {
			return nil
		}

		if t.Inherit == nil {
			return nil
		}

		return *t.Inherit
	}

	// 3. Build and validate the DAG
	return dag.Build(nodes, parentFn)
}

// ResolveInheritance validates inheritance and mutates every *tool.Tool so that all
// parents' fields are applied.
func (d Defaults) ResolveInheritance() error {
	// Build the dependency DAG once.
	graph, err := d.BuildGraph()
	if err != nil {
		return fmt.Errorf("building inheritance tree for defaults: %w", err) // bad parent name or cycle
	}

	// Iterate in parents-first order so parents are final
	// before any child is processed.
	for _, name := range graph.Topo() {
		t := d[name]

		if t.Inherit == nil {
			debug.Debug("No inheritance for %q", name)

			continue
		}

		debug.Debug("Processing %q with inheritance: %v", name, t.Inherit)

		for _, p := range *t.Inherit { // direct parents, in declared order
			debug.Debug("constructing merge %q -> %q", name, p)

			if err := t.MergeInto(d[p]); err != nil {
				return fmt.Errorf("merging %q into %q: %w", p, name, err)
			}
		}
	}

	return nil
}
