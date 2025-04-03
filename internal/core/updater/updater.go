// Package updater provides functionality for updating tools and managing update strategies.
package updater

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/inconshreveable/go-update"

	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/folder"
	"github.com/idelchi/godyl/pkg/logger"
)

// Updater handles tool self-updating functionality.
type Updater struct {
	defaults    tools.Defaults
	noVerifySSL bool
	log         *logger.Logger
	template    []byte // Used for Windows cleanup batch script
}

// Versions contains the current and requested versions of the tool.
type Versions struct {
	Current   string
	Requested string
	Pre       bool
}

// New creates a new Updater with the specified configuration.
func New(defaults tools.Defaults, noVerifySSL bool, template []byte, level logger.Level) *Updater {
	return &Updater{
		defaults:    defaults,
		noVerifySSL: noVerifySSL,
		log:         logger.New(level),
		template:    template,
	}
}

// Update performs the self-update process for the godyl tool.
func (u *Updater) Update(versions Versions) error {
	tool, currentVersion, err := u.prepareToolInfo(versions)
	if err != nil {
		return err
	}

	// Skip if already up to date
	if tool.Version.Version == currentVersion {
		u.log.Info("godyl (%v) is already up-to-date", currentVersion)
		return nil
	}

	u.log.Info("Update requested from %q -> %q", currentVersion, tool.Version.Version)

	// TODO(Idelchi): Use `dry` flag here if set.

	return u.performUpdate(tool)
}

// prepareToolInfo gathers information about the current binary and prepares the tool configuration.
func (u *Updater) prepareToolInfo(versions Versions) (tools.Tool, string, error) {
	// Get path and version from build info
	path := "idelchi/godyl" // Default

	if info, ok := debug.ReadBuildInfo(); ok {
		path = strings.TrimPrefix(info.Main.Path, "github.com/")
		if versions.Current == "" {
			versions.Current = info.Main.Version
		}
	}

	// Create tool configuration
	tool := tools.Tool{
		Name: path,
		Version: tools.Version{
			Version:  versions.Requested,
			Patterns: []string{`.*?(\d+\.\d+\.\d+(?:-beta)?).*`},
		},
		Source: sources.Source{
			Type: sources.GITHUB,
			Github: github.GitHub{
				Pre: versions.Pre,
			},
		},
		Strategy:    tools.Upgrade,
		NoVerifySSL: u.noVerifySSL,
	}

	// Apply defaults and resolve configuration
	tool.ApplyDefaults(u.defaults)

	if err := tool.Resolve(nil, nil); err != nil && !(errors.Is(err, tools.ErrRequiresUpdate) || errors.Is(err, tools.ErrUpToDate)) {
		return tool, versions.Current, fmt.Errorf("resolving tool: %w", err)
	}

	return tool, versions.Current, nil
}

// performUpdate downloads and applies the update.
func (u *Updater) performUpdate(tool tools.Tool) error {
	// Download the tool to a temporary directory
	outputDir, err := u.downloadTool(tool)
	if err != nil {
		return err
	}

	// Clean up the temporary directory when done
	defer func() {
		if err := folder.New(outputDir).Remove(); err != nil {
			u.log.Warn("Failed to remove temporary folder: %v", err)
		}
	}()

	// Replace the existing binary with the newly downloaded version
	newBinaryPath := filepath.Join(outputDir, tool.Exe.Name)
	if err := u.replaceBinary(newBinaryPath); err != nil {
		return fmt.Errorf("replacing binary: %w", err)
	}

	// Handle platform-specific cleanup
	if IsWindows() {
		if err := u.cleanupWindows(); err != nil {
			return fmt.Errorf("windows cleanup: %w", err)
		}

		u.log.Debug("Windows cleanup completed successfully")
	}

	u.log.Info("Godyl updated successfully")
	return nil
}

// downloadTool downloads the tool to a temporary directory.
func (u *Updater) downloadTool(tool tools.Tool) (string, error) {
	// Create a temporary directory based on the platform
	dir, err := u.createTempDir()
	if err != nil {
		return "", err
	}

	// Configure the tool for download
	tool.Output = dir

	// Download the tool
	_, msg, err := tool.Download()
	if err != nil {
		return "", fmt.Errorf("downloading tool: %w: %s", err, msg)
	}

	return tool.Output, nil
}

// createTempDir creates an appropriate temporary directory based on the platform.
func (u *Updater) createTempDir() (string, error) {
	if IsWindows() {
		// On Windows, create temp dir in the same directory as the executable
		exePath, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("getting executable path: %w", err)
		}

		dir, err := tmp.GodylCreateRandomDirIn(folder.New(file.New(exePath).Dir()))
		if err != nil {
			return "", fmt.Errorf("creating temporary directory: %w", err)
		}

		return dir.Path(), nil

	}
	// On other platforms, use system temp directory
	dir, err := tmp.GodylCreateRandomDir()
	if err != nil {
		return "", fmt.Errorf("creating temporary directory: %w", err)
	}

	return dir.Path(), nil
}

// replaceBinary replaces the current executable with the new binary.
func (u *Updater) replaceBinary(newBinaryPath string) error {
	file, err := os.Open(filepath.Clean(newBinaryPath))
	if err != nil {
		return fmt.Errorf("opening new binary: %w", err)
	}
	defer file.Close()

	options := update.Options{}
	if err := update.Apply(file, options); err != nil {
		return fmt.Errorf("applying update: %w", err)
	}

	return nil
}

// cleanupWindows handles Windows-specific cleanup after an update.
func (u *Updater) cleanupWindows() error {
	return createAndRunCleanupScript(u.template, u.log)
}
