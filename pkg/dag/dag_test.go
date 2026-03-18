package dag_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/idelchi/godyl/pkg/dag"
)

// parentMap is a test helper that builds a parentsFn from a map.
func parentMap(m map[string][]string) func(string) []string {
	return func(k string) []string {
		return m[k]
	}
}

func TestBuild(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		nodes         []string
		parents       map[string][]string
		wantError     bool
		wantCycleErr  bool
		wantErrSubstr string
	}{
		{
			name:  "linear A depends on B depends on C",
			nodes: []string{"A", "B", "C"},
			parents: map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": {},
			},
		},
		{
			name:  "diamond A depends on B and C both depend on D",
			nodes: []string{"A", "B", "C", "D"},
			parents: map[string][]string{
				"A": {"B", "C"},
				"B": {"D"},
				"C": {"D"},
				"D": {},
			},
		},
		{
			name:  "no edges",
			nodes: []string{"A", "B", "C"},
			parents: map[string][]string{
				"A": {},
				"B": {},
				"C": {},
			},
		},
		{
			name:  "cycle A depends on B depends on A",
			nodes: []string{"A", "B"},
			parents: map[string][]string{
				"A": {"B"},
				"B": {"A"},
			},
			wantError:     true,
			wantCycleErr:  true,
			wantErrSubstr: "cycle detected",
		},
		{
			name:  "self-cycle A depends on A",
			nodes: []string{"A"},
			parents: map[string][]string{
				"A": {"A"},
			},
			wantError:     true,
			wantCycleErr:  true,
			wantErrSubstr: "A",
		},
		{
			name:      "single node no parents",
			nodes:     []string{"A"},
			parents:   map[string][]string{"A": {}},
			wantError: false,
		},
		{
			name:  "undefined parent referenced by node",
			nodes: []string{"A"},
			parents: map[string][]string{
				"A": {"B"}, // B is not in nodes
			},
			wantError:     true,
			wantErrSubstr: "undefined parent",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			g, err := dag.Build(tc.nodes, parentMap(tc.parents))

			switch {
			case tc.wantError && err == nil:
				t.Fatal("Build() error = nil, want error")
			case tc.wantError:
				if tc.wantCycleErr {
					var cycleErr *dag.CycleError[string]
					if !errors.As(err, &cycleErr) {
						t.Errorf("Build() error type = %T, want *dag.CycleError[string]; err = %v", err, err)
					}
				}

				if tc.wantErrSubstr != "" && !strings.Contains(err.Error(), tc.wantErrSubstr) {
					t.Errorf("Build() error = %q, want it to contain %q", err.Error(), tc.wantErrSubstr)
				}

				return
			case err != nil:
				t.Fatalf("Build() unexpected error: %v", err)
			}

			if g == nil {
				t.Fatal("Build() returned nil DAG without error")
			}
		})
	}
}

func TestTopo(t *testing.T) {
	t.Parallel()

	// A depends on B, B depends on C. Topo order is parents-first: C, B, A.
	g, err := dag.Build(
		[]string{"A", "B", "C"},
		parentMap(map[string][]string{
			"A": {"B"},
			"B": {"C"},
			"C": {},
		}),
	)
	if err != nil {
		t.Fatalf("Build() unexpected error: %v", err)
	}

	got := g.Topo()
	want := []string{"C", "B", "A"}

	if !slices.Equal(got, want) {
		t.Errorf("Topo() = %v, want %v", got, want)
	}

	t.Run("four independent roots all appear in topo", func(t *testing.T) {
		t.Parallel()

		g2, err := dag.Build(
			[]string{"W", "X", "Y", "Z"},
			parentMap(map[string][]string{
				"W": {},
				"X": {},
				"Y": {},
				"Z": {},
			}),
		)
		if err != nil {
			t.Fatalf("Build() unexpected error: %v", err)
		}

		topo := g2.Topo()

		if len(topo) != 4 {
			t.Fatalf("Topo() len = %d, want 4; got %v", len(topo), topo)
		}

		sorted := slices.Clone(topo)
		slices.Sort(sorted)

		want := []string{"W", "X", "Y", "Z"}

		if !slices.Equal(sorted, want) {
			t.Errorf("Topo() sorted = %v, want %v", sorted, want)
		}
	})
}

func TestChain(t *testing.T) {
	t.Parallel()

	tests := []struct {
		node string
		want []string
	}{
		{node: "A", want: []string{"C", "B", "A"}},
		{node: "B", want: []string{"C", "B"}},
		{node: "C", want: []string{"C"}},
	}

	// Each subtest builds its own *DAG so that concurrent subtests do not
	// race on the shared ancestorCache inside a single DAG instance.
	for _, tc := range tests {
		t.Run(tc.node, func(t *testing.T) {
			t.Parallel()

			// A depends on B, B depends on C.
			g, err := dag.Build(
				[]string{"A", "B", "C"},
				parentMap(map[string][]string{
					"A": {"B"},
					"B": {"C"},
					"C": {},
				}),
			)
			if err != nil {
				t.Fatalf("Build() unexpected error: %v", err)
			}

			got, err := g.Chain(tc.node)
			if err != nil {
				t.Fatalf("Chain(%q) unexpected error: %v", tc.node, err)
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("Chain(%q) = %v, want %v", tc.node, got, tc.want)
			}
		})
	}

	t.Run("unknown", func(t *testing.T) {
		t.Parallel()

		g, err := dag.Build(
			[]string{"A", "B", "C"},
			parentMap(map[string][]string{
				"A": {"B"},
				"B": {"C"},
				"C": {},
			}),
		)
		if err != nil {
			t.Fatalf("Build() unexpected error: %v", err)
		}

		_, err = g.Chain("Z")
		if err == nil {
			t.Fatal("Chain(unknown node) should error")
		}
	})
}

func TestCondense(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "adjacent duplicates collapsed",
			input: []string{"A", "A", "B", "B", "C"},
			want:  []string{"A", "B", "C"},
		},
		{
			name:  "no duplicates unchanged",
			input: []string{"A", "B", "C"},
			want:  []string{"A", "B", "C"},
		},
		{
			name:  "non-adjacent duplicates preserved",
			input: []string{"A", "B", "A"},
			want:  []string{"A", "B", "A"},
		},
		{
			name:  "all same collapses to one",
			input: []string{"A", "A", "A"},
			want:  []string{"A"},
		},
		{
			name:  "nil input returns empty",
			input: nil,
			want:  nil,
		},
		{
			name:  "single element unchanged",
			input: []string{"A"},
			want:  []string{"A"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := dag.Condense(tc.input)

			// Use len-based comparison so that both nil and []string{} satisfy
			// an "empty" expectation without leaking nil-vs-empty distinctions.
			if len(tc.want) == 0 {
				if len(got) != 0 {
					t.Errorf("Condense() = %v, want empty", got)
				}

				return
			}

			if !slices.Equal(got, tc.want) {
				t.Errorf("Condense() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestBuildEmpty(t *testing.T) {
	t.Parallel()

	// Build with an empty nodes slice must succeed and produce a valid,
	// empty DAG (no error, non-nil return, zero-length Topo).
	g, err := dag.Build([]string{}, parentMap(map[string][]string{}))
	if err != nil {
		t.Fatalf("Build(empty): unexpected error: %v", err)
	}

	if g == nil {
		t.Fatal("Build(empty): returned nil DAG, want non-nil")
	}

	topo := g.Topo()
	if len(topo) != 0 {
		t.Errorf("Build(empty): Topo() = %v, want empty slice", topo)
	}
}

func TestChainDiamond(t *testing.T) {
	t.Parallel()

	// Diamond: A→B, A→C, B→D, C→D
	// In parent-map terms: A has parents [B, C], B has parent [D], C has parent [D], D has no parents.
	//
	// Chain walks each parent path in full, so the shared ancestor D will appear
	// once per path that includes it.  The resulting chain is [D B D C A].
	// Condense (which removes adjacent duplicates only) does not collapse it
	// further because the two D entries are non-adjacent.
	// The key invariants are:
	//   1. "A" is the final element.
	//   2. Every occurrence of D precedes the B/C entry in its respective sub-path.
	//   3. Condensing the chain yields [D B D C A] unchanged (no adjacent dups).
	g, err := dag.Build(
		[]string{"A", "B", "C", "D"},
		parentMap(map[string][]string{
			"A": {"B", "C"},
			"B": {"D"},
			"C": {"D"},
			"D": {},
		}),
	)
	if err != nil {
		t.Fatalf("Build(diamond): unexpected error: %v", err)
	}

	chain, err := g.Chain("A")
	if err != nil {
		t.Fatalf("Chain(\"A\"): unexpected error: %v", err)
	}

	// "A" must be the last element.
	if len(chain) == 0 {
		t.Fatal("Chain(\"A\"): returned empty slice")
	}

	if chain[len(chain)-1] != "A" {
		t.Errorf("Chain(\"A\"): last element = %q, want \"A\"", chain[len(chain)-1])
	}

	// The first occurrence of D must come before B (D is B's parent).
	// The second occurrence of D must come before C (D is C's parent).
	// We verify this by scanning for the sub-sequences [D … B] and [D … C].
	foundDBeforeB := false
	foundDBeforeC := false

	lastD := -1

	for i, n := range chain {
		switch n {
		case "D":
			lastD = i
		case "B":
			if lastD >= 0 && lastD < i {
				foundDBeforeB = true
			}
		case "C":
			if lastD >= 0 && lastD < i {
				foundDBeforeC = true
			}
		}
	}

	if !foundDBeforeB {
		t.Errorf("Chain(\"A\"): no occurrence of \"D\" before \"B\" in chain %v", chain)
	}

	if !foundDBeforeC {
		t.Errorf("Chain(\"A\"): no occurrence of \"D\" before \"C\" in chain %v", chain)
	}

	// Condense must not alter the chain because none of D's occurrences are adjacent.
	condensed := dag.Condense(chain)

	if !slices.Equal(condensed, chain) {
		t.Errorf("Condense(Chain(\"A\")) = %v, want unchanged %v", condensed, chain)
	}
}

func TestCondenseEmpty(t *testing.T) {
	t.Parallel()

	// Condense with a non-nil but empty slice must return nil (or empty),
	// matching the documented behaviour for len(chain)==0.
	got := dag.Condense([]string{})

	if len(got) != 0 {
		t.Errorf("Condense([]string{}): got %v, want empty", got)
	}
}
