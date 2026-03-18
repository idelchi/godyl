package defaults_test

import (
	"slices"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"

	"github.com/idelchi/godyl/internal/defaults"
)

func TestPickHappyPath(t *testing.T) {
	t.Parallel()

	yaml := []byte(heredoc.Doc(`
		default:
		  name: default
		  output: /usr/local/bin

		other:
		  name: other
		  output: /opt/bin
	`))

	d, err := defaults.NewDefaultsFromBytes(yaml)
	if err != nil {
		t.Fatalf("NewDefaultsFromBytes() unexpected error: %v", err)
	}

	got, err := d.Pick("default", "other")
	if err != nil {
		t.Fatalf("Pick() unexpected error: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("Pick() returned %d items, want 2", len(got))
	}
}

func TestPickMissingName(t *testing.T) {
	t.Parallel()

	yaml := []byte(heredoc.Doc(`
		default:
		  name: default
		  output: /usr/local/bin
	`))

	d, err := defaults.NewDefaultsFromBytes(yaml)
	if err != nil {
		t.Fatalf("NewDefaultsFromBytes() unexpected error: %v", err)
	}

	_, err = d.Pick("nonexistent")
	if err == nil {
		t.Fatal("Pick() expected error for missing name, got nil")
	}
}

func TestNewDefaultsFromBytesInvalidYAML(t *testing.T) {
	t.Parallel()

	invalid := []byte(`{this is: [not valid yaml`)

	_, err := defaults.NewDefaultsFromBytes(invalid)
	if err == nil {
		t.Fatal("NewDefaultsFromBytes() expected error for invalid YAML, got nil")
	}
}

func TestBuildGraph(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		yaml       []byte
		wantErrMsg string
		wantNodes  []string // expected nodes in topological order (parents-first), checked when no error
	}{
		{
			name: "linear inheritance",
			yaml: []byte(heredoc.Doc(`
				default:
				  name: default
				  output: /usr/local/bin

				child:
				  name: child
				  inherit: default
			`)),
			wantNodes: []string{"default", "child"},
		},
		{
			name: "cycle a inherits b, b inherits a",
			yaml: []byte(heredoc.Doc(`
				a:
				  name: a
				  inherit: b

				b:
				  name: b
				  inherit: a
			`)),
			wantErrMsg: "cycle detected",
		},
		{
			name: "missing parent",
			yaml: []byte(heredoc.Doc(`
				child:
				  name: child
				  inherit: nonexistent
			`)),
			wantErrMsg: "undefined parent",
		},
		{
			name: "diamond: child has two parents sharing a grandparent",
			yaml: []byte(heredoc.Doc(`
				grandparent:
				  name: grandparent
				  output: /grand/bin

				parent1:
				  name: parent1
				  inherit: grandparent

				parent2:
				  name: parent2
				  inherit: grandparent

				child:
				  name: child
				  inherit:
				    - parent1
				    - parent2
			`)),
			wantNodes: []string{"grandparent", "parent1", "parent2", "child"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d, err := defaults.NewDefaultsFromBytes(tc.yaml)
			if err != nil {
				t.Fatalf("NewDefaultsFromBytes() unexpected error: %v", err)
			}

			g, err := d.BuildGraph()

			if tc.wantErrMsg != "" {
				if err == nil {
					t.Fatal("BuildGraph() expected error, got nil")
				}

				if !strings.Contains(err.Error(), tc.wantErrMsg) {
					t.Fatalf("BuildGraph() error = %q, want it to contain %q", err.Error(), tc.wantErrMsg)
				}

				return
			}

			if err != nil {
				t.Fatalf("BuildGraph() unexpected error: %v", err)
			}

			if g == nil {
				t.Fatal("BuildGraph() returned nil DAG without error")
			}

			topo := g.Topo()

			// Verify all expected nodes are present in the topo order.
			for _, node := range tc.wantNodes {
				if !slices.Contains(topo, node) {
					t.Fatalf("BuildGraph() Topo() = %v, missing expected node %q", topo, node)
				}
			}

			// Verify parents appear before children in topological order.
			// For each pair (node at index i, later at index j > i), node must not
			// be a descendant of later — i.e. later must not appear in node's ancestor chain.
			for i, node := range topo {
				for _, later := range topo[i+1:] {
					chain, err := g.Chain(node)
					if err != nil {
						t.Fatalf("Chain(%q) unexpected error: %v", node, err)
					}

					// chain includes node itself as the last element; ancestors are chain[:len(chain)-1].
					if slices.Contains(chain[:len(chain)-1], later) {
						t.Errorf(
							"topological order violated: %q appears before %q but %q is an ancestor of %q",
							node,
							later,
							later,
							node,
						)
					}
				}
			}
		})
	}
}

func TestResolveInheritance(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		yaml             []byte
		childName        string
		wantOutput       string
		parentName       string
		wantParentOutput string
	}{
		{
			name: "child inherits output from parent",
			yaml: []byte(heredoc.Doc(`
				default:
				  name: default
				  output: /usr/local/bin

				child:
				  name: child
				  inherit: default
			`)),
			childName:  "child",
			wantOutput: "/usr/local/bin",
		},
		{
			name: "child output overrides parent",
			yaml: []byte(heredoc.Doc(`
				default:
				  name: default
				  output: /usr/local/bin

				child:
				  name: child
				  output: /home/user/bin
				  inherit: default
			`)),
			childName:  "child",
			wantOutput: "/home/user/bin",
		},
		{
			name: "last declared parent wins when both set output",
			yaml: []byte(heredoc.Doc(`
				base1:
				  name: base1
				  output: /opt/bin

				base2:
				  name: base2
				  output: /usr/bin

				child:
				  name: child
				  inherit:
				    - base1
				    - base2
			`)),
			childName:  "child",
			wantOutput: "/usr/bin",
		},
		{
			name: "diamond: child inherits from two parents sharing a grandparent",
			yaml: []byte(heredoc.Doc(`
				grandparent:
				  name: grandparent
				  output: /grand/bin

				parent1:
				  name: parent1
				  inherit: grandparent

				parent2:
				  name: parent2
				  output: /parent2/bin
				  inherit: grandparent

				child:
				  name: child
				  inherit:
				    - parent1
				    - parent2
			`)),
			childName:        "child",
			wantOutput:       "/parent2/bin",
			parentName:       "parent1",
			wantParentOutput: "/grand/bin",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			d, err := defaults.NewDefaultsFromBytes(tc.yaml)
			if err != nil {
				t.Fatalf("NewDefaultsFromBytes() unexpected error: %v", err)
			}

			if err := d.ResolveInheritance(); err != nil {
				t.Fatalf("ResolveInheritance() unexpected error: %v", err)
			}

			child := d.Get(tc.childName)
			if child == nil {
				t.Fatalf("Get(%q) returned nil after ResolveInheritance()", tc.childName)
			}

			if child.Output != tc.wantOutput {
				t.Errorf("child.Output = %q, want %q", child.Output, tc.wantOutput)
			}

			if tc.parentName != "" {
				parent := d.Get(tc.parentName)
				if parent == nil {
					t.Fatalf("Get(%q) returned nil after ResolveInheritance()", tc.parentName)
				}

				if parent.Output != tc.wantParentOutput {
					t.Errorf("%s.Output = %q, want %q", tc.parentName, parent.Output, tc.wantParentOutput)
				}
			}
		})
	}
}
