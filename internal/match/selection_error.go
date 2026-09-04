package match

// SelectionError preserves the structured context for a deterministic asset
// matching failure.
type SelectionError struct {
	// Cause classifies and describes the matching failure.
	Cause error
	// Candidates contains every release asset considered by the matcher.
	Candidates []string
	// Results contains the results relevant to the failure. For an ambiguous
	// match, these are the equally ranked best candidates.
	Results Results
	// Requirements contains the effective platform and hints used for matching.
	Requirements Requirements
}

// Error returns the deterministic matching error.
func (e *SelectionError) Error() string {
	return e.Cause.Error()
}

// Unwrap exposes the matching error for errors.Is and errors.As.
func (e *SelectionError) Unwrap() error {
	return e.Cause
}
