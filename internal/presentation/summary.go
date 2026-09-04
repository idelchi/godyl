package presentation

import (
	"errors"

	"github.com/idelchi/godyl/internal/match"
	"github.com/idelchi/godyl/internal/processor"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
)

// ShowConfig configures how results are presented.
type ShowConfig struct {
	// Verbose controls the amount of result detail shown.
	Verbose int
	// ErrorFile receives machine-readable errors when set.
	ErrorFile file.File
	// Suggest reports that results came from non-mutating suggestion mode.
	Suggest bool
}

// ShowSummary formats and displays the processing results.
func ShowSummary(summary processor.Summary, cfg ShowConfig, log *logger.Logger) {
	const tableMaxWidth = 100

	// Create table formatter
	tableFormatter := NewTableFormatter(TableConfig{
		MaxWidth: tableMaxWidth,
		Verbose:  cfg.Verbose > 0,
	})

	// Render table
	tableOutput := tableFormatter.RenderResults(summary.Results)

	if tableOutput == "" {
		log.Info("Nothing of interest to show")

		return
	}

	// Display results
	if cfg.Verbose > 0 {
		log.Info("")
		log.Info("Installation Summary:")
		log.Info(tableOutput)

		if suggestions := RenderSuggestions(summary.Results); suggestions != "" {
			log.Info("")
			log.Info("AI Suggestions:")
			log.Info(suggestions)
		}

		if cfg.Suggest {
			log.Info("")
			log.Info("Nothing was downloaded or installed by suggestion mode.")
		}

		log.Infof("%d tools processed", len(summary.Results))
	} else {
		log.Info("Done!")
	}

	// Handle errors
	if summary.HasErrors() {
		if cfg.Suggest && cfg.ErrorFile.Path() == "" {
			summary.Errors = nonMatchingErrors(summary.Errors)
			summary.Failed = len(summary.Errors)
		}

		if !summary.HasErrors() {
			return
		}

		showErrors(summary, cfg, log)
	}
}

func nonMatchingErrors(details []processor.ErrorDetail) []processor.ErrorDetail {
	filtered := make([]processor.ErrorDetail, 0, len(details))

	for _, detail := range details {
		if _, ok := errors.AsType[*match.SelectionError](detail.Error); ok {
			continue
		}

		filtered = append(filtered, detail)
	}

	return filtered
}

// showErrors formats and displays error messages.
func showErrors(summary processor.Summary, cfg ShowConfig, log *logger.Logger) {
	// Determine error format
	format := ErrorFormatText

	if cfg.ErrorFile.Path() != "" {
		format = ErrorFormatJSON
	}

	const errorWrapWidth = 120

	// Create error formatter
	errorFormatter := NewErrorFormatter(ErrorConfig{
		WrapWidth: errorWrapWidth,
		Format:    format,
	})

	// Format errors
	errorOutput, err := errorFormatter.FormatErrors(summary.Errors)
	if err != nil {
		log.Errorf("failed to format errors: %v", err)

		return
	}

	// Output errors
	if cfg.ErrorFile.Path() == "" {
		log.Error(errorOutput)
	} else {
		if err := cfg.ErrorFile.Write([]byte(errorOutput)); err != nil {
			log.Errorf("failed to write error output to %q: %v", cfg.ErrorFile.Path(), err)
		} else {
			log.Errorf("See error file %q for details", cfg.ErrorFile.Path())
		}
	}
}
