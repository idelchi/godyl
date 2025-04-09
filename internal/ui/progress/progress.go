package progress

import (
	"fmt"
	"io"
	"sync"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	"github.com/idelchi/godyl/internal/tools" // Import tools package
)

// MpbProgressTracker implements getter.ProgressTracker using mpb for visual feedback.
type MpbProgressTracker struct {
	progress *mpb.Progress
	bars     map[string]*mpb.Bar // Map src URL to bar
	mu       sync.Mutex
	tool     *tools.Tool // Store the tool being processed
}

// NewMpbProgressTracker creates a new tracker associated with an mpb.Progress container
// and the specific tool being downloaded.
func NewMpbProgressTracker(p *mpb.Progress, tool *tools.Tool) *MpbProgressTracker {
	return &MpbProgressTracker{
		progress: p,
		bars:     make(map[string]*mpb.Bar),
		tool:     tool, // Store the tool pointer
	}
}

// TrackProgress is called by go-getter to monitor a download stream.
func (t *MpbProgressTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	if stream == nil {
		// Handle cases where the stream might be nil (e.g., getter error before stream starts)
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Format display name using the stored tool info
	name := "download" // Default/fallback name
	if t.tool != nil {
		name = t.tool.Name
		if t.tool.Version.Version != "" {
			name = fmt.Sprintf("%s %s", name, t.tool.Version.Version)
		}
	}

	// Limit name length for display
	const maxNameLen = 35 // Increased length slightly for tool name + version
	if len(name) > maxNameLen {
		name = "..." + name[len(name)-maxNameLen+3:]
	}

	// Create a new bar or update existing one
	bar, ok := t.bars[src] // Use original src as key for mapping
	if !ok || bar == nil {
		// Only set total size if known (>0)
		bar = t.progress.AddBar(totalSize,
			// mpb.BarOptional(mpb.BarRemoveOnComplete(), totalSize <= 0), // Remove if size unknown - let's keep completed bars for now
			// Removed styling attempts for now
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("%-*.*s:", maxNameLen, maxNameLen, name), decor.WC{W: maxNameLen + 1}),
				decor.CountersKibiByte("% .2f / % .2f"), // Display in KiB
			),
			mpb.AppendDecorators(
				decor.NewPercentage("%d%%"),
				decor.Name(" | "),
				decor.AverageETA(decor.ET_STYLE_GO, decor.WC{W: 4}),
				decor.Name(" | "),
				decor.AverageSpeed(decor.SizeB1024(0), "% .1f", decor.WC{W: 7}), // Use SizeB1024(0) for KiB/s etc.
			),
		)
		t.bars[src] = bar
	} else {
		// If bar exists, maybe update total size if it was unknown initially
		if bar.Current() == 0 && totalSize > 0 {
			bar.SetTotal(totalSize, false) // Don't trigger completion yet
		}
	}

	// Set initial progress if resuming
	if currentSize > 0 {
		bar.SetCurrent(currentSize)
	}

	// Create proxy reader
	reader := bar.ProxyReader(stream)

	// Wrap the reader's Close method to potentially remove the bar or mark as complete
	return &progressReader{
		ReadCloser:       reader,
		bar:              bar,
		tracker:          t,
		src:              src,
		initialTotalSize: totalSize, // Store initial total size
	}
}

// progressReader wraps the bar's proxy reader to handle Close.
type progressReader struct {
	io.ReadCloser
	bar              *mpb.Bar
	tracker          *MpbProgressTracker
	src              string
	initialTotalSize int64 // Store the total size known at creation
}

// Close closes the underlying reader and potentially cleans up the bar.
func (r *progressReader) Close() error {
	// Ensure the bar completes if it hasn't already (e.g., due to error before EOF)
	// This might not be strictly necessary if Wait() handles it, but can be safer.
	// Check if the bar is already completed to avoid panic
	if !r.bar.Completed() {
		// Force completion if total size was known initially
		if r.initialTotalSize > 0 {
			r.bar.SetTotal(r.bar.Current(), true) // Mark as complete at current position
		}
		// Abort the bar (true flag attempts removal) regardless of whether total was known.
		// This should prevent the final redraw after completion/EOF.
		r.bar.Abort(true)
	}

	return r.ReadCloser.Close()
}
