package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/idelchi/godyl/internal/folder"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/sources"
	"github.com/inconshreveable/go-update"
)

type GodylUpdater struct {
	Strategy tools.Strategy

	Defaults tools.Defaults
}

func (gu GodylUpdater) Update() error {
	if gu.Strategy == tools.None {
		gu.Strategy = tools.Upgrade
	}

	fmt.Printf("Updating godyl with strategy: %q\n", gu.Strategy)

	path := "idelchi/godyl"
	info, ok := debug.ReadBuildInfo()
	if ok {
		path = strings.TrimPrefix(info.Main.Path, "github.com/")
	}

	tool := tools.Tool{
		Name: path,
		Source: sources.Source{
			Type: sources.GITHUB,
		},
		Strategy: gu.Strategy,
	}

	tool.ApplyDefaults(gu.Defaults)

	output, err := gu.Get(tool)
	if err != nil {
		return fmt.Errorf("geting godyl: %w", err)
	}

	if err := gu.Replace(filepath.Join(output, "godyl")); err != nil {
		return fmt.Errorf("replacing godyl: %w", err)
	}

	fmt.Println("godyl updated successfully")

	return nil
}

func (gu GodylUpdater) Replace(path string) error {
	body, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("opening file %q: %w", path, err)
	}
	defer body.Close()
	if err := update.Apply(body, update.Options{}); err != nil {
		return err
	}
	return err
}

func (gu GodylUpdater) Get(tool tools.Tool) (string, error) {
	var dir folder.Folder
	if err := dir.CreateRandomInTempDir(); err != nil {
		return "", fmt.Errorf("creating temporary directory: %w", err)
	}

	tool.Output = dir.Path()

	if err := tool.Resolve(nil, nil); err != nil {
		return "", fmt.Errorf("resolving tool: %w", err)
	}

	if output, msg, err := tool.Download(); err != nil {
		return "", fmt.Errorf("downloading tool: %w: %s", err, output)
	} else {
		fmt.Println(msg)
		fmt.Println(output)
	}

	fmt.Printf("Downloading %q from %q\n", tool.Name, tool.Path)

	return tool.Output, nil
}
