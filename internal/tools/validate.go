package tools

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/idelchi/godyl/internal/executable"
	"github.com/idelchi/godyl/internal/folder"
	stringlike "github.com/idelchi/godyl/internal/generic"
	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
)

var (
	ErrAlreadyExists   = fmt.Errorf("tool already exists")
	ErrDoesNotHaveTags = fmt.Errorf("tool does not have required tags")
	ErrSkipped         = fmt.Errorf("tool skipped")
	ErrFailed          = fmt.Errorf("tool failed")
)

func (t *Tool) Resolve(withTags []string, withoutTags []string) error {
	// Normalize values given in .Values map such that the first letter of keys
	// are capitalized.
	t.NormalizeValues()

	// Create instance of folder.Folder for the output
	output := folder.Folder(t.Output)
	if err := output.Expand(); err != nil {
		return err
	}

	t.Output = output.Path()

	// Save path to be templated last
	path := t.Path

	fallbacks := append([]string{t.Source.Type}, t.Fallbacks...)

	var lastErr error

	for _, fallback := range fallbacks {
		if err := t.tryResolveFallback(fallback, path, withTags, withoutTags); err != nil {
			lastErr = err

			continue // Try the next fallback
		}

		// If successful, return nil
		return nil
	}

	// If none of the fallbacks worked, return the last error encountered
	return lastErr
}

func (t *Tool) CheckSkipConditions(withTags []string, withoutTags []string) error {
	if !t.Tags.Has(withTags) {
		return ErrDoesNotHaveTags
	}

	if !t.Tags.HasNot(withoutTags) {
		return ErrDoesNotHaveTags
	}

	if err := t.Strategy.Check(t); err != nil {
		return err
	}

	if t.Skip.Skip {
		return ErrSkipped
	}

	return nil
}

func (t *Tool) tryResolveFallback(fallback string, path string, withTags []string, withoutTags []string) error {
	t.Tags.Append(t.Name)

	if err := t.CheckSkipConditions(withTags, withoutTags); err != nil {
		return err
	}

	t.Source.Type = fallback

	populator, err := t.Source.Installer()
	if err != nil {
		return err
	}

	if err := populator.Initialize(t.Name); err != nil {
		return err
	}

	if err := populator.Exe(); err != nil {
		return err
	}

	stringlike.SetIfEmpty(&t.Exe.Name, populator.Get("exe"))

	if err := t.Template(); err != nil {
		return err
	}

	if err := t.CheckSkipConditions(withTags, withoutTags); err != nil {
		return err
	}

	if stringlike.IsEmpty(t.Version) {
		if err := populator.Version(t.Name); err != nil {
			return err
		}
	}

	stringlike.SetIfEmpty(&t.Version, populator.Get("version"))

	if err := t.Template(); err != nil {
		return err
	}

	if stringlike.IsEmpty(t.Path) {
		if err := populator.Path(t.Name, t.Extensions, t.Version, match.Requirements{
			Platform: t.Platform,
			Hints:    t.Hints,
		}); err != nil {
			return err
		}
	}

	stringlike.SetIfEmpty(&t.Path, populator.Get("path"))
	stringlike.SetIfEmpty(&t.Path, path)

	if err := t.Template(); err != nil {
		return err
	}

	stringlike.SetIfEmpty(&t.Exe.Name, t.Name)

	t.Exe.Name += t.Platform.Extension.String()

	stringlike.SetSliceIfNil(&t.Exe.Patterns, fmt.Sprintf("^%s$", t.Exe.Name))

	for i, alias := range t.Aliases {
		t.Aliases[i] = alias + t.Platform.Extension.String()
	}

	t.Name = t.Name + t.Platform.Extension.String()

	if err := t.Strategy.Upgrade(t); err != nil {
		return err
	}

	// Validation
	return t.Validate()
}

func (t *Tool) Validate() error {
	validate := validator.New()
	if err := validate.Struct(t); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}
	return nil
}

func (t *Tool) Exists() bool {
	return executable.New(t.Output, t.Exe.Name).Exists()
}

type Installer interface {
	Install(d sources.InstallData) (output string, err error)
}

func (t *Tool) Download() (string, error) {
	installer, err := t.Source.Installer()
	if err != nil {
		return "", err
	}

	data := sources.InstallData{
		Path:     t.Path,
		Name:     t.Name,
		Exe:      t.Exe.Name,
		Patterns: t.Exe.Patterns,
		Output:   t.Output,
		Aliases:  t.Aliases,
	}

	return installer.Install(data)
}
