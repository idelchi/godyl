// Package presentation handles all UI and formatting logic.
package presentation

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/idelchi/godyl/internal/runner"
	"github.com/idelchi/godyl/pkg/path/file"
)

// TableFormatter handles table rendering for results.
type TableFormatter struct {
	config TableConfig
}

// TableConfig configures the table formatter.
type TableConfig struct {
	MaxWidth int
	Verbose  bool
}

// NewTableFormatter creates a new table formatter.
func NewTableFormatter(config TableConfig) *TableFormatter {
	return &TableFormatter{
		config: config,
	}
}

// RenderResults renders a collection of results as a formatted table.
func (f *TableFormatter) RenderResults(results []runner.Result) string {
	if len(results) == 0 {
		return ""
	}

	t := table.NewWriter()

	// Configure headers
	headers := []HeaderConfig{
		{Name: "Tool", WidthMax: f.config.MaxWidth},
		{Name: "Version", WidthMax: f.config.MaxWidth},
		{Name: "Output Path", WidthMax: f.config.MaxWidth},
		{Name: "OS/ARCH", WidthMax: f.config.MaxWidth},
		{Name: "File", WidthMax: f.config.MaxWidth},
		{Name: "Status", WidthMax: f.config.MaxWidth, Bold: true},
	}

	// Set up headers
	headerRow := make(table.Row, 0, len(headers))
	columnConfigs := make([]table.ColumnConfig, 0, len(headers))

	for i, h := range headers {
		headerRow = append(headerRow, h.Name)
		config := table.ColumnConfig{
			Number:   i + 1,
			WidthMax: h.WidthMax,
		}

		if h.Bold {
			config.Colors = text.Colors{text.Bold}
		}

		columnConfigs = append(columnConfigs, config)
	}

	t.AppendHeader(headerRow)
	t.SetColumnConfigs(columnConfigs)

	// Set style
	t.SetStyle(table.StyleRounded)

	t.Style().Color.Header = text.Colors{text.FgBlue, text.Bold}

	t.SortBy([]table.SortBy{{Name: "Tool", Mode: table.Asc}})

	// Track row colors
	rowColors := make(map[int]text.Colors)
	rowNum := 0

	// Set up row painter
	t.SetRowPainter(func(row table.Row, attr table.RowAttributes) text.Colors {
		if colors, exists := rowColors[attr.Number]; exists {
			return colors
		}

		return nil
	})

	// Add results
	for _, result := range results {
		rowNum++
		color := f.getColorForStatus(result.Status)
		rowColors[rowNum] = color

		row := f.formatResultRow(result)
		t.AppendRow(row)
	}

	return t.Render()
}

// formatResultRow formats a single result into a table row.
func (f *TableFormatter) formatResultRow(result runner.Result) table.Row {
	tool := result.Tool

	// Format file/URL
	fileDisplay := "Not Applicable"
	if tool.URL != "" {
		fileDisplay = file.File(tool.URL).Unescape().Base()
	}

	// Format executable name
	exeName := file.New(tool.Exe.Name).WithoutExtension().String()
	if tool.Mode == "extract" && fileDisplay != "Not Applicable" {
		exeName = fileDisplay
	}

	// Format message
	message := result.Message
	if result.Status == runner.StatusFailed && f.config.Verbose {
		message = "failed, see below for details"
	}

	return table.Row{
		exeName,
		tool.Version.Version,
		tool.Output,
		fmt.Sprintf("%s/%s", tool.Platform.OS.Name, tool.Platform.Architecture.Name),
		fileDisplay,
		message,
	}
}

// getColorForStatus returns the appropriate color for a given status.
func (f *TableFormatter) getColorForStatus(status runner.Status) text.Colors {
	switch status {
	case runner.StatusOK:
		return text.Colors{text.FgGreen}
	case runner.StatusFailed:
		return text.Colors{text.FgRed}
	case runner.StatusSkipped:
		return text.Colors{text.FgYellow}
	default:
		return text.Colors{text.BgBlack}
	}
}

// HeaderConfig defines the configuration for a table header column.
type HeaderConfig struct {
	Name     any
	WidthMax int
	Bold     bool
}
