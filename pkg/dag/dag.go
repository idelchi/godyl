// Package dag implements a generic directed acyclic graph with cycle detection,
// topological sorting, and ancestor chain computation.
package dag

import (
	"errors"
	"fmt"
	"slices"
	"strings"
)

// ErrDag is returned when a DAG operation fails.
var ErrDag = errors.New("dag error")

// DAG is an immutable, validated representation of a dependency graph.
// The type parameter K must be comparable so it can act as a map key.
type DAG[K comparable] struct {
	// child → direct parents
	parents map[K][]K
	// memoised Chain results
	ancestorCache map[K][]K
	// every node, parents-first
	topo []K
}

// CycleError represents a cycle in the graph with the path that forms the cycle.
type CycleError[K comparable] struct {
	// The node where the cycle was detected
	Node K
	// The path that forms the cycle
	Path []K
}

// Error implements the error interface for CycleError.
func (e *CycleError[K]) Error() string {
	if len(e.Path) == 0 {
		return fmt.Sprintf("cycle detected at %v", e.Node)
	}

	// Find the node's position in the path (it should appear twice)
	startIdx := -1

	for i, node := range e.Path {
		if node == e.Node {
			startIdx = i

			break
		}
	}

	// If we found the start node in the path, extract the cycle
	if startIdx >= 0 {
		// Extract just the cycle portion
		cyclePath := append(slices.Clone(e.Path[startIdx:]), e.Node)
		parts := make([]string, len(cyclePath))

		for i, node := range cyclePath {
			parts[i] = fmt.Sprintf("%v", node)
		}

		return fmt.Sprintf("%v: cycle detected: %s", ErrDag, strings.Join(parts, " -> "))
	}

	// Fallback if we couldn't find a clean cycle
	parts := make([]string, len(e.Path))
	for i, node := range e.Path {
		parts[i] = fmt.Sprintf("%v", node)
	}

	return fmt.Sprintf("%v: cycle detected: %s -> %v", ErrDag, strings.Join(parts, " -> "), e.Node)
}

// Build constructs the graph and performs **two** validations:
//
//  1. Every parent returned by parentsFn(node) must itself be present in
//     the `nodes` slice.
//  2. The graph must be acyclic.
//
// On success a *DAG is returned; otherwise an error explains the problem.
func Build[K comparable](
	nodes []K,
	parentsFn func(K) []K,
) (*DAG[K], error) {
	graph := &DAG[K]{
		parents:       make(map[K][]K, len(nodes)),
		ancestorCache: make(map[K][]K, len(nodes)),
	}

	// Populate parent table (defensive copies).
	for _, n := range nodes {
		p := parentsFn(n)
		cp := make([]K, len(p))
		copy(cp, p)

		graph.parents[n] = cp
	}

	// Quick lookup: does node exist?
	set := make(map[K]struct{}, len(nodes))
	for _, n := range nodes {
		set[n] = struct{}{}
	}

	// DFS for validation + topological sort.
	color := make(map[K]uint8, len(nodes)) // 0 white,1 gray,2 black

	path := make(map[K]int, len(nodes)) // Track position in current path

	var currentPath []K // Current DFS path

	var dfs func(K) error

	dfs = func(node K) error {
		switch color[node] {
		case 1: // gray - we found a cycle
			// Reconstruct the cycle path
			cycleStart := path[node]
			cyclePath := make([]K, len(currentPath)-cycleStart)
			copy(cyclePath, currentPath[cycleStart:])

			return &CycleError[K]{
				Node: node,
				Path: cyclePath,
			}
		case 2: //nolint:mnd // Constant 2 represents DFS black state
			// black - already processed
			return nil
		}

		color[node] = 1 // mark as gray (being processed)
		path[node] = len(currentPath)
		currentPath = append(currentPath, node)

		for _, p := range graph.parents[node] {
			if _, ok := set[p]; !ok {
				return fmt.Errorf("%w: undefined parent %v referenced by %v", ErrDag, p, node)
			}

			if err := dfs(p); err != nil {
				return err
			}
		}

		color[node] = 2 // mark as black (fully processed)

		currentPath = currentPath[:len(currentPath)-1] // remove from current path

		graph.topo = append(graph.topo, node)

		return nil
	}

	for _, node := range nodes {
		if color[node] == 0 {
			currentPath = nil // Reset path for each root node

			if err := dfs(node); err != nil {
				return nil, err
			}
		}
	}

	return graph, nil
}

// Topo returns every node exactly once in "parents first" order.
func (g *DAG[K]) Topo() []K {
	out := make([]K, len(g.topo))
	copy(out, g.topo)

	return out
}

// Chain returns the linearised chain `[grandParent … parent, node]`.
//
// The slice is newly allocated on each call; modifying it will not affect the
// DAG or future calls.
func (g *DAG[K]) Chain(node K) ([]K, error) {
	if c, ok := g.ancestorCache[node]; ok {
		cp := make([]K, len(c))

		copy(cp, c)

		return cp, nil
	}

	parents, ok := g.parents[node]
	if !ok {
		return nil, fmt.Errorf("%w: unknown node %v", ErrDag, node)
	}

	var out []K

	for _, parent := range parents {
		sub, err := g.Chain(parent)
		if err != nil {
			return nil, err
		}

		out = append(out, sub...)
	}

	out = append(out, node)

	g.ancestorCache[node] = out

	cp := make([]K, len(out))

	copy(cp, out)

	return cp, nil
}

// Condense returns a copy of `chain` with only *adjacent* duplicates removed.
//
//	["base","base","custom"]        → ["base","custom"]
//	["base","custom","base"]        → (unchanged)
//	["a","a","b","b","b","c","c"]   → ["a","b","c"]
//
// The input slice is never modified.
func Condense[K comparable](chain []K) []K {
	if len(chain) == 0 {
		return nil
	}

	out := make([]K, 0, len(chain))

	prev := chain[0]

	out = append(out, prev)

	for i := 1; i < len(chain); i++ {
		if chain[i] == prev {
			continue // adjacent duplicate → skip
		}

		prev = chain[i]

		out = append(out, prev)
	}

	return out
}
