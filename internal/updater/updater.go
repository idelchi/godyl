// Package updater provides functionality for updating tools and managing update strategies.
package updater

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/inconshreveable/go-update"

	"github.com/idelchi/godyl/internal/config/root"
	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/internal/tools/mode"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/idelchi/godyl/internal/tools/sources/github"
	"github.com/idelchi/godyl/internal/tools/strategy"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/internal/tools/version"
	"github.com/idelchi/godyl/pkg/download/progress"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/idelchi/godyl/pkg/path/folder"
	vversion "github.com/idelchi/godyl/pkg/version"
)

type Godyl struct {
	Tool    *tool.Tool
	Version string
}

func (g *Godyl) IsUpToDate() bool {
	return vversion.Equal(g.Version, g.Tool.Version.Version)
}

func NewGodyl(v string, cfg *root.Config) Godyl {
	path := "idelchi/godyl" // Default

	if info, ok := debug.ReadBuildInfo(); ok {
		path = strings.TrimPrefix(info.Main.Path, "github.com/")

		if v == "" {
			v = info.Main.Version
		}
	}

	return Godyl{
		Version: v,
		Tool: &tool.Tool{
			Name: path,
			Version: version.Version{
				Version:  cfg.Update.Version,
				Patterns: &version.Patterns{`.*?(\d+\.\d+\.\d+(?:-beta)?).*`},
			},
			Mode: mode.Extract,
			Source: sources.Source{
				Type: sources.GITHUB,
				GitHub: github.GitHub{
					Pre: cfg.Update.Pre,
				},
			},
			NoCache:  true,
			Strategy: cfg.Common.Strategy,
		},
	}
}

// Updater manages the self-update process for the godyl tool.
type Updater struct {
	godyl    *Godyl
	log      *logger.Logger
	template []byte
}

// New creates a new Updater instance with the provided configuration.
// Takes default settings, SSL verification flag, cleanup script template, and logger.
func New(godyl *Godyl, template []byte, log *logger.Logger) *Updater {
	return &Updater{
		godyl:    godyl,
		log:      log,
		template: template,
	}
}

func (u *Updater) IsForced() bool {
	return u.godyl.Tool.Strategy == strategy.Force
}

// Update performs the self-update process for the godyl tool.
// Downloads the new version, replaces the current binary, and handles
// platform-specific cleanup. Returns an error if any step fails.
func (u *Updater) Update(check bool) error {
	res := u.godyl.Tool.Resolve(tags.IncludeTags{}, tool.WithUpUntilVersion())
	if res.AsError() != nil {
		return res.AsError()
	}

	if u.godyl.IsUpToDate() && !u.IsForced() {
		u.log.Infof("godyl (%v) is already up-to-date", u.godyl.Version)

		return nil
	}

	if check {
		u.log.Infof("A new version %q is available! You are at %q", u.godyl.Tool.Version.Version, u.godyl.Version)

		if body := u.godyl.Tool.GetPopulator().Get("body"); body != "" {
			u.log.Info("")
			u.log.Info(strings.TrimSpace(body))
		}

		u.log.Info("")
		u.log.Info("You can update with:")

		if u.godyl.Tool.Source.GitHub.Pre {
			u.log.Info("  godyl update --pre")
		} else {
			u.log.Info("  godyl update")
		}

		return nil
	}

	u.log.Infof("Update requested from %q -> %q", u.godyl.Version, u.godyl.Tool.Version.Version)

	return u.performUpdate(u.godyl.Tool)
}

// PerformUpdate downloads the new version and applies the update.
// Handles temporary file management and platform-specific cleanup.
func (u *Updater) performUpdate(tool *tool.Tool) error {
	if res := tool.Resolve(tags.IncludeTags{}); !res.IsOK() {
		return res.AsError()
	}

	// Download the tool to a temporary directory
	outputDir, err := u.downloadTool(tool)
	if err != nil {
		return err
	}

	// Clean up the temporary directory when done
	defer func() {
		if err := folder.New(outputDir).Remove(); err != nil {
			u.log.Warnf("Failed to remove temporary folder: %v", err)
		}
	}()

	// Replace the existing binary with the newly downloaded version
	newBinaryPath := filepath.Join(outputDir, tool.Exe.Name)
	if err := u.replaceBinary(newBinaryPath); err != nil {
		return fmt.Errorf("replacing binary: %w", err)
	}

	// Handle platform-specific cleanup
	if IsWindows() && u.template != nil {
		u.log.Debug("Performing Windows cleanup")

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
func (u *Updater) downloadTool(tool *tool.Tool) (string, error) {
	// Create a temporary directory based on the platform
	dir, err := u.createTempDir()
	if err != nil {
		return "", err
	}

	// Configure the tool for download
	tool.Output = dir

	// Download the tool, passing the progress tracker
	// := progress.New()
	// tracker.Start()

	res := tool.Download(progress.NewNoop())

	// tracker.Wait()

	if err := res.AsError(); err != nil {
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
