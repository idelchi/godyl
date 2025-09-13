package updater

import (
	"context"
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/idelchi/godyl/internal/data"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
)

// CleanupData contains the paths and filenames needed for Windows cleanup.
// Used to populate the cleanup batch script template with the correct paths.
type cleanupData struct {
	// OldBinary is the path to the old executable to be removed.
	OldBinary string

	// BatchFile is the path to the cleanup batch script.
	BatchFile string

	// Folder is the temporary directory path for cleanup files.
	Folder string

	// LogFile is the path where cleanup logs will be written.
	LogFile string
}

// CreateAndRunCleanupScript handles Windows-specific cleanup after an update.
// Creates a batch script from the template, populates it with the necessary paths,
// and executes it in a minimized window. Returns an error if any step fails.
func createAndRunCleanupScript(templateContent []byte, log *logger.Logger) error {
	log.Debug("Issuing a delete command for the old godyl binary")

	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	log.Debugf("Executable path: %q", exePath)

	// Create a temporary folder for cleanup files
	folder, err := data.CreateUniqueDirIn()
	if err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	// Prepare file paths
	oldBinary := file.New(file.New(exePath).Dir(), ".godyl.exe.old")
	batchFile := file.New(folder.Path(), "cleanup.bat")
	logFile := file.New(folder.Path(), "cleanup_debug.log")

	log.Debugf("Batch file stored in: %s", batchFile.Path())

	// Create cleanup script
	if err := createBatchFile(templateContent, batchFile.Path(), cleanupData{
		OldBinary: oldBinary.Path(),
		BatchFile: batchFile.Path(),
		Folder:    folder.Path(),
		LogFile:   logFile.Path(),
	}); err != nil {
		return err
	}

	// Execute the cleanup script
	return executeScript(batchFile.Path())
}

// CreateBatchFile generates a cleanup batch script from the template.
// Takes the template content, output path, and cleanup data as input.
// Returns an error if the file cannot be created or the template fails.
func createBatchFile(templateContent []byte, batchFilePath string, data cleanupData) error {
	// Parse the template
	tmpl, err := template.New("cleanup").Parse(string(templateContent))
	if err != nil {
		return fmt.Errorf("parsing cleanup template: %w", err)
	}

	batchFile := file.New(batchFilePath)
	if err := batchFile.Create(); err != nil {
		return fmt.Errorf("creating batch file: %w", err)
	}

	// Create the batch file
	batchFileHandle, err := batchFile.OpenForWriting()
	if err != nil {
		return fmt.Errorf("creating batch file: %w", err)
	}
	defer batchFileHandle.Close()

	// Execute the template with the data
	if err := tmpl.Execute(batchFileHandle, data); err != nil {
		return fmt.Errorf("executing cleanup template: %w", err)
	}

	return nil
}

// ExecuteScript runs the cleanup batch script in a minimized window.
// Uses cmd.exe to start the script with minimal UI visibility.
func executeScript(scriptPath string) error {
	cmd := exec.CommandContext(context.Background(), "cmd", "/C", "start", "/MIN", scriptPath)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting cleanup script: %w", err)
	}

	return nil
}
