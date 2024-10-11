package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/idelchi/godyl/internal/folder"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/inconshreveable/go-update"
)

func doUpdate(file string) error {
	body, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer body.Close()
	if err := update.Apply(body, update.Options{}); err != nil {
		return err
	}
	return err
}

func updater(cfg Config) error {
	if cfg.Update.Update {
		fmt.Printf("Updating godyl with strategy: %v\n", cfg.Update.Strategy)
		if err := cfg.Defaults.Defaults(); err != nil {
			return fmt.Errorf("Error setting defaults: %w", err)
		}

		var dir folder.Folder
		if err := dir.CreateRandomInTempDir(); err != nil {
			return fmt.Errorf("Error creating temporary directory: %w", err)
		}

		tool := tools.Tool{
			Output: dir.Path(),
			Name:   "idelchi/godyl",
			Source: sources.Source{
				Type: "github",
			},
		}
		tool.ApplyDefaults(cfg.Defaults)

		if err := tool.Resolve(nil, nil); err != nil {
			return fmt.Errorf("Error resolving tool: %w", err)
		}

		if output, _, err := tool.Download(); err != nil {
			return fmt.Errorf("Error downloading tool: %w: %s", err, output)
		}

		fmt.Printf("Downloading %q from %q\n", tool.Name, tool.Path)

		if err := doUpdate(filepath.Join(tool.Output, "godyl")); err != nil {
			return fmt.Errorf("Error updating godyl: %w", err)
		}

		fmt.Println("godyl updated successfully")

		os.Exit(0)

	}

	return nil
}
