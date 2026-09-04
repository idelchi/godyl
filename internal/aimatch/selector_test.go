//nolint:testpackage // Tests exercise the unexported exact-response validator.
package aimatch

import "testing"

func TestExactCandidate(t *testing.T) {
	t.Parallel()

	candidates := []string{"tool-linux-amd64.tar.gz", "tool-darwin-arm64.tar.gz"}

	selected, err := exactCandidate("```\ntool-linux-amd64.tar.gz\n```", candidates)
	if err != nil {
		t.Fatalf("exactCandidate() unexpected error: %v", err)
	}

	if selected != candidates[0] {
		t.Errorf("exactCandidate() = %q, want %q", selected, candidates[0])
	}
}

func TestExactCandidateRejectsExplanation(t *testing.T) {
	t.Parallel()

	candidates := []string{"tool-linux-amd64.tar.gz"}

	if _, err := exactCandidate("Use tool-linux-amd64.tar.gz", candidates); err == nil {
		t.Fatal("exactCandidate() expected an error for a non-exact response")
	}
}
