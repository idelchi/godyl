// Package updater provides functionality for updating tools and managing update strategies.
package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/inconshreveable/go-update"

	"github.com/idelchi/godyl/internal/cache/cache"
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/internal/ui/progress"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
)

// Updater manages the self-update process for the godyl tool.
type Updater struct {
	// Defaults contains default configuration values.
	defaults tools.Defaults

	// NoVerifySSL disables SSL certificate verification.
	noVerifySSL bool

	// Log is the logger instance for output.
	log *logger.Logger

	// Template contains the Windows cleanup script template.
	template []byte

	// Cache manages version caching.
	cache *cache.Cache
}

// Versions defines the version information for an update operation.
type Versions struct {
	Current   string
	Requested string
	Pre       bool
}

// New creates a new Updater instance with the provided configuration.
// Takes default settings, SSL verification flag, cleanup script template,
// logging level, and cache instance.
func New(defaults tools.Defaults, noVerifySSL bool, template []byte, level logger.Level, cache *cache.Cache) *Updater {
	return &Updater{
		defaults:    defaults,
		noVerifySSL: noVerifySSL,
		log:         logger.New(level),
		template:    template,
	}
}

// Update performs the self-update process for the godyl tool.
// Downloads the new version, replaces the current binary, and handles
// platform-specific cleanup. Returns an error if any step fails.
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

// PrepareToolInfo gathers information about the current binary and creates
// a tool configuration for the update. Returns the tool configuration,
// current version, and any errors encountered.
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
			GitHub: github.GitHub{
				Pre: versions.Pre,
			},
		},
		Strategy:    tools.Upgrade,
		NoVerifySSL: u.noVerifySSL,
	}

	// Apply defaults and resolve configuration
	tool.ApplyDefaults(u.defaults)

	if result := tool.Resolve(tools.IncludeTags{}); result.Unsuccessful() {
		return tool, versions.Current, result.Error()
	}

	return tool, versions.Current, nil
}

// PerformUpdate downloads the new version and applies the update.
// Handles temporary file management and platform-specific cleanup.
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

// DownloadTool retrieves the new version and stores it in a temporary directory.
// Sets up progress tracking for the download operation.
func (u *Updater) downloadTool(tool tools.Tool) (string, error) {
	// Create a temporary directory based on the platform
	dir, err := u.createTempDir()
	if err != nil {
		return "", err
	}

	// Configure the tool for download
	tool.Output = dir

	// Setup progress bar for download
	progressMgr := progress.NewProgressManager(progress.DefaultOptions())
	progressTracker := progressMgr.NewTracker(&tool)

	// Download the tool, passing the progress tracker
	result := tool.Download(progressTracker)

	// Wait for progress bar to finish rendering
	progressTracker.Wait()

	if err := result.Error(); err != nil {
		return "", fmt.Errorf("downloading tool: %w", err)
	}

	return tool.Output, nil
}

// CreateTempDir creates a temporary directory for the update process.
// Uses platform-specific logic to determine the appropriate location.
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

// ReplaceBinary replaces the current executable with the new version.
// Uses go-update library to handle the replacement process safely.
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

// CleanupWindows performs Windows-specific post-update cleanup operations.
// Creates and executes a cleanup script to handle file replacement.
func (u *Updater) cleanupWindows() error {
	return createAndRunCleanupScript(u.template, u.log)
}
