package result_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/idelchi/godyl/internal/tools/result"
)

func TestResultStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		r           result.Result
		wantOK      bool
		wantSkipped bool
		wantFailed  bool
	}{
		{
			name:        "WithOK is only OK",
			r:           result.WithOK("success"),
			wantOK:      true,
			wantSkipped: false,
			wantFailed:  false,
		},
		{
			name:        "WithSkipped is only Skipped",
			r:           result.WithSkipped("skipping"),
			wantOK:      false,
			wantSkipped: true,
			wantFailed:  false,
		},
		{
			name:        "WithFailed is only Failed",
			r:           result.WithFailed("oops"),
			wantOK:      false,
			wantSkipped: false,
			wantFailed:  true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if got := tc.r.IsOK(); got != tc.wantOK {
				t.Errorf("IsOK() = %v, want %v", got, tc.wantOK)
			}

			if got := tc.r.IsSkipped(); got != tc.wantSkipped {
				t.Errorf("IsSkipped() = %v, want %v", got, tc.wantSkipped)
			}

			if got := tc.r.IsFailed(); got != tc.wantFailed {
				t.Errorf("IsFailed() = %v, want %v", got, tc.wantFailed)
			}
		})
	}
}

func TestResultWrapped(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		base        result.Result
		wrapMsg     string
		checkStatus func(result.Result) bool
		statusLabel string
	}{
		{
			name:        "Wrapped OK preserves OK status",
			base:        result.WithOK("original"),
			wrapMsg:     "extra context",
			checkStatus: result.Result.IsOK,
			statusLabel: "IsOK",
		},
		{
			name:        "Wrapped Skipped preserves Skipped status",
			base:        result.WithSkipped("original"),
			wrapMsg:     "skip reason",
			checkStatus: result.Result.IsSkipped,
			statusLabel: "IsSkipped",
		},
		{
			name:        "Wrapped Failed preserves Failed status",
			base:        result.WithFailed("original"),
			wrapMsg:     "more detail",
			checkStatus: result.Result.IsFailed,
			statusLabel: "IsFailed",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wrapped := tc.base.Wrapped(tc.wrapMsg)

			if !tc.checkStatus(wrapped) {
				t.Errorf("%s() = false after Wrapped, want true", tc.statusLabel)
			}

			if !strings.Contains(wrapped.Message, tc.wrapMsg) {
				t.Errorf("Wrapped(%q).Message = %q, want it to contain %q", tc.wrapMsg, wrapped.Message, tc.wrapMsg)
			}

			// The original message must also be present (Wrapped appends, not replaces)
			if !strings.Contains(wrapped.Message, tc.base.Message) {
				t.Errorf(
					"Wrapped(%q).Message = %q, want it to still contain original %q",
					tc.wrapMsg,
					wrapped.Message,
					tc.base.Message,
				)
			}
		})
	}
}

func TestResultWrap(t *testing.T) {
	t.Parallel()

	sentinel := errors.New("sentinel error")

	tests := []struct {
		name           string
		base           result.Result
		err            error
		wantAsErr      bool
		checkUnwrap    bool
		checkUnchanged bool
		unwrapSentinel error
	}{
		{
			name:      "Wrap with error on Failed makes AsError non-nil",
			base:      result.WithFailed("bad"),
			err:       errors.New("inner error"),
			wantAsErr: true,
		},
		{
			// Wrap(nil) is a no-op: the result is returned unchanged.
			// WithOK("x").Wrap(nil) stays OK, so AsError returns nil — distinguishable
			// from a Failed result which would return non-nil even without Wrap.
			name:           "Wrap(nil) on OK is a no-op; AsError stays nil",
			base:           result.WithOK("x"),
			err:            nil,
			wantAsErr:      false,
			checkUnchanged: true,
		},
		{
			name:      "Wrap with error on OK does not change AsError to non-nil",
			base:      result.WithOK("fine"),
			err:       errors.New("wrapped but ok"),
			wantAsErr: false,
		},
		{
			name:           "Wrap preserves inner error via errors.Is",
			base:           result.WithFailed("bad"),
			err:            sentinel,
			wantAsErr:      true,
			checkUnwrap:    true,
			unwrapSentinel: sentinel,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			wrapped := tc.base.Wrap(tc.err)

			gotAsErr := wrapped.AsError()
			if tc.wantAsErr && gotAsErr == nil {
				t.Errorf("AsError() = nil, want non-nil error")
			}

			if !tc.wantAsErr && gotAsErr != nil {
				t.Errorf("AsError() = %v, want nil", gotAsErr)
			}

			if tc.checkUnwrap && !errors.Is(wrapped, tc.unwrapSentinel) {
				t.Errorf("errors.Is(result, sentinel) = false, want true; Unwrap should expose inner error")
			}

			if tc.checkUnchanged && wrapped.Message != tc.base.Message {
				t.Errorf("Wrap(nil).Message = %q, want unchanged %q", wrapped.Message, tc.base.Message)
			}
		})
	}
}

func TestResultAsError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		r       result.Result
		wantNil bool
	}{
		{
			name:    "Failed returns non-nil error",
			r:       result.WithFailed("it broke"),
			wantNil: false,
		},
		{
			// AsError only returns an error for Failed status; OK and Skipped are not errors.
			name:    "non-Failed returns nil error",
			r:       result.WithOK("all good"),
			wantNil: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.r.AsError()

			if tc.wantNil && err != nil {
				t.Errorf("AsError() = %v, want nil", err)
			}

			if !tc.wantNil && err == nil {
				t.Errorf("AsError() = nil, want non-nil error")
			}
		})
	}
}

func TestResultAsErrorMessage(t *testing.T) {
	t.Parallel()

	const msg = "something went wrong"

	r := result.WithFailed(msg)

	err := r.AsError()
	if err == nil {
		t.Fatal("AsError() = nil, want non-nil for Failed result")
	}

	if !strings.Contains(err.Error(), msg) {
		t.Errorf("AsError().Error() = %q, want it to contain %q", err.Error(), msg)
	}
}

func TestResultAsErrorJoined(t *testing.T) {
	t.Parallel()

	err1 := errors.New("first error")
	err2 := errors.New("second error")

	// WithFailed joins multiple errors via errors.Join.
	r := result.WithFailed("multi", err1, err2)

	if !errors.Is(r, err1) {
		t.Errorf("errors.Is(result, err1) = false, want true")
	}

	if !errors.Is(r, err2) {
		t.Errorf("errors.Is(result, err2) = false, want true")
	}
}

// TestResultAsErrorVariants is a table-driven test covering the AsError()
// message format and the joined-errors case in a single, consolidated test.
// It complements the individual TestResultAsErrorMessage and
// TestResultAsErrorJoined tests above.
func TestResultAsErrorVariants(t *testing.T) {
	t.Parallel()

	err1 := errors.New("first error")
	err2 := errors.New("second error")

	tests := []struct {
		name         string
		r            result.Result
		wantNilErr   bool
		wantContains string
		wantIs       []error // errors that should be found via errors.Is
	}{
		{
			name:         "WithFailed no inner error: message appears in error string",
			r:            result.WithFailed("download failed"),
			wantContains: "download failed",
		},
		{
			name:         "WithFailed with one inner error: inner error text appears",
			r:            result.WithFailed("wrap test", errors.New("inner")),
			wantContains: "inner",
		},
		{
			name:   "WithFailed with two joined errors: both retrievable via errors.Is",
			r:      result.WithFailed("multi-error", err1, err2),
			wantIs: []error{err1, err2},
		},
		{
			name:       "WithOK: AsError returns nil",
			r:          result.WithOK("all fine"),
			wantNilErr: true,
		},
		{
			name:       "WithSkipped: AsError returns nil",
			r:          result.WithSkipped("intentional skip"),
			wantNilErr: true,
		},
		{
			// Wrap(nil) on a Failed result must not change AsError to nil
			// (the Failed status drives AsError, not the inner error).
			name:         "Failed Wrap(nil): AsError still non-nil because status is Failed",
			r:            result.WithFailed("base").Wrap(nil),
			wantContains: "base",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.r.AsError()

			if tc.wantNilErr {
				if err != nil {
					t.Errorf("AsError() = %v, want nil", err)
				}

				return
			}

			if err == nil {
				t.Fatalf("AsError() = nil, want non-nil error")
			}

			if tc.wantContains != "" && !strings.Contains(err.Error(), tc.wantContains) {
				t.Errorf("AsError().Error() = %q, want it to contain %q", err.Error(), tc.wantContains)
			}

			for _, want := range tc.wantIs {
				if !errors.Is(tc.r, want) {
					t.Errorf("errors.Is(result, %v) = false, want true", want)
				}
			}
		})
	}
}

func TestResultError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		r           result.Result
		wantContain []string
	}{
		{
			name:        "OK without inner error contains status and message",
			r:           result.WithOK("done"),
			wantContain: []string{"ok", "done"},
		},
		{
			name:        "Failed without inner error contains status and message",
			r:           result.WithFailed("broken"),
			wantContain: []string{"failed", "broken"},
		},
		{
			name:        "Failed with inner error includes all three",
			r:           result.WithFailed("oops", errors.New("root cause")),
			wantContain: []string{"failed", "oops", "root cause"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.r.Error()

			for _, want := range tc.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("Error() = %q, want it to contain %q", got, want)
				}
			}
		})
	}
}

func TestResultUnwrap(t *testing.T) {
	t.Parallel()

	t.Run("no inner error returns nil", func(t *testing.T) {
		t.Parallel()

		r := result.WithFailed("no inner")

		if got := r.Unwrap(); got != nil {
			t.Errorf("Unwrap() = %v, want nil", got)
		}
	})

	t.Run("with inner error returns it", func(t *testing.T) {
		t.Parallel()

		inner := errors.New("inner")
		r := result.WithFailed("msg", inner)

		got := r.Unwrap()
		if !errors.Is(got, inner) {
			t.Errorf("Unwrap() = %v, want errors.Is(inner)", got)
		}
	})

	t.Run("Wrap sets inner error accessible via Unwrap", func(t *testing.T) {
		t.Parallel()

		sentinel := errors.New("sentinel")
		r := result.WithFailed("base").Wrap(sentinel)

		got := r.Unwrap()
		if got == nil {
			t.Fatal("Unwrap() = nil after Wrap, want non-nil")
		}

		if !errors.Is(got, sentinel) {
			t.Errorf("Unwrap() does not wrap sentinel: %v", got)
		}
	})
}

func TestResultNew(t *testing.T) {
	t.Parallel()

	r := result.New("hello", result.OK)

	if !r.IsOK() {
		t.Error("New(OK).IsOK() = false, want true")
	}

	if r.Message != "hello" {
		t.Errorf("New(OK).Message = %q, want %q", r.Message, "hello")
	}
}

func TestResultString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		r           result.Result
		wantContain []string
	}{
		{
			name:        "OK result string contains status and message",
			r:           result.WithOK("installed"),
			wantContain: []string{"ok", "installed"},
		},
		{
			name:        "Skipped result string contains status and message",
			r:           result.WithSkipped("already exists"),
			wantContain: []string{"skipped", "already exists"},
		},
		{
			name:        "Failed result string contains status and message",
			r:           result.WithFailed("download error"),
			wantContain: []string{"failed", "download error"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			s := tc.r.String()

			for _, want := range tc.wantContain {
				if !strings.Contains(s, want) {
					t.Errorf("String() = %q, want it to contain %q", s, want)
				}
			}
		})
	}
}
