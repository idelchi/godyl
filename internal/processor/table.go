package processor

import (
	"fmt"

	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/path/file"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// ResultTable encapsulates the table creation and rendering logic for displaying tool processing results.
type ResultTable struct {
	// Writer is the underlying table writer instance.
	writer table.Writer

	// RowColors maps row numbers to their display colors.
	rowColors map[int]text.Colors

	// CurrentRowNum tracks the current row being processed.
	currentRowNum int
}

// HeaderConfig defines the configuration for a table header column.
type HeaderConfig struct {
	// Name is the header column name.
	Name any

	// WidthMax is the maximum width for the column.
	WidthMax int

	// Bold indicates if the header should be displayed in bold.
	Bold bool
}

// NewResultTable creates a new result table with properly configured styling.
func NewResultTable(headers ...HeaderConfig) *ResultTable {
	t := table.NewWriter()
	rt := &ResultTable{
		writer:    t,
		rowColors: make(map[int]text.Colors),
	}

	// Extract just the header names for the table
	headerRow := make(table.Row, 0, len(headers))
	columnConfigs := make([]table.ColumnConfig, 0, len(headers))

	for i, h := range headers {
		headerRow = append(headerRow, h.Name)

		// Build column config with the specified width
		config := table.ColumnConfig{
			Number:   i + 1, // Column numbers are 1-based
			WidthMax: h.WidthMax,
		}

		// Apply bold if specified
		if h.Bold {
			config.Colors = text.Colors{text.Bold}
		}

		columnConfigs = append(columnConfigs, config)
	}

	// Add header to the table
	t.AppendHeader(headerRow)

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

	// Apply the column configurations
	t.SetColumnConfigs(columnConfigs)

	return rt
}

// trackRowWithColor assigns colors to the next row.
func (rt *ResultTable) trackRowWithColor(color text.Colors) int {
	rt.currentRowNum++
	rt.rowColors[rt.currentRowNum] = color

	return rt.currentRowNum
}

// AddResult adds a single result to the table with specified coloring.
func (rt *ResultTable) AddResult(tool *tool.Tool, color text.Colors, message string) {
	// Track row and associate colors
	rt.trackRowWithColor(color)

	base := file.File(tool.URL).Unescape().Base()
	if base == "." {
		base = "Not Applicable"
	}

	// Append the row to the table
	rt.writer.AppendRow(table.Row{
		file.New(tool.Exe.Name).WithoutExtension(),
		tool.Version.Version,
		tool.Output,
		fmt.Sprintf("%s/%s", tool.Platform.OS.Name, tool.Platform.Architecture.Name),
		base,
		message,
	})
}

// AddResult adds a single result to the table with specified coloring.
func (rt *ResultTable) AddFail(tool *tool.Tool, color text.Colors, message string) {
	// Track row and associate colors
	rt.trackRowWithColor(color)

	// Append the row to the table
	rt.writer.AppendRow(table.Row{
		tool.Exe.Name,
		message,
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

// Common color providers for different result states.
var (
	ErrorColors   = text.Colors{text.FgRed}
	InfoColors    = text.Colors{text.FgYellow}
	SuccessColors = text.Colors{text.FgGreen}
	WarningColors = text.Colors{text.FgHiYellow}
	DefaultColors = text.Colors{text.BgBlack}
)
