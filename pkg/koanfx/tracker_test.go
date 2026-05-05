package koanfx_test

import (
	"slices"
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"

	"github.com/idelchi/godyl/pkg/koanfx"
)

func TestTrackerBasics(t *testing.T) {
	t.Parallel()

	tr := koanfx.NewTracker()

	if tr.Exists("key") {
		t.Error("Exists(\"key\") = true on empty tracker, want false")
	}

	if tr.IsSet("key") {
		t.Error("IsSet(\"key\") = true on empty tracker, want false")
	}

	if got := len(tr.Names()); got != 0 {
		t.Errorf("len(Names()) = %d on empty tracker, want 0", got)
	}
}

func TestTrackerTrackAll(t *testing.T) {
	t.Parallel()

	k := koanf.New(".")
	if err := k.Set("foo", "bar"); err != nil {
		t.Fatalf("k.Set(\"foo\", \"bar\") error: %v", err)
	}

	if err := k.Set("baz", "qux"); err != nil {
		t.Fatalf("k.Set(\"baz\", \"qux\") error: %v", err)
	}

	tr := koanfx.NewTracker()
	tr.TrackAll(k)

	tests := []struct {
		key       string
		wantExist bool
		wantSet   bool
	}{
		{key: "foo", wantExist: true, wantSet: true},
		{key: "baz", wantExist: true, wantSet: true},
		{key: "missing", wantExist: false, wantSet: false},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			t.Parallel()

			if got := tr.Exists(tc.key); got != tc.wantExist {
				t.Errorf("Exists(%q) = %v, want %v", tc.key, got, tc.wantExist)
			}

			if got := tr.IsSet(tc.key); got != tc.wantSet {
				t.Errorf("IsSet(%q) = %v, want %v", tc.key, got, tc.wantSet)
			}
		})
	}
}

func TestTrackerTrackFlags(t *testing.T) {
	t.Parallel()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("output", "./bin", "output path")
	flags.String("source", "github", "source type")

	if err := flags.Parse([]string{"--output", "/usr/bin"}); err != nil {
		t.Fatalf("flags.Parse error: %v", err)
	}

	tr := koanfx.NewTracker()
	tr.TrackFlags(flags)

	tests := []struct {
		key       string
		wantExist bool
		wantSet   bool
	}{
		// output was explicitly changed via Parse, so it must be both present and set.
		{key: "output", wantExist: true, wantSet: true},
		// source was registered but not changed, so it exists but is not set.
		{key: "source", wantExist: true, wantSet: false},
	}

	for _, tc := range tests {
		t.Run(tc.key, func(t *testing.T) {
			t.Parallel()

			if got := tr.Exists(tc.key); got != tc.wantExist {
				t.Errorf("Exists(%q) = %v, want %v", tc.key, got, tc.wantExist)
			}

			if got := tr.IsSet(tc.key); got != tc.wantSet {
				t.Errorf("IsSet(%q) = %v, want %v", tc.key, got, tc.wantSet)
			}
		})
	}
}

// TestTrackerTrackAllThenFlags verifies that a key already marked as set via
// TrackAll is not downgraded to unset when TrackFlags visits the same key
// with Changed=false.
func TestTrackerTrackAllThenFlags(t *testing.T) {
	t.Parallel()

	k := koanf.New(".")
	if err := k.Set("output", "/usr/bin"); err != nil {
		t.Fatalf("k.Set(\"output\") error: %v", err)
	}

	tr := koanfx.NewTracker()
	tr.TrackAll(k)

	// Sanity: output is set after TrackAll.
	if !tr.IsSet("output") {
		t.Fatal("IsSet(\"output\") = false after TrackAll, want true")
	}

	// Now register a flag for "output" but do NOT parse it (Changed=false).
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("output", "./bin", "output path")

	tr.TrackFlags(flags)

	// The set status from TrackAll must be preserved; TrackFlags must not
	// downgrade it even though the flag was not explicitly changed.
	if !tr.Exists("output") {
		t.Error("Exists(\"output\") = false after TrackFlags, want true")
	}

	if !tr.IsSet("output") {
		t.Error("IsSet(\"output\") = false after TrackFlags, want true (must not downgrade set keys)")
	}
}

// TestTrackerFlagsThenAll verifies the reverse order: TrackFlags first (unchanged flag),
// then TrackAll with the same key present → IsSet must be true afterward.
func TestTrackerFlagsThenAll(t *testing.T) {
	t.Parallel()

	// Register a flag for "output" but do NOT parse it (Changed=false).
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("output", "./bin", "output path")

	tr := koanfx.NewTracker()
	tr.TrackFlags(flags)

	// After TrackFlags with unchanged flag, output exists but is not set.
	if !tr.Exists("output") {
		t.Fatal("Exists(\"output\") = false after TrackFlags, want true")
	}

	if tr.IsSet("output") {
		t.Error("IsSet(\"output\") = true after unchanged TrackFlags, want false")
	}

	// Now call TrackAll with a koanf that has the key.
	k := koanf.New(".")
	if err := k.Set("output", "/usr/bin"); err != nil {
		t.Fatalf("k.Set(\"output\") error: %v", err)
	}

	tr.TrackAll(k)

	// TrackAll must upgrade the key to set=true.
	if !tr.Exists("output") {
		t.Error("Exists(\"output\") = false after TrackAll, want true")
	}

	if !tr.IsSet("output") {
		t.Error("IsSet(\"output\") = false after TrackAll, want true")
	}
}

func TestTrackerNames(t *testing.T) {
	t.Parallel()

	k := koanf.New(".")
	if err := k.Set("a", 1); err != nil {
		t.Fatalf("k.Set(\"a\", 1) error: %v", err)
	}

	if err := k.Set("b", 2); err != nil {
		t.Fatalf("k.Set(\"b\", 2) error: %v", err)
	}

	tr := koanfx.NewTracker()
	tr.TrackAll(k)

	got := tr.Names()

	if len(got) != 2 {
		t.Fatalf("len(Names()) = %d, want 2; got %v", len(got), got)
	}

	slices.Sort(got)

	want := []string{"a", "b"}

	if !slices.Equal(got, want) {
		t.Errorf("Names() = %v, want %v", got, want)
	}
}

// TestTrackerTrackFlagsDuplicate verifies that calling TrackFlags twice with
// the same FlagSet (once with the flag unchanged, once after marking it
// Changed) results in the final IsSet value reflecting the Changed state
// from the second call.
func TestTrackerTrackFlagsDuplicate(t *testing.T) {
	t.Parallel()

	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("output", "./bin", "output path")

	tr := koanfx.NewTracker()

	// First call: flag not yet changed → exists but not set.
	tr.TrackFlags(flags)

	if !tr.Exists("output") {
		t.Fatal("Exists(\"output\") = false after first TrackFlags, want true")
	}

	if tr.IsSet("output") {
		t.Error("IsSet(\"output\") = true after first TrackFlags (unchanged), want false")
	}

	// Simulate the flag being changed (e.g., by Parse) and call TrackFlags again.
	if err := flags.Parse([]string{"--output", "/usr/local/bin"}); err != nil {
		t.Fatalf("flags.Parse error: %v", err)
	}

	// Second call: flag is now Changed → must be marked as set.
	tr.TrackFlags(flags)

	if !tr.Exists("output") {
		t.Error("Exists(\"output\") = false after second TrackFlags, want true")
	}

	if !tr.IsSet("output") {
		t.Error("IsSet(\"output\") = false after second TrackFlags (changed), want true")
	}
}

// ---------------------------------------------------------------------------
// Koanf wrapper tests
// ---------------------------------------------------------------------------

func TestKoanfNew(t *testing.T) {
	t.Parallel()

	k := koanfx.New()

	if k == nil {
		t.Fatal("New() returned nil")
	}

	if k.Tracker == nil {
		t.Error("New().Tracker is nil, want initialized")
	}

	// Fresh instance has no keys.
	if got := k.Map(); len(got) != 0 {
		t.Errorf("New().Map() has %d entries, want 0", len(got))
	}
}

func TestKoanfFromStruct(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name  string `yaml:"name"`
		Count int    `yaml:"count"`
	}

	k, err := koanfx.FromStruct(cfg{Name: "test", Count: 5}, "yaml")
	if err != nil {
		t.Fatalf("FromStruct() unexpected error: %v", err)
	}

	if got := k.String("name"); got != "test" {
		t.Errorf("String(\"name\") = %q, want %q", got, "test")
	}

	if got := k.Int("count"); got != 5 {
		t.Errorf("Int(\"count\") = %d, want %d", got, 5)
	}
}

func TestKoanfWithFlags(t *testing.T) {
	t.Parallel()

	k := koanfx.New()
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("output", "./bin", "output path")

	k2 := k.WithFlags(flags)

	// Must be a new instance sharing the same Tracker.
	if k2.Tracker != k.Tracker {
		t.Error("WithFlags() should share the same Tracker")
	}
}

func TestKoanfClearTracker(t *testing.T) {
	t.Parallel()

	k := koanfx.New()
	if err := k.Set("key", "val"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	k.TrackAll().Track()

	if !k.Tracker.IsSet("key") {
		t.Fatal("IsSet(\"key\") = false after TrackAll, want true")
	}

	k2 := k.ClearTracker()

	if k2.Tracker.IsSet("key") {
		t.Error("ClearTracker().Tracker.IsSet(\"key\") = true, want false (fresh tracker)")
	}
}

func TestKoanfWithKoanf(t *testing.T) {
	t.Parallel()

	k := koanfx.New()

	newK := koanf.New(".")
	if err := newK.Set("x", 42); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	k2 := k.WithKoanf(newK)

	if got := k2.Int("x"); got != 42 {
		t.Errorf("WithKoanf().Int(\"x\") = %d, want %d", got, 42)
	}
}

func TestKoanfWithKoanfNil(t *testing.T) {
	t.Parallel()

	k := koanfx.New()
	k2 := k.WithKoanf(nil)

	// Nil input should produce a fresh empty koanf, not panic.
	if got := k2.Map(); len(got) != 0 {
		t.Errorf("WithKoanf(nil).Map() has %d entries, want 0", len(got))
	}
}

func TestKoanfIsSet(t *testing.T) {
	t.Parallel()

	k := koanfx.New()
	if err := k.Set("present", "yes"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	k.TrackAll().Track()

	if !k.IsSet("present") {
		t.Error("IsSet(\"present\") = false, want true")
	}

	if k.IsSet("absent") {
		t.Error("IsSet(\"absent\") = true, want false")
	}
}

func TestKoanfFiltered(t *testing.T) {
	t.Parallel()

	k := koanfx.New()

	for key, val := range map[string]any{"a": 1, "b": 2, "c": 3} {
		if err := k.Set(key, val); err != nil {
			t.Fatalf("Set(%q) error: %v", key, err)
		}
	}

	filtered := k.Filtered("a", "c", "nonexistent")
	m := filtered.Map()

	if len(m) != 2 {
		t.Fatalf("Filtered().Map() has %d entries, want 2; got %v", len(m), m)
	}

	if _, ok := m["a"]; !ok {
		t.Error("Filtered() missing key \"a\"")
	}

	if _, ok := m["c"]; !ok {
		t.Error("Filtered() missing key \"c\"")
	}

	if _, ok := m["b"]; ok {
		t.Error("Filtered() should not contain key \"b\"")
	}
}

func TestKoanfMap(t *testing.T) {
	t.Parallel()

	k := koanfx.New()
	if err := k.Set("x", "y"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	m := k.Map()
	if m["x"] != "y" {
		t.Errorf("Map()[\"x\"] = %v, want \"y\"", m["x"])
	}
}

func TestKoanfTrackAllAndTrackFlags(t *testing.T) {
	t.Parallel()

	k := koanfx.New()
	if err := k.Set("key", "val"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	// TrackAll sets active, Track executes it.
	k.TrackAll().Track()

	if !k.IsSet("key") {
		t.Error("IsSet(\"key\") = false after TrackAll+Track, want true")
	}

	// TrackFlags with no flags set should not panic.
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("output", "./bin", "output")

	k2 := k.WithFlags(flags)
	k2.TrackFlags().Track()

	// "output" was registered but not changed, so not set.
	if k2.IsSet("output") {
		t.Error("IsSet(\"output\") = true after unchanged TrackFlags, want false")
	}
}

func TestKoanfTrackNilActive(t *testing.T) {
	t.Parallel()

	// Track() with no active function set should be a no-op (no panic).
	k := koanfx.New()
	k.Track() // must not panic
}

func TestKoanfUnmarshal(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `mapstructure:"name"`
	}

	k := koanfx.New()
	if err := k.Set("name", "hello"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	var got cfg

	if err := k.Unmarshal(&got); err != nil {
		t.Fatalf("Unmarshal() error: %v", err)
	}

	if got.Name != "hello" {
		t.Errorf("Unmarshal() Name = %q, want %q", got.Name, "hello")
	}
}

func TestKoanfUnmarshalWithMetadata(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `mapstructure:"name"`
	}

	k := koanfx.New()
	if err := k.Set("name", "hello"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	var got cfg

	meta, err := k.UnmarshalWithMetadata(&got)
	if err != nil {
		t.Fatalf("UnmarshalWithMetadata() error: %v", err)
	}

	if got.Name != "hello" {
		t.Errorf("UnmarshalWithMetadata() Name = %q, want %q", got.Name, "hello")
	}

	if len(meta.Keys) == 0 {
		t.Error("UnmarshalWithMetadata() metadata.Keys is empty, want non-empty")
	}
}

func TestKoanfOptions(t *testing.T) {
	t.Parallel()

	t.Run("WithFlatPaths sets FlatPaths", func(t *testing.T) {
		t.Parallel()

		conf := koanfx.NewDefaultUnmarshalConfig()
		koanfx.WithFlatPaths()(&conf)

		if !conf.FlatPaths {
			t.Error("WithFlatPaths() did not set FlatPaths to true")
		}
	})

	t.Run("WithErrorUnused sets ErrorUnused", func(t *testing.T) {
		t.Parallel()

		conf := koanfx.NewDefaultUnmarshalConfig()
		koanfx.WithErrorUnused()(&conf)

		if conf.DecoderConfig == nil {
			t.Fatal("WithErrorUnused() did not initialize DecoderConfig")
		}

		if !conf.DecoderConfig.ErrorUnused {
			t.Error("WithErrorUnused() did not set ErrorUnused to true")
		}
	})

	t.Run("WithSquash sets Squash", func(t *testing.T) {
		t.Parallel()

		conf := koanfx.NewDefaultUnmarshalConfig()
		koanfx.WithSquash()(&conf)

		if conf.DecoderConfig == nil {
			t.Fatal("WithSquash() did not initialize DecoderConfig")
		}

		if !conf.DecoderConfig.Squash {
			t.Error("WithSquash() did not set Squash to true")
		}
	})
}

func TestKoanfNewDefaultUnmarshalConfig(t *testing.T) {
	t.Parallel()

	conf := koanfx.NewDefaultUnmarshalConfig()

	if conf.Tag != "mapstructure" {
		t.Errorf("NewDefaultUnmarshalConfig().Tag = %q, want %q", conf.Tag, "mapstructure")
	}
}

func TestKoanfNewUnmarshalConfig(t *testing.T) {
	t.Parallel()

	conf := koanfx.NewUnmarshalConfig()

	if conf.DecoderConfig == nil {
		t.Fatal("NewUnmarshalConfig().DecoderConfig is nil")
	}

	if !conf.DecoderConfig.WeaklyTypedInput {
		t.Error("NewUnmarshalConfig().DecoderConfig.WeaklyTypedInput = false, want true")
	}
}

// ---------------------------------------------------------------------------
// Standalone unmarshal function tests
// ---------------------------------------------------------------------------

func TestUnmarshalStandalone(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name  string `mapstructure:"name"`
		Count int    `mapstructure:"count"`
	}

	k := koanf.New(".")
	if err := k.Set("name", "standalone"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	if err := k.Set("count", 7); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	var got cfg

	if err := koanfx.Unmarshal(k, "", &got); err != nil {
		t.Fatalf("Unmarshal() error: %v", err)
	}

	if got.Name != "standalone" {
		t.Errorf("Name = %q, want %q", got.Name, "standalone")
	}

	if got.Count != 7 {
		t.Errorf("Count = %d, want %d", got.Count, 7)
	}
}

func TestUnmarshalAllStandalone(t *testing.T) {
	t.Parallel()

	type cfg struct {
		X string `mapstructure:"x"`
	}

	k := koanf.New(".")
	if err := k.Set("x", "value"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	var got cfg

	if err := koanfx.UnmarshalAll(k, &got); err != nil {
		t.Fatalf("UnmarshalAll() error: %v", err)
	}

	if got.X != "value" {
		t.Errorf("X = %q, want %q", got.X, "value")
	}
}

func TestUnmarshalWithMetadataStandalone(t *testing.T) {
	t.Parallel()

	type cfg struct {
		A string `mapstructure:"a"`
	}

	k := koanf.New(".")
	if err := k.Set("a", "val"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	var got cfg

	meta, err := koanfx.UnmarshalWithMetadata(k, "", &got)
	if err != nil {
		t.Fatalf("UnmarshalWithMetadata() error: %v", err)
	}

	if got.A != "val" {
		t.Errorf("A = %q, want %q", got.A, "val")
	}

	if len(meta.Keys) == 0 {
		t.Error("metadata.Keys is empty, want non-empty")
	}
}

func TestUnmarshalAllWithMetadataStandalone(t *testing.T) {
	t.Parallel()

	type cfg struct {
		B int `mapstructure:"b"`
	}

	k := koanf.New(".")
	if err := k.Set("b", 99); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	var got cfg

	meta, err := koanfx.UnmarshalAllWithMetadata(k, &got)
	if err != nil {
		t.Fatalf("UnmarshalAllWithMetadata() error: %v", err)
	}

	if got.B != 99 {
		t.Errorf("B = %d, want %d", got.B, 99)
	}

	if len(meta.Keys) == 0 {
		t.Error("metadata.Keys is empty, want non-empty")
	}
}

func TestUnmarshalWithOptions(t *testing.T) {
	t.Parallel()

	type cfg struct {
		Name string `mapstructure:"name"`
	}

	k := koanf.New(".")
	if err := k.Set("name", "test"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	if err := k.Set("extra", "unused"); err != nil {
		t.Fatalf("Set() error: %v", err)
	}

	var got cfg

	err := koanfx.Unmarshal(k, "", &got, koanfx.WithErrorUnused())
	if err == nil {
		t.Fatal("Unmarshal with WithErrorUnused() expected error for unused key, got nil")
	}
}
