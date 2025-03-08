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
)

// UpdateStrategy defines how updates are applied.
type UpdateStrategy string

const (
	// None means no updates will be applied.
	None UpdateStrategy = "none"

	// Upgrade means only newer versions will be applied.
	Upgrade UpdateStrategy = "upgrade"

	// Force means updates will be applied regardless of version.
	Force UpdateStrategy = "force"
)

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
	Strategy    tools.Strategy // Strategy defines how updates are applied (e.g., Upgrade, Downgrade, None).
	Defaults    tools.Defaults // Defaults holds tool-specific default values for the update process.
	NoVerifySSL bool           // NoVerifySSL disables SSL verification for the update process.
	Template    []byte

	downloader ToolDownloader
	replacer   BinaryReplacer
}

// NewUpdater creates a new Updater with the specified strategy and defaults.
func NewUpdater(strategy tools.Strategy, defaults tools.Defaults, noVerifySSL bool) *Updater {
	return &Updater{
		Strategy:    strategy,
		Defaults:    defaults,
		NoVerifySSL: noVerifySSL,
		downloader:  &DefaultDownloader{},
		replacer:    &DefaultReplacer{},
	}
}

// DefaultDownloader is the default implementation of ToolDownloader.
type DefaultDownloader struct{}

// DefaultReplacer is the default implementation of BinaryReplacer.
type DefaultReplacer struct{}

// Update performs the update process for the godyl tool, applying the specified strategy.
func (u *Updater) Update(version string) error {
	// Set default strategy if none is provided.
	if u.Strategy == tools.None {
		u.Strategy = tools.Upgrade
	}

	// Determine the tool path from build info, defaulting to "idelchi/godyl" if not available.
	path := "idelchi/godyl"
	info, ok := debug.ReadBuildInfo()
	if ok {
		path = strings.TrimPrefix(info.Main.Path, "github.com/")
	}

	// Create a new Tool object with the appropriate strategy and source.
	tool := tools.Tool{
		Name: path,
		Source: sources.Source{
			Type: sources.GITHUB,
		},
		Strategy:    u.Strategy,
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
	if u.Strategy == tools.Force {
		fmt.Println("Forcing update...")
		return true
	}

	if tool.Version.Version == currentVersion {
		fmt.Printf("godyl (%v) is already up-to-date\n", currentVersion)
		return false
	}

	fmt.Printf("Update requested from %q -> %q\n", currentVersion, tool.Version.Version)
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
		folder.Remove()
	}()

	// Replace the existing godyl binary with the newly downloaded version
	if err := u.Replace(filepath.Join(output, tool.Exe.Name)); err != nil {
		return fmt.Errorf("replacing godyl: %w", err)
	}

	// Perform platform-specific cleanup
	if runtime.GOOS == "windows" {
		if err := winCleanup(u.Template); err != nil {
			return fmt.Errorf("issuing delete command: %w", err)
		}
	}

	fmt.Println("Godyl updated successfully")
	return nil
}

// Replace applies the new godyl binary by replacing the current executable with the downloaded one.
func (u *Updater) Replace(path string) error {
	return u.replacer.Replace(path)
}

// Replace implements the BinaryReplacer interface.
func (r *DefaultReplacer) Replace(path string) error {
	body, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file %q: %w", path, err)
	}
	defer body.Close()

	options := update.Options{}
	if runtime.GOOS == "windows" {
		// options.OldSavePath = filepath.Join(filepath.Dir(path), ".godyl.exe.old")
	}

	if err := update.Apply(body, options); err != nil {
		return err
	}

	return nil
}

// Get downloads the tool based on its source, placing it in a temporary directory, and returns the output path.
func (u *Updater) Get(tool tools.Tool) (string, error) {
	return u.downloader.Download(tool)
}

// Download implements the ToolDownloader interface.
func (d *DefaultDownloader) Download(tool tools.Tool) (string, error) {
	// Create a temporary directory to store the downloaded tool
	var dir file.Folder

	// For Windows, get the directory of the current executable
	if runtime.GOOS == "windows" {
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

	fmt.Printf("Downloading %q from %q\n", tool.Name, tool.Path)

	return tool.Output, nil
}
