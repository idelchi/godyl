package validator_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/idelchi/godyl/pkg/validator"
)

type testConfig struct {
	Name  string `validate:"required"`
	Count int    `validate:"min=1,max=100"`
}

func TestValidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   testConfig
		wantErr bool
	}{
		{
			name:    "valid config",
			input:   testConfig{Name: "hello", Count: 50},
			wantErr: false,
		},
		{
			name:    "missing required Name",
			input:   testConfig{Name: "", Count: 50},
			wantErr: true,
		},
		{
			name:    "Count below min",
			input:   testConfig{Name: "hello", Count: 0},
			wantErr: true,
		},
		{
			name:    "Count above max",
			input:   testConfig{Name: "hello", Count: 101},
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			v := validator.New()
			errs := v.Validate(tc.input)

			if tc.wantErr && len(errs) == 0 {
				t.Fatalf("expected validation errors but got none")
			}

			if !tc.wantErr && len(errs) != 0 {
				t.Fatalf("expected no validation errors but got: %v", errs)
			}
		})
	}

	t.Run("empty Name and zero Count produce exactly two errors", func(t *testing.T) {
		t.Parallel()

		v := validator.New()
		errs := v.Validate(testConfig{Name: "", Count: 0})

		if len(errs) != 2 {
			t.Errorf("expected exactly 2 errors, got %d: %v", len(errs), errs)
		}
	})
}

func TestValidateHelper(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		input       testConfig
		wantErr     bool
		wantWrapped bool
	}{
		{
			name:        "valid config returns nil",
			input:       testConfig{Name: "hello", Count: 50},
			wantErr:     false,
			wantWrapped: false,
		},
		{
			name:        "invalid config wraps ErrValidation",
			input:       testConfig{Name: "", Count: 0},
			wantErr:     true,
			wantWrapped: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := validator.Validate(tc.input)

			if tc.wantErr && err == nil {
				t.Fatal("expected error but got nil")
			}

			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil error but got: %v", err)
			}

			if tc.wantWrapped && !errors.Is(err, validator.ErrValidation) {
				t.Errorf("expected error to wrap ErrValidation, got: %v", err)
			}
		})
	}
}

func TestRegisterValidationAndTranslation(t *testing.T) {
	t.Parallel()

	type evenConfig struct {
		Value int `validate:"even"`
	}

	validateEven := func(fl validator.FieldLevel) bool {
		return fl.Field().Int()%2 == 0
	}

	tests := []struct {
		name            string
		input           evenConfig
		wantErr         bool
		wantMsgContains string
	}{
		{
			name:    "even number passes",
			input:   evenConfig{Value: 4},
			wantErr: false,
		},
		{
			name:            "odd number fails with must be even message",
			input:           evenConfig{Value: 3},
			wantErr:         true,
			wantMsgContains: "must be even",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Each subtest gets its own validator so registration is isolated.
			v := validator.New()

			if err := v.RegisterValidationAndTranslation("even", validateEven, "{0} must be even"); err != nil {
				t.Fatalf("RegisterValidationAndTranslation: %v", err)
			}

			errs := v.Validate(tc.input)

			if tc.wantErr && len(errs) == 0 {
				t.Fatal("expected validation errors but got none")
			}

			if !tc.wantErr && len(errs) != 0 {
				t.Errorf("expected no validation errors but got: %v", errs)
			}

			if tc.wantMsgContains != "" {
				found := false

				for _, e := range errs {
					if strings.Contains(e.Error(), tc.wantMsgContains) {
						found = true

						break
					}
				}

				if !found {
					t.Errorf("expected an error containing %q, got: %v", tc.wantMsgContains, errs)
				}
			}
		})
	}
}

func TestValidator_Accessor(t *testing.T) {
	t.Parallel()

	v := validator.New()

	if v.Validator() == nil {
		t.Error("Validator().Validator() = nil, want non-nil")
	}
}

func TestFormatErrors(t *testing.T) {
	t.Parallel()

	t.Run("plain error is returned as-is in a slice", func(t *testing.T) {
		t.Parallel()

		v := validator.New()
		plain := errors.New("something")
		errs := v.FormatErrors(plain)

		if len(errs) != 1 {
			t.Fatalf("FormatErrors(plain): got %d errors, want 1", len(errs))
		}

		if !errors.Is(errs[0], plain) {
			t.Errorf("FormatErrors(plain): got %v, want %v", errs[0], plain)
		}
	})
}

func TestFormatErrorsNil(t *testing.T) {
	t.Parallel()

	// FormatErrors(nil) must not panic.  nil is not a validator.ValidationErrors
	// value, so the function wraps it in a single-element slice.  The slice has
	// length 1 and contains a nil error element.
	v := validator.New()
	errs := v.FormatErrors(nil)

	// The contract: no panic, result is a length-1 slice containing nil.
	if len(errs) != 1 {
		t.Fatalf("FormatErrors(nil): got %d errors, want 1", len(errs))
	}

	if errs[0] != nil {
		t.Errorf("FormatErrors(nil)[0] = %v, want nil", errs[0])
	}
}

func TestValidateErrorFieldInspection(t *testing.T) {
	t.Parallel()

	// Validate a struct where Name is required but left empty.
	// The error message produced by go-playground/validator with English
	// translations must mention the field name "Name".
	v := validator.New()
	errs := v.Validate(testConfig{Name: "", Count: 50})

	if len(errs) == 0 {
		t.Fatal("Validate(empty Name): expected errors, got none")
	}

	fieldMentioned := false

	for _, err := range errs {
		if strings.Contains(err.Error(), "Name") {
			fieldMentioned = true

			break
		}
	}

	if !fieldMentioned {
		t.Errorf("Validate(empty Name): expected at least one error mentioning \"Name\", got: %v", errs)
	}
}
