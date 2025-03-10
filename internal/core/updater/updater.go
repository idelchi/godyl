// Package updater provides functionality for updating tools and managing update strategies.
package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/inconshreveable/go-update"

	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
)

func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// ToolDownloader is responsible for downloading tools.
type ToolDownloader interface {
	Download(tool tools.Tool) (string, error)
}

// BinaryReplacer is responsible for replacing the current binary with a new one.
type BinaryReplacer interface {
	Replace(path string) error
}

// Updater is responsible for updating the godyl tool using the specified update strategy and defaults.
type Updater struct {
	Defaults    tools.Defaults // Defaults holds tool-specific default values for the update process.
	NoVerifySSL bool           // NoVerifySSL disables SSL verification for the update process.
	Template    []byte

	downloader ToolDownloader
	replacer   BinaryReplacer
	log        *logger.Logger
}

// NewUpdater creates a new Updater with the specified strategy and defaults.
func NewUpdater(defaults tools.Defaults, noVerifySSL bool, template []byte) *Updater {
	return &Updater{
		Defaults:    defaults,
		NoVerifySSL: noVerifySSL,
		Template:    template,
		downloader:  &DefaultDownloader{},
		replacer:    &DefaultReplacer{},
		log:         logger.New(logger.INFO),
	}
}

// DefaultDownloader is the default implementation of ToolDownloader.
type DefaultDownloader struct{}

// DefaultReplacer is the default implementation of BinaryReplacer.
type DefaultReplacer struct{}

// Update performs the update process for the godyl tool, applying the specified strategy.
func (u *Updater) Update() error {
	// Determine the tool path from build info, defaulting to "idelchi/godyl" if not available.
	path := "idelchi/godyl"
	info, ok := debug.ReadBuildInfo()

	var version string

	if ok {
		path = strings.TrimPrefix(info.Main.Path, "github.com/")
		version = info.Main.Version
	}

	// Create a new Tool object with the appropriate strategy and source.
	tool := tools.Tool{
		Name: path,
		Source: sources.Source{
			Type: sources.GITHUB,
		},
		Strategy:    tools.Upgrade,
		NoVerifySSL: u.NoVerifySSL,
	}

	// Apply any default values to the tool.
	tool.ApplyDefaults(u.Defaults)

	if err := tool.Resolve(nil, nil); err != nil {
		return fmt.Errorf("resolving tool: %w", err)
	}

	// Check if update is needed based on strategy
	if u.shouldUpdate(tool, version) {
		return u.performUpdate(tool)
	}

	return nil
}

// shouldUpdate determines if an update should be performed based on the strategy and versions.
func (u *Updater) shouldUpdate(tool tools.Tool, currentVersion string) bool {
	if tool.Version.Version == currentVersion {
		u.log.Info("godyl (%v) is already up-to-date", currentVersion)

		return false
	}

	u.log.Info("Update requested from %q -> %q", currentVersion, tool.Version.Version)

	return true
}

// performUpdate downloads and applies the update.
func (u *Updater) performUpdate(tool tools.Tool) error {
	// Download the tool
	output, err := u.Get(tool)
	if err != nil {
		return fmt.Errorf("getting godyl: %w", err)
	}

	// Clean up the temporary directory when done
	defer func() {
		folder := file.Folder(output)
		if err := folder.Remove(); err != nil {
			u.log.Warn("Failed to remove temporary folder: %v", err)
		}
	}()

	// Replace the existing godyl binary with the newly downloaded version
	if err := u.Replace(filepath.Join(output, tool.Exe.Name)); err != nil {
		return fmt.Errorf("replacing godyl: %w", err)
	}

	// Perform platform-specific cleanup
	if IsWindows() {
		if err := winCleanup(u.Template); err != nil {
			return fmt.Errorf("issuing delete command: %w", err)
		}
	}

	u.log.Info("Godyl updated successfully")

	return nil
}

// Replace applies the new godyl binary by replacing the current executable with the downloaded one.
func (u *Updater) Replace(path string) error {
	if err := u.replacer.Replace(path); err != nil {
		return fmt.Errorf("replacing binary: %w", err)
	}

	return nil
}

// Replace implements the BinaryReplacer interface.
func (r *DefaultReplacer) Replace(path string) error {
	body, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("opening file %q: %w", path, err)
	}
	defer body.Close()

	options := update.Options{}

	// Removed empty block - uncomment if needed in the future
	// if IsWindows() {
	//	options.OldSavePath = filepath.Join(filepath.Dir(path), ".godyl.exe.old")
	// }

	if err := update.Apply(body, options); err != nil {
		return fmt.Errorf("applying update: %w", err)
	}

	return nil
}

// Get downloads the tool based on its source, placing it in a temporary directory, and returns the output path.
func (u *Updater) Get(tool tools.Tool) (string, error) {
	path, err := u.downloader.Download(tool)
	if err != nil {
		return "", fmt.Errorf("downloading tool: %w", err)
	}

	return path, nil
}

// Download implements the ToolDownloader interface.
func (d *DefaultDownloader) Download(tool tools.Tool) (string, error) {
	// Create a temporary directory to store the downloaded tool
	var dir file.Folder

	// For Windows, get the directory of the current executable
	if IsWindows() {
		current, err := os.Executable()
		if err != nil {
			return "", fmt.Errorf("getting current executable: %w", err)
		}

		folder := filepath.Dir(current)
		if err := dir.CreateRandomInDir(folder); err != nil {
			return "", fmt.Errorf("creating temporary directory: %w", err)
		}
	} else {
		if err := dir.CreateRandomInTempDir(); err != nil {
			return "", fmt.Errorf("creating temporary directory: %w", err)
		}
	}

	tool.Output = dir.Path()

	// Resolve any dependencies or settings for the tool
	if err := tool.Resolve(nil, nil); err != nil {
		return "", fmt.Errorf("resolving tool: %w", err)
	}

	// Download the tool and capture any messages or errors
	if output, msg, err := tool.Download(); err != nil {
		return "", fmt.Errorf("downloading tool: %w: %s: %s", err, output, msg)
	}

	// Using logger would be ideal, but DefaultDownloader doesn't have access to it
	// For now, we'll just return the output path and let the caller log anything needed

	return tool.Output, nil
}
