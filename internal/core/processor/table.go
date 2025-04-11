package processor

import (
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ColorProvider is a function that returns text colors for row styling
type ColorProvider func() text.Colors

// ResultTable encapsulates the table creation and rendering logic.
type ResultTable struct {
	writer        table.Writer
	rowColors     map[int]text.Colors
	currentRowNum int
}

// NewResultTable creates a new result table with properly configured styling.
func NewResultTable(headers ...any) *ResultTable {
	t := table.NewWriter()
	rt := &ResultTable{
		writer:    t,
		rowColors: make(map[int]text.Colors),
	}

	// Add header to the table
	t.AppendHeader(table.Row(headers))

	// Set up row painter with custom coloring
	t.SetRowPainter(func(row table.Row, attr table.RowAttributes) text.Colors {
		if colors, exists := rt.rowColors[attr.Number]; exists {
			return colors
		}
		return nil // default color
	})

	// Set style and colors
	t.SetStyle(table.StyleRounded)
	t.Style().Color.Header = text.Colors{text.FgBlue, text.Bold}

	// Configure column widths and formatting
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMax: 20},                                 // Tool name column
		{Number: 2, WidthMax: 15},                                 // Version column
		{Number: 3, WidthMax: 40},                                 // Output Path column
		{Number: 4, WidthMax: 30},                                 // Aliases column
		{Number: 5, WidthMax: 50, Colors: text.Colors{text.Bold}}, // Status column always bold
	})

	return rt
}

// trackRowWithColor assigns colors to the next row
func (rt *ResultTable) trackRowWithColor(colorProvider ColorProvider) int {
	rt.currentRowNum++
	if colorProvider != nil {
		rt.rowColors[rt.currentRowNum] = colorProvider()
	}
	return rt.currentRowNum
}

// AddResult adds a single result to the table with specified coloring.
func (rt *ResultTable) AddResult(r result, status string, colorProvider ColorProvider) {
	tool := r.Tool

	// Format aliases as a comma-separated string if present
	aliases := ""
	if tool.Aliases != nil && len(tool.Aliases) > 0 {
		aliases = strings.Join(tool.Aliases, ", ")
	}

	// Track row and associate colors
	rt.trackRowWithColor(colorProvider)

	// Append the row to the table
	rt.writer.AppendRow(table.Row{
		tool.Exe.Name,
		tool.Version.Version,
		tool.Output,
		aliases,
		status,
	})
}

// AddCustomRow adds a custom row to the table with specified coloring.
func (rt *ResultTable) AddCustomRow(toolName, version, output, aliases, status string, colorProvider ColorProvider) {
	// Track row and associate colors
	rt.trackRowWithColor(colorProvider)

	rt.writer.AppendRow(table.Row{
		toolName,
		version,
		output,
		aliases,
		status,
	})
}

// Clear removes all rows from the table (except the header).
func (rt *ResultTable) Clear() {
	rt.writer.ResetRows()
	rt.rowColors = make(map[int]text.Colors)
	rt.currentRowNum = 0
}

// Render returns the rendered table as a string.
func (rt *ResultTable) Render() string {
	return rt.writer.Render()
}

// Common color providers for convenience
var (
	ErrorColors   = func() text.Colors { return text.Colors{text.FgRed} }
	InfoColors    = func() text.Colors { return text.Colors{text.FgYellow} }
	SuccessColors = func() text.Colors { return text.Colors{text.FgGreen} }
	WarningColors = func() text.Colors { return text.Colors{text.FgHiYellow} }
	DefaultColors = func() text.Colors { return nil }
)
