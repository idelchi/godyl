package env_test

import (
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/idelchi/godyl/pkg/env"
)

func TestAdd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		pairs     []string
		wantKey   string
		wantValue string
		wantErr   bool
	}{
		{
			name:      "simple key=value",
			pairs:     []string{"KEY=VALUE"},
			wantKey:   "KEY",
			wantValue: "VALUE",
			wantErr:   false,
		},
		{
			name:      "value contains equals sign",
			pairs:     []string{"KEY=VAL=UE"},
			wantKey:   "KEY",
			wantValue: "VAL=UE",
			wantErr:   false,
		},
		{
			name:      "empty value",
			pairs:     []string{"KEY="},
			wantKey:   "KEY",
			wantValue: "",
			wantErr:   false,
		},
		{
			name:    "missing equals sign",
			pairs:   []string{"INVALID"},
			wantErr: true,
		},
		{
			name:    "empty key (=VALUE)",
			pairs:   []string{"=VALUE"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			e := make(env.Env)
			err := e.Add(tc.pairs...)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("Add(%v): expected error, got nil", tc.pairs)
				}

				if !errors.Is(err, env.ErrEnvMalformed) {
					t.Errorf("Add(%v): error = %v, want wrapping ErrEnvMalformed", tc.pairs, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("Add(%v): unexpected error: %v", tc.pairs, err)
			}

			got, ok := e[tc.wantKey]
			if !ok {
				t.Fatalf("Add(%v): key %q not found in env", tc.pairs, tc.wantKey)
			}

			if got != tc.wantValue {
				t.Errorf("Add(%v): env[%q] = %q, want %q", tc.pairs, tc.wantKey, got, tc.wantValue)
			}
		})
	}

	t.Run("multiple pairs added together", func(t *testing.T) {
		t.Parallel()

		e := make(env.Env)

		if err := e.Add("A=1", "B=2"); err != nil {
			t.Fatalf("Add: unexpected error: %v", err)
		}

		for key, want := range map[string]string{"A": "1", "B": "2"} {
			if got := e.Get(key); got != want {
				t.Errorf("env[%q] = %q, want %q", key, got, want)
			}
		}
	})

	t.Run("partial success: valid pairs set, error still returned", func(t *testing.T) {
		t.Parallel()

		e := make(env.Env)

		err := e.Add("A=1", "INVALID", "B=2")
		if err == nil {
			t.Fatal("Add(partial): expected error, got nil")
		}

		if !errors.Is(err, env.ErrEnvMalformed) {
			t.Errorf("Add(partial): error = %v, want wrapping ErrEnvMalformed", err)
		}

		if e.Get("A") != "1" {
			t.Errorf("Add(partial): env[\"A\"] = %q, want \"1\"", e.Get("A"))
		}

		if e.Get("B") != "2" {
			t.Errorf("Add(partial): env[\"B\"] = %q, want \"2\"", e.Get("B"))
		}
	})
}

func TestGet(t *testing.T) {
	t.Parallel()

	e := env.Env{"PRESENT": "hello"}

	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "existing key", key: "PRESENT", want: "hello"},
		{name: "missing key", key: "ABSENT", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := e.Get(tc.key)
			if got != tc.want {
				t.Errorf("Get(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestGetAny(t *testing.T) {
	t.Parallel()

	e := env.Env{"B": "second"}

	tests := []struct {
		name string
		keys []string
		want string
	}{
		{
			name: "returns first found",
			keys: []string{"A", "B", "C"},
			want: "second",
		},
		{
			name: "none found returns empty",
			keys: []string{"X", "Y"},
			want: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := e.GetAny(tc.keys...)
			if got != tc.want {
				t.Errorf("GetAny(%v) = %q, want %q", tc.keys, got, tc.want)
			}
		})
	}
}

func TestGetOrDefault(t *testing.T) {
	t.Parallel()

	e := env.Env{"KEY": "value"}

	tests := []struct {
		name         string
		key          string
		defaultValue string
		want         string
	}{
		{name: "existing key returns value", key: "KEY", defaultValue: "fallback", want: "value"},
		{name: "missing key returns default", key: "MISSING", defaultValue: "fallback", want: "fallback"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := e.GetOrDefault(tc.key, tc.defaultValue)
			if got != tc.want {
				t.Errorf("GetOrDefault(%q, %q) = %q, want %q", tc.key, tc.defaultValue, got, tc.want)
			}
		})
	}
}

func TestExists(t *testing.T) {
	t.Parallel()

	e := env.Env{"PRESENT": "yes"}

	tests := []struct {
		name string
		key  string
		want bool
	}{
		{name: "present key", key: "PRESENT", want: true},
		{name: "absent key", key: "ABSENT", want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := e.Exists(tc.key)
			if got != tc.want {
				t.Errorf("Exists(%q) = %v, want %v", tc.key, got, tc.want)
			}
		})
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("deleted key is gone and other keys intact", func(t *testing.T) {
		t.Parallel()

		e := env.Env{"A": "1", "B": "2", "C": "3"}
		e.Delete("B")

		if e.Exists("B") {
			t.Error("Delete(\"B\"): key still present after deletion")
		}

		if !e.Exists("A") {
			t.Error("Delete(\"B\"): key \"A\" unexpectedly removed")
		}

		if !e.Exists("C") {
			t.Error("Delete(\"B\"): key \"C\" unexpectedly removed")
		}
	})

	t.Run("delete nonexistent key", func(t *testing.T) {
		t.Parallel()

		e := env.Env{"A": "1"}
		e.Delete("MISSING") // should not panic

		if !e.Exists("A") {
			t.Error("existing key should survive")
		}
	})
}

func TestMerge(t *testing.T) {
	t.Parallel()

	t.Run("receiver wins on conflict", func(t *testing.T) {
		t.Parallel()

		// Merge copies the incoming envs first, then overwrites with receiver,
		// so receiver keys take precedence over incoming keys.
		e := env.Env{"A": "1", "B": "2"}
		e.Merge(env.Env{"B": "3", "C": "4"})

		cases := map[string]string{
			"A": "1", // receiver-only key preserved
			"B": "2", // receiver wins over incoming
			"C": "4", // incoming-only key added
		}

		for key, want := range cases {
			got := e.Get(key)
			if got != want {
				t.Errorf("after Merge, env[%q] = %q, want %q", key, got, want)
			}
		}
	})

	t.Run("multiple incoming envs", func(t *testing.T) {
		t.Parallel()

		e := env.Env{"A": "receiver"}
		e.Merge(
			env.Env{"B": "from-first", "A": "from-first"},
			env.Env{"C": "from-second", "B": "from-second"},
		)

		// receiver wins for A
		if got := e.Get("A"); got != "receiver" {
			t.Errorf("env[\"A\"] = %q, want \"receiver\"", got)
		}

		// B: first incoming sets it, second incoming overwrites (last write from Merge wins
		// among incoming), then receiver overwrites — but receiver doesn't have B, so last
		// incoming wins. The Merge implementation processes envs in order and then *e last,
		// so B ends up as "from-second" (second incoming overwrites first).
		if got := e.Get("B"); got != "from-second" {
			t.Errorf("env[\"B\"] = %q, want \"from-second\"", got)
		}

		// C comes only from second incoming
		if got := e.Get("C"); got != "from-second" {
			t.Errorf("env[\"C\"] = %q, want \"from-second\"", got)
		}
	})
}

func TestMergedWith(t *testing.T) {
	t.Parallel()

	t.Run("returns new env without mutating original", func(t *testing.T) {
		t.Parallel()

		original := env.Env{"A": "1"}
		incoming := env.Env{"B": "2"}

		merged := original.MergedWith(incoming)

		// merged must contain both keys
		if merged.Get("A") != "1" {
			t.Errorf("merged[\"A\"] = %q, want \"1\"", merged.Get("A"))
		}

		if merged.Get("B") != "2" {
			t.Errorf("merged[\"B\"] = %q, want \"2\"", merged.Get("B"))
		}

		// original must not have gained the incoming key
		if original.Exists("B") {
			t.Error("MergedWith mutated the original env")
		}
	})

	t.Run("receiver wins on shared key", func(t *testing.T) {
		t.Parallel()

		original := env.Env{"SHARED": "original"}
		incoming := env.Env{"SHARED": "incoming"}

		merged := original.MergedWith(incoming)

		// MergedWith calls Merge on a clone of original, so original wins.
		if got := merged.Get("SHARED"); got != "original" {
			t.Errorf("MergedWith conflict: merged[\"SHARED\"] = %q, want \"original\"", got)
		}
	})
}

func TestAsSlice(t *testing.T) {
	t.Parallel()

	t.Run("sorted key=value pairs", func(t *testing.T) {
		t.Parallel()

		e := env.Env{"B": "2", "A": "1"}
		got := e.AsSlice()

		want := []string{"A=1", "B=2"}

		if !slices.Equal(got, want) {
			t.Errorf("AsSlice() = %v, want %v", got, want)
		}
	})
}

func TestKeys(t *testing.T) {
	t.Parallel()

	t.Run("sorted keys", func(t *testing.T) {
		t.Parallel()

		e := env.Env{"C": "3", "A": "1", "B": "2"}
		got := e.Keys()

		want := []string{"A", "B", "C"}

		if !slices.Equal(got, want) {
			t.Errorf("Keys() = %v, want %v", got, want)
		}
	})
}

func TestMustGet(t *testing.T) {
	t.Parallel()

	e := env.Env{"FOUND": "yes"}

	tests := []struct {
		name    string
		key     string
		want    string
		wantErr bool
	}{
		{name: "existing key", key: "FOUND", want: "yes", wantErr: false},
		{name: "missing key", key: "MISSING", want: "", wantErr: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := e.MustGet(tc.key)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("MustGet(%q): expected error, got nil", tc.key)
				}

				if !errors.Is(err, env.ErrEnvVarNotFound) {
					t.Errorf("MustGet(%q): error = %v, want wrapping ErrEnvVarNotFound", tc.key, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("MustGet(%q): unexpected error: %v", tc.key, err)
			}

			if got != tc.want {
				t.Errorf("MustGet(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestGetWithPredicates(t *testing.T) {
	t.Parallel()

	t.Run("filter by key prefix", func(t *testing.T) {
		t.Parallel()

		e := env.Env{
			"APP_HOST": "localhost",
			"APP_PORT": "8080",
			"DB_HOST":  "dbserver",
		}

		prefixPredicate := func(key, _ string) bool {
			return strings.HasPrefix(key, "APP_")
		}

		got := e.GetWithPredicates(prefixPredicate)

		if len(got) != 2 {
			t.Fatalf("GetWithPredicates: got %d entries, want 2", len(got))
		}

		if got.Get("APP_HOST") != "localhost" {
			t.Errorf("GetWithPredicates: APP_HOST = %q, want \"localhost\"", got.Get("APP_HOST"))
		}

		if got.Get("APP_PORT") != "8080" {
			t.Errorf("GetWithPredicates: APP_PORT = %q, want \"8080\"", got.Get("APP_PORT"))
		}

		if got.Exists("DB_HOST") {
			t.Error("GetWithPredicates: DB_HOST should not be in result")
		}
	})
}

func TestAddPair(t *testing.T) {
	t.Parallel()

	t.Run("valid pair sets value", func(t *testing.T) {
		t.Parallel()

		e := make(env.Env)

		if err := e.AddPair("KEY", "value"); err != nil {
			t.Fatalf("AddPair(\"KEY\", \"value\"): unexpected error: %v", err)
		}

		if got := e.Get("KEY"); got != "value" {
			t.Errorf("AddPair: env[\"KEY\"] = %q, want \"value\"", got)
		}
	})

	t.Run("empty key returns ErrEnvMalformed", func(t *testing.T) {
		t.Parallel()

		e := make(env.Env)

		err := e.AddPair("", "value")
		if err == nil {
			t.Fatal("AddPair(\"\", \"value\"): expected error, got nil")
		}

		if !errors.Is(err, env.ErrEnvMalformed) {
			t.Errorf("AddPair(\"\", \"value\"): error = %v, want wrapping ErrEnvMalformed", err)
		}
	})
}

func TestGetAsEnv(t *testing.T) {
	t.Parallel()

	e := env.Env{"KEY": "value"}

	tests := []struct {
		name string
		key  string
		want string
	}{
		{name: "existing key returns KEY=value format", key: "KEY", want: "KEY=value"},
		{name: "missing key returns empty string", key: "MISSING", want: ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := e.GetAsEnv(tc.key)
			if got != tc.want {
				t.Errorf("GetAsEnv(%q) = %q, want %q", tc.key, got, tc.want)
			}
		})
	}
}

func TestAsEnv(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   []string
		want    env.Env
		wantErr bool
	}{
		{
			name:    "valid pairs",
			input:   []string{"A=1", "B=2"},
			want:    env.Env{"A": "1", "B": "2"},
			wantErr: false,
		},
		{
			name:    "malformed entry",
			input:   []string{"INVALID"},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := env.AsEnv(tc.input...)

			if tc.wantErr {
				if err == nil {
					t.Fatalf("AsEnv(%v): expected error, got nil", tc.input)
				}

				if !errors.Is(err, env.ErrEnvMalformed) {
					t.Errorf("AsEnv(%v): error = %v, want wrapping ErrEnvMalformed", tc.input, err)
				}

				return
			}

			if err != nil {
				t.Fatalf("AsEnv(%v): unexpected error: %v", tc.input, err)
			}

			for key, want := range tc.want {
				if got.Get(key) != want {
					t.Errorf("AsEnv(%v): env[%q] = %q, want %q", tc.input, key, got.Get(key), want)
				}
			}

			if len(got) != len(tc.want) {
				t.Errorf("AsEnv(%v): got %d entries, want %d", tc.input, len(got), len(tc.want))
			}
		})
	}
}

func TestAddZeroPairs(t *testing.T) {
	t.Parallel()

	// Add() with no arguments should be a no-op and return nil.
	e := make(env.Env)

	if err := e.Add(); err != nil {
		t.Fatalf("Add(): unexpected error: %v", err)
	}

	if len(e) != 0 {
		t.Errorf("Add(): env has %d entries after zero-pair call, want 0", len(e))
	}
}

func TestAsSliceEmpty(t *testing.T) {
	t.Parallel()

	// AsSlice on an empty Env must return an empty (not nil) slice.
	e := make(env.Env)
	got := e.AsSlice()

	if len(got) != 0 {
		t.Errorf("AsSlice() on empty Env = %v, want empty slice", got)
	}
}

func TestKeysEmpty(t *testing.T) {
	t.Parallel()

	// Keys on an empty Env must return an empty (not nil) slice.
	e := make(env.Env)
	got := e.Keys()

	if len(got) != 0 {
		t.Errorf("Keys() on empty Env = %v, want empty slice", got)
	}
}
