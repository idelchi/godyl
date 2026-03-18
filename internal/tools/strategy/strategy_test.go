package strategy_test

import (
	"testing"

	"github.com/idelchi/godyl/internal/tools/strategy"
)

// mockTool implements strategy.Tool for use in tests.
type mockTool struct {
	exists         bool
	currentVersion string
	start          strategy.Strategy
	targetVersion  string
}

func (m mockTool) Exists() bool                   { return m.exists }
func (m mockTool) GetCurrentVersion() string      { return m.currentVersion }
func (m mockTool) GetStrategy() strategy.Strategy { return m.start }
func (m mockTool) GetTargetVersion() string       { return m.targetVersion }

func TestSync(t *testing.T) {
	t.Parallel()

	// s is a zero-value Strategy; the receiver is unused by Sync — dispatch
	// reads t.GetStrategy() from the tool.
	var s strategy.Strategy

	tests := []struct {
		name        string
		tool        mockTool
		wantOK      bool
		wantSkipped bool
		wantFailed  bool
	}{
		// None strategy
		{
			name: "None + tool exists → Skipped",
			tool: mockTool{
				exists: true,
				start:  strategy.None,
			},
			wantSkipped: true,
		},
		{
			name: "None + tool does not exist → OK",
			tool: mockTool{
				exists: false,
				start:  strategy.None,
			},
			wantOK: true,
		},

		// Sync strategy
		{
			name: "Sync + exists + versions match → Skipped",
			tool: mockTool{
				exists:         true,
				start:          strategy.Sync,
				currentVersion: "v1.0.0",
				targetVersion:  "v1.0.0",
			},
			wantSkipped: true,
		},
		{
			name: "Sync + exists + versions differ → OK",
			tool: mockTool{
				exists:         true,
				start:          strategy.Sync,
				currentVersion: "v1.0.0",
				targetVersion:  "v1.1.0",
			},
			wantOK: true,
		},
		{
			name: "Sync + exists + empty current version → OK",
			tool: mockTool{
				exists:         true,
				start:          strategy.Sync,
				currentVersion: "",
				targetVersion:  "v1.0.0",
			},
			wantOK: true,
		},
		{
			name: "Sync + exists + empty target version → OK",
			tool: mockTool{
				exists:         true,
				start:          strategy.Sync,
				currentVersion: "v1.0.0",
				targetVersion:  "",
			},
			wantOK: true,
		},
		{
			name: "Sync + tool does not exist → OK",
			tool: mockTool{
				exists:        false,
				start:         strategy.Sync,
				targetVersion: "v1.0.0",
			},
			wantOK: true,
		},
		{
			// version.Parse returns nil for strings with no semver substring,
			// exercising the source == nil early-return path.
			name: "Sync + exists + unparseable current version → OK",
			tool: mockTool{
				exists:         true,
				start:          strategy.Sync,
				currentVersion: "not-a-version",
				targetVersion:  "v1.0.0",
			},
			wantOK: true,
		},

		// Force strategy
		{
			name: "Force + tool exists → OK",
			tool: mockTool{
				exists:         true,
				start:          strategy.Force,
				currentVersion: "v1.0.0",
				targetVersion:  "v1.0.0",
			},
			wantOK: true,
		},
		{
			name: "Force + tool does not exist → OK",
			tool: mockTool{
				exists: false,
				start:  strategy.Force,
			},
			wantOK: true,
		},

		// Existing strategy
		{
			// The early-return for !t.Exists() fires before the strategy switch,
			// so Existing with a missing tool returns OK, not Skipped.
			name: "Existing + tool does not exist → OK",
			tool: mockTool{
				exists: false,
				start:  strategy.Existing,
			},
			wantOK: true,
		},
		{
			name: "Existing + exists + versions match → Skipped",
			tool: mockTool{
				exists:         true,
				start:          strategy.Existing,
				currentVersion: "v1.0.0",
				targetVersion:  "v1.0.0",
			},
			wantSkipped: true,
		},
		{
			name: "Existing + exists + versions differ → OK",
			tool: mockTool{
				exists:         true,
				start:          strategy.Existing,
				currentVersion: "v1.0.0",
				targetVersion:  "v1.1.0",
			},
			wantOK: true,
		},
		{
			// currentVersion == "" triggers the early-return path before version comparison.
			name: "Existing + exists + empty current version → OK",
			tool: mockTool{
				exists:         true,
				start:          strategy.Existing,
				currentVersion: "",
				targetVersion:  "v1.0.0",
			},
			wantOK: true,
		},

		// Unknown strategy
		{
			name: "unknown strategy → Failed",
			tool: mockTool{
				exists: true,
				start:  strategy.Strategy("bogus"),
			},
			wantFailed: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := s.Sync(tc.tool)

			switch {
			case tc.wantOK && !got.IsOK():
				t.Errorf("Sync(%+v): got status %v, want OK", tc.tool, got)
			case tc.wantSkipped && !got.IsSkipped():
				t.Errorf("Sync(%+v): got status %v, want Skipped", tc.tool, got)
			case tc.wantFailed && !got.IsFailed():
				t.Errorf("Sync(%+v): got status %v, want Failed", tc.tool, got)
			default:
				if !tc.wantOK && !tc.wantSkipped && !tc.wantFailed {
					t.Errorf("Sync(%+v): no expected status set in test case", tc.tool)
				}
			}
		})
	}
}
