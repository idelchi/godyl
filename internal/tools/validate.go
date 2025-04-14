package tools

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/hashicorp/go-getter/v2"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/common"
	"github.com/idelchi/godyl/pkg/env"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
	"github.com/idelchi/godyl/pkg/utils"
)

// Resolve attempts to resolve the tool's source and strategy based on the provided tags.
// It handles environment variables, fallbacks, templating, and validation of the tool's
// configuration. Returns a Result indicating success or failure with detailed messages.
func (t *Tool) Resolve(tags IncludeTags) Result {
	// Load environment variables from the system.
	t.Env.Merge(env.FromEnv())

	// Expand and set the output folder path.
	t.Output = folder.New(t.Output).Expanded().Path()

	// Save the path for templating later.
	path := t.Path

	// Expand environment variables.
	t.Env.Expand()

	if err := t.TemplateFirst(); err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("templating first: %s", err)}
	}

	// Must atleast after first set of templates.
	if err := t.Validate(); err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("validating config: %s", err)}
	}

	var result Result
	// Try resolving with each fallback in order.
	for _, fallback := range t.Fallbacks.Build(t.Source.Type) {
		if result = t.resolve(fallback, path, tags.Include, tags.Exclude); !result.Successful() {
			continue // Move on to the next fallback.
		}

		break
	}

	return result
}

// CheckSkipConditions determines if a tool should be skipped based on its tags,
// strategy, and skip conditions. Returns a Result with the skip status and reason.
func (t *Tool) CheckSkipConditions(withTags, withoutTags []string) Result {
	result := Result{Status: Skipped, Message: "skipped due to"}

	if !t.Tags.Has(withTags) {
		return result.Wrapped(fmt.Sprintf("does not have required tags: %v", withTags))
	}

	if !t.Tags.HasNot(withoutTags) {
		return result.Wrapped(fmt.Sprintf("has excluded tags: %v", withoutTags))
	}

	if skip, err := t.Skip.Evaluate(); err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("evaluating skip conditions: %s", err)}
	} else if skip.Has() {
		return result.Wrapped(fmt.Sprintf("condition: %q", skip[0].Reason))
	}

	if t.Strategy == None && t.Exists() {
		return result.Wrapped("already exists")
	}

	return Result{Status: OK, Message: "passed all conditions"}
}

// Resolve attempts to resolve a tool using a specific fallback source type.
// It handles executable details, version information, path resolution, and
// platform-specific configurations. Returns a Result indicating success or failure.
//
//nolint:cyclop,funlen 	// TODO(Idelchi): Refactor this function to reduce cyclomatic complexity.
func (t *Tool) resolve(fallback sources.Type, path string, withTags, withoutTags []string) Result {
	// Append the tool's name as a tag.
	t.Tags.Append(t.Name)

	// Set the source type to the current fallback.
	t.Source.Type = fallback

	// Get the installer for the current source type.
	populator, err := t.Source.Installer()
	if err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("getting installer: %s", err)}
	}

	// Initialize the installer.
	if err := populator.Initialize(t.Name); err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("initializing installer: %s", err)}
	}

	// Retrieve executable details from the installer.
	if err := populator.Exe(); err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("getting executable details: %s", err)}
	}

	// Apply templating to the tool's fields.
	utils.SetIfZeroValue(&t.Exe.Name, populator.Get("exe"))
	utils.SetIfZeroValue(&t.Exe.Name, t.Name)

	// Re-check skip conditions after applying templates.
	if result := t.CheckSkipConditions(withTags, withoutTags); !result.Successful() {
		return result
	}

	// Retrieve the tool's version from the installer if it is not already set.
	if utils.IsZeroValue(t.Version.Version) {
		if err := populator.Version(t.Name); err != nil {
			return Result{Status: Failed, Message: fmt.Sprintf("getting version: %s", err)}
		}
	}

	utils.SetIfZeroValue(&t.Version.Version, populator.Get("version"))

	if err := t.TemplateLast(); err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("templating last: %s", err)}
	}

	t.Extensions = t.Extensions.Compacted()
	t.Aliases = t.Aliases.Compacted()

	// Determine the tool's path if not already set.
	if utils.IsZeroValue(t.Path) {
		if err := t.Hints.Parse(); err != nil {
			return Result{Status: Failed, Message: fmt.Sprintf("parsing hints: %s", err)}
		}

		hints := t.Hints
		hints.Add(t.Extensions.ToHint())

		if err := populator.Path(t.Name, nil, t.Version.Version, match.Requirements{
			Platform: t.Platform,
			Hints:    hints,
		}); err != nil {
			return Result{Status: Failed, Message: fmt.Sprintf("getting path: %s", err)}
		}
	}

	utils.SetIfZeroValue(&t.Path, populator.Get("path"))
	utils.SetIfZeroValue(&t.Path, path)

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

	// Attempt to upgrade the tool using the current strategy.
	if result := t.Strategy.Upgrade(t); !result.Successful() {
		return result
	}

	// Validate the tool's configuration.
	if err := t.Validate(); err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("validating config: %s", err)}
	}

	return Result{Status: OK, Message: "resolved successfully"}
}

// Validate performs structural validation of the Tool's configuration using
// the validator package. Returns an error if validation fails.
func (t *Tool) Validate() error {
	validate := validator.New()
	if err := validate.Struct(t); err != nil {
		return fmt.Errorf("validating config: %w", err)
	}

	return nil
}

// Exists checks if the tool's executable exists in the configured output path.
// Returns true if the file exists and is a regular file.
func (t *Tool) Exists() bool {
	f := file.New(t.Output, t.Exe.Name)

	return f.Exists() && f.IsFile()
}

// Download retrieves and installs the tool using its configured source and installer.
// It handles progress tracking and executes any post-installation commands.
// Returns a Result indicating success or failure with detailed messages.
func (t *Tool) Download(progressListener getter.ProgressTracker) Result {
	installer, err := t.Source.Installer()
	if err != nil {
		return Result{Status: Failed, Message: fmt.Sprintf("getting installer: %s", err)}
	}

	data := common.InstallData{
		Path:        t.Path,
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
		return Result{Status: Failed, Message: fmt.Sprintf("installing tool: %s", err)}.Wrapped(output)
	}

	// Execute post-installation commands if any exist
	if len(t.Commands.Commands) > 0 {
		if output, err := t.Commands.Exe(t.Env); err != nil {
			return Result{Status: Failed, Message: fmt.Sprintf("executing post-installation commands: %s", err)}.Wrapped(output)
		}
	}

	return Result{Status: OK, Message: "installed successfully"}
}
