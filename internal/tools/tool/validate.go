// Package tool provides core functionality for managing tool configurations.
package tool

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/result"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/folder"
	"github.com/idelchi/godyl/pkg/utils"
	"github.com/idelchi/godyl/pkg/validator"
)

// Resolve attempts to resolve the tool's source and strategy based on the provided tags.
// It handles environment variables, fallbacks, templating, and validation of the tool's
// configuration. Returns a Result indicating success or failure with detailed messages.
func (t *Tool) Resolve(tags tags.IncludeTags, dry bool) result.Result {
	// Load environment variables from the system.
	t.Env.Merge(env.FromEnv())

	// Expand and set the output folder path.
	t.Output = folder.New(t.Output).Expanded().Path()

	// Expand environment variables.
	t.Env.Expand()

	if err := t.TemplateFirst(); err != nil {
		return result.WithFailed("templating first").Wrap(err)
	}

	// Must pass atleast after first set of templates.
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
		if err != nil {
			return result.WithFailed(fmt.Sprintf("getting populator: %s", err))
		}

		// Initialize the installer.
		if err := populator.Initialize(t.Name); err != nil {
			return result.WithFailed(fmt.Sprintf("initializing populator: %s", err))
		}

		// Apply templating to the tool's fields.
		utils.SetIfZeroValue(&t.Exe.Name, populator.Get("exe"))
		utils.SetIfZeroValue(&t.Exe.Name, t.Name)

		// Re-check skip conditions after applying templates.
		if res = t.CheckSkipConditions(tags); !res.IsOK() {
			return res
		}

		if dry {
			break
		}

		if res = t.resolve(populator); res.IsFailed() {
			continue // Move on to the next fallback.
		}

		return res
	}

	return res
}

func (t *Tool) resolve(populator sources.Populator) result.Result {
	// Retrieve the tool's version from the installer if it is not already set.
	if utils.IsZeroValue(t.Version.Version) {
		if err := populator.Version(t.Name); err != nil {
			return result.WithFailed(fmt.Sprintf("getting version: %s", err))
		}
	}

	utils.SetIfZeroValue(&t.Version.Version, populator.Get("version"))

	if err := t.TemplateLast(); err != nil {
		return result.WithFailed(fmt.Sprintf("templating last: %s", err))
	}

	t.Extensions = t.Extensions.Compacted()
	t.Aliases = t.Aliases.Compacted()

	// Determine the tool's path if not already set.
	if utils.IsZeroValue(t.URL) {
		if err := t.Hints.Parse(); err != nil {
			return result.WithFailed(fmt.Sprintf("parsing hints: %s", err))
		}

		hints := t.Hints
		hints.Add(t.Extensions.ToHint())

		if err := populator.Path(t.Name, nil, t.Version.Version, match.Requirements{
			Platform: t.Platform,
			Hints:    hints,
		}); err != nil {
			return result.WithFailed(fmt.Sprintf("getting path: %s", err))
		}
	}

	utils.SetIfZeroValue(&t.URL, populator.Get("path"))

	// Append platform-specific file extension to the executable name.
	if !strings.HasSuffix(t.Exe.Name, t.Platform.Extension.String()) {
		t.Exe.Name += t.Platform.Extension.String()
	}

	// Set patterns for finding the executable.
	utils.SetSliceIfZero(&t.Exe.Patterns, fmt.Sprintf("^%s$", t.Exe.Name))

	// Append platform-specific extensions to aliases.
	for i, alias := range t.Aliases {
		t.Aliases[i] = alias + t.Platform.Extension.String()
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
		Patterns:    t.Exe.Patterns,
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
