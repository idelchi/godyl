package updater

import (
	"fmt"
	"html/template"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/idelchi/godyl/pkg/file"
	"github.com/idelchi/godyl/pkg/logger"
)

type cleanupData struct {
	OldBinary string
	BatchFile string
	Folder    string
	LogFile   string
}

// winCleanup handles Windows-specific cleanup after an update.
func winCleanup(cleanupTemplate []byte) error {
	log := logger.New(logger.INFO)
	log.Info("Issuing a delete command for the old godyl binary")

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("getting executable path: %w", err)
	}

	exeDir := file.NewFile(exePath).Dir()

	var folder file.Folder
	if err := folder.CreateRandomInTempDir(); err != nil {
		return fmt.Errorf("creating temporary directory: %w", err)
	}

	oldBinary := file.NewFile(exeDir.Path(), ".godyl.exe.old")
	batchFile := file.NewFile(folder.Path(), "cleanup.bat")
	logFile := file.NewFile(folder.Path(), "cleanup_debug.log")

	log.Info("Batch file stored in: %s", batchFile.Path())

	// Read and parse the template
	tmpl, err := template.New("cleanup").Parse(string(cleanupTemplate))
	if err != nil {
		return fmt.Errorf("parsing cleanup template: %w", err)
	}

	// Create the batch file
	batchFilePath := filepath.Clean(batchFile.Path())

	batchFileHandle, err := os.Create(batchFilePath)
	if err != nil {
		return fmt.Errorf("creating batch file: %w", err)
	}

	defer batchFileHandle.Close()

	// Execute the template with the data
	data := cleanupData{
		OldBinary: oldBinary.Path(),
		BatchFile: batchFile.Path(),
		Folder:    folder.Path(),
		LogFile:   logFile.Path(),
	}

	if err := tmpl.Execute(batchFileHandle, data); err != nil {
		return fmt.Errorf("executing cleanup template: %w", err)
	}

	// Fire and forget, run minimized
	safeCmd := "cmd"
	cmdArgs := []string{"/C", "start", "/MIN", batchFilePath}

	cmd := exec.Command(safeCmd, cmdArgs...) //nolint:gosec	// TODO(Idelchi): Keep this in mind.
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting cleanup script: %w", err)
	}

	return nil
}
