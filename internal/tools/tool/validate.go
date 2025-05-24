// Package tool provides core functionality for managing tool configurations.
package tool

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/templates"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
	"github.com/idelchi/godyl/pkg/utils"
	"github.com/idelchi/godyl/pkg/validator"
)

// Resolve attempts to resolve the tool's source and strategy based on the provided tags.
// It handles environment variables, fallbacks, templating, and validation of the tool's
// configuration. Returns a Result indicating success or failure with detailed messages.
func (t *Tool) Resolve(tags tags.IncludeTags, options ...ResolveOption) result.Result {
	// Initialize default options
	opts := resolveOptions{}

	// Apply all options
	for _, option := range options {
		option(&opts)
	}

	if err := t.Platform.Parse(); err != nil {
		return result.WithFailed("parsing platform").Wrap(err)
	}

	// Load environment variables from the system.
	t.Env.Merge(env.FromEnv())

	// Expand environment variables.
	t.Env.Expand()

	// Create a basic template engine with missing key errors and a map of platform values
	tmpl := templates.New(templates.WithMissingKeyError(), templates.WithSlimSprig()).
		WithValues(t.Platform.ToMap(), t.ToTemplateMap())

	if err := t.TemplateFirst(tmpl); err != nil {
		return result.WithFailed("templating first").Wrap(err)
	}

	// Expand and set the output folder path.
	t.Output = folder.New(t.Output).Expanded().Path()

	// Must pass at least after first set of templates.
	if err := t.Validate(); err != nil {
		return result.WithFailed("validating config").Wrap(err)
	}

	// Append the tool's name as a tag.
	t.Tags.Append(t.Name)

	// Try resolving with each fallback in order.
	var res result.Result

	for _, fallback := range t.Fallbacks.Build(t.Source.Type) {
		// Set the source type to the current fallback.
		t.Source.Type = fallback

		// Get the installer for the current source type.
		populator, err := t.Source.Installer()
		t.populator = populator

		if err != nil {
			return result.WithFailed(fmt.Sprintf("getting populator: %s", err))
		}

		// Initialize the installer.
		if err := populator.Initialize(t.Name); err != nil {
			return result.WithFailed(fmt.Sprintf("initializing populator: %s", err))
		}

		// Set the executable name according the source type's rules.
		utils.SetIfZero(&t.Exe.Name, populator.Get("exe"))
		utils.SetIfZero(&t.Exe.Name, t.Name)

		// Update the template engine with .exe
		tmpl.AddValue("Exe", t.Exe.Name)

		// Re-check skip conditions after applying templates.
		if res = t.CheckSkipConditions(tags); !res.IsOK() {
			return res
		}

		if opts.skipVersion {
			return result.WithSkipped("skipped version resolution")
		}

		if res = t.resolve(populator, tmpl, opts); res.IsFailed() {
			continue // Move on to the next fallback.
		}

		return res
	}

	return res
}

func (t *Tool) resolve(populator sources.Populator, tmpl *templates.Processor, opts resolveOptions) result.Result {
	// Retrieve the tool's version from the installer if it is not already set.
	if utils.IsZero(t.Version.Version) {
		if err := populator.Version(t.Name); err != nil {
			return result.WithFailed(fmt.Sprintf("getting version: %s", err))
		}

		t.Version.Version = populator.Get("version")
	}

	// Update the version to the template engine.
	tmpl.AddValue("Version", t.Version.Version)

	if opts.upUntilVersion {
		return result.WithSkipped("skipped after version resolution")
	}

	if err := t.TemplateLast(tmpl); err != nil {
		return result.WithFailed(fmt.Sprintf("templating last: %s", err))
	}

	// Attempt to sync the tool using the current strategy.
	outcome := t.Strategy.Sync(t)
	if !outcome.IsOK() {
		return outcome
	}

	// Validate the tool's configuration.
	if err := t.Validate(); err != nil {
		return result.WithFailed(fmt.Sprintf("validating config: %s", err))
	}

	if opts.skipURL {
		return outcome.Wrapped("requires download")
	}

	// Determine the tool's path if not already set.
	if utils.IsZero(t.URL) {
		if err := t.Hints.Parse(); err != nil {
			return result.WithFailed(fmt.Sprintf("parsing hints: %s", err))
		}

		if err := populator.URL(t.Name, nil, t.Version.Version, match.Requirements{
			Platform: t.Platform,
			Hints:    *t.Hints.Reduced(),
		}); err != nil {
			return result.WithFailed(fmt.Sprintf("getting url: %s", err))
		}

		t.URL = populator.Get("url")
	}

	// Update the URL to the template engine.
	tmpl.AddValue("URL", t.URL)

	// Append platform-specific extensions to aliases.
	for i, alias := range t.Aliases {
		if !file.File(alias).HasExtension() {
			t.Aliases[i] = alias + t.Platform.Extension.String()
		}
	}

	if !strings.HasSuffix(t.Exe.Name, t.Platform.Extension.String()) && !file.File(t.Exe.Name).HasExtension() {
		t.Exe.Name += t.Platform.Extension.String()
	}

	return outcome.Wrapped("requires download")
}

// Validate performs structural validation of the Tool's configuration using
// the validator package. Returns an error if validation fails.
func (t *Tool) Validate() error {
	if err := validator.Validate(t); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}

// CheckSkipConditions determines if a tool should be skipped based on its tags,
// strategy, and skip conditions. Returns a Result with the skip status and reason.
func (t *Tool) CheckSkipConditions(tags tags.IncludeTags) result.Result {
	res := result.WithSkipped("skipped")

	if !t.Tags.Include(tags.Include) {
		return res.Wrapped(fmt.Sprintf("does not have required tags: %v", tags.Include))
	}

	if !t.Tags.Exclude(tags.Exclude) {
		return res.Wrapped(fmt.Sprintf("has excluded tags: %v", tags.Exclude))
	}

	if skip, err := t.Skip.Evaluate(); err != nil {
		return result.WithFailed(fmt.Sprintf("evaluating skip conditions: %s", err))
	} else if skip.Has() {
		return res.Wrapped(fmt.Sprintf("condition: %q", skip[0].Reason))
	}

	if t.Strategy == strategy.None && t.Exists() {
		return res.Wrapped("already exists")
	}

	return result.WithOK("passed all conditions")
}

// Download retrieves and installs the tool using its configured source and installer.
// It handles progress tracking and executes any post-installation commands.
// Returns a Result indicating success or failure with detailed messages.
func (t *Tool) Download(progressListener getter.ProgressTracker) result.Result {
	installer, err := t.Source.Installer()
	if err != nil {
		return result.WithFailed("getting installer").Wrap(err)
	}

	data := common.InstallData{
		Path:        t.URL,
		Name:        t.Name,
		Exe:         t.Exe.Name,
		Patterns:    *t.Exe.Patterns,
		Output:      t.Output,
		Aliases:     t.Aliases,
		Mode:        t.Mode.String(),
		Env:         t.Env,
		NoVerifySSL: t.NoVerifySSL,
	}

	// Pass the progress listener to the specific source's Install method
	output, _, err := installer.Install(data, progressListener)
	if err != nil {
		return result.WithFailed("installing tool").Wrap(err).Wrapped(output)
	}

	// Execute post-installation commands if any exist
	if len(t.Commands.Commands) > 0 {
		if output, err := t.Commands.Run(t.Env); err != nil {
			return result.WithFailed("executing post-installation commands").Wrap(err).Wrapped(output)
		}
	}

	return result.WithOK("installed successfully")
}
