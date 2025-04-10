package updater

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"

	"github.com/idelchi/godyl/internal/tmp"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
)

// cleanupData contains data needed for the cleanup batch script template.
type cleanupData struct {
	OldBinary string
	BatchFile string
	Folder    string
	LogFile   string
}

// createAndRunCleanupScript creates and executes a Windows cleanup batch script.
func createAndRunCleanupScript(templateContent []byte, log *logger.Logger) error {
	log.Debug("Issuing a delete command for the old godyl binary")

	// Get the current executable path
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	log.Debug("Executable path: %q", exePath)

	// Create a temporary folder for cleanup files
	folder, err := tmp.GodylCreateRandomDir()
	if err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	// Prepare file paths
	oldBinary := file.New(file.New(exePath).Dir(), ".godyl.exe.old")
	batchFile := file.New(folder.Path(), "cleanup.bat")
	logFile := file.New(folder.Path(), "cleanup_debug.log")

	log.Debug("Batch file stored in: %s", batchFile.Path())

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

// createBatchFile creates a batch file from the provided template and data.
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

// executeScript runs the cleanup script in a minimized window.
func executeScript(scriptPath string) error {
	cmd := exec.Command("cmd", "/C", "start", "/MIN", scriptPath) //nolint:gosec
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting cleanup script: %w", err)
	}
	return nil
}
