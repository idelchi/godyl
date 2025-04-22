package common

import (
	"fmt"

	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/defaults"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/utils"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/files"
)

func Common(cfg *config.Config, embedded config.Embedded, args []string) (tools.Tools, *logger.Logger, error) {
	lvl, err := logger.LevelString(cfg.Root.Log)
	if err != nil {
		return nil, nil, fmt.Errorf("parsing log level: %w", err)
	}

	// Set the tools file if provided as an argument
	if len(args) > 0 {
		cfg.Tool.Tools = files.New("", args...)
	} else {
		cfg.Tool.Tools = files.New(".", "tools.yml")
	}

	log := logger.New(lvl)

	// Load defaults
	defaults, err := defaults.Load(cfg.Root.Defaults, embedded, *cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("loading defaults: %w", err)
	}

	toolsList := tools.Tools{}

	// Load tools
	for _, file := range cfg.Tool.Tools {
		tools, err := utils.LoadTools(file, defaults, cfg.Root.Default)
		if err != nil {
			return nil, nil, fmt.Errorf("loading tools: %w", err)
		}

		log.Info("loaded %d tools from %q", len(tools), file)

		toolsList = append(toolsList, tools...)
	}

	return toolsList, log, nil
}
