package progress

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/idelchi/godyl/internal/tools" // Import tools package
)

// formatBytes converts bytes to a human-readable string (KB, MB, GB, etc.)
func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// Global shared progress writer
var (
	sharedWriter progress.Writer
	writerOnce   sync.Once
	activeCount  int
	writerMu     sync.Mutex
)

// initSharedWriter initializes the shared progress writer
func initSharedWriter() progress.Writer {
	writerOnce.Do(func() {
		// Create a new progress writer
		pw := progress.NewWriter()

		// Configure the progress writer
		pw.SetAutoStop(true)
		pw.SetMessageLength(35)
		pw.SetTrackerLength(40)
		pw.SetStyle(progress.StyleDefault)
		pw.SetTrackerPosition(progress.PositionRight)
		pw.SetUpdateFrequency(time.Millisecond * 100)

		// Configure visibility options
		pw.Style().Visibility.Percentage = true
		pw.Style().Visibility.Time = true
		pw.Style().Visibility.Tracker = true
		pw.Style().Visibility.Value = true
		pw.Style().Visibility.ETA = true
		pw.Style().Visibility.Speed = true

		// Note: Total bytes will be shown in the tracker message

		// Configure colors
		pw.Style().Colors.Message = text.Colors{text.FgWhite}
		pw.Style().Colors.Tracker = text.Colors{text.FgYellow}
		pw.Style().Colors.Value = text.Colors{text.FgCyan}
		pw.Style().Colors.Time = text.Colors{text.FgGreen}
		pw.Style().Colors.Percent = text.Colors{text.FgHiRed}
		pw.Style().Colors.Speed = text.Colors{text.FgMagenta}

		// Start rendering in a goroutine
		go pw.Render()

		sharedWriter = pw
	})

	return sharedWriter
}

// PrettyProgressTracker implements getter.ProgressTracker using go-pretty for visual feedback.
type PrettyProgressTracker struct {
	trackers map[string]*progress.Tracker // Map src URL to tracker
	mu       sync.Mutex
	tool     *tools.Tool // Store the tool being processed
}

// NewPrettyProgressTracker creates a new tracker associated with a go-pretty Progress container
// and the specific tool being downloaded.
func NewPrettyProgressTracker(tool *tools.Tool) *PrettyProgressTracker {
	// Initialize the shared writer if needed
	pw := initSharedWriter()

	// Increment active count
	writerMu.Lock()
	activeCount++
	pw.SetNumTrackersExpected(activeCount)
	writerMu.Unlock()

	return &PrettyProgressTracker{
		trackers: make(map[string]*progress.Tracker),
		tool:     tool, // Store the tool pointer
	}
}

// TrackProgress is called by go-getter to monitor a download stream.
func (t *PrettyProgressTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	if stream == nil {
		// Handle cases where the stream might be nil (e.g., getter error before stream starts)
		return nil
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	// Format display name using the stored tool info
	name := t.tool.Exe.Name
	if t.tool.Version.Version != "" {
		name = fmt.Sprintf("%s %s", name, t.tool.Version.Version)
	}

	// Limit name length for display
	const maxNameLen = 35 // Increased length slightly for tool name + version
	if len(name) > maxNameLen {
		name = "..." + name[len(name)-maxNameLen+3:]
	}

	// Create a new tracker or update existing one
	tracker, ok := t.trackers[src] // Use original src as key for mapping
	if !ok || tracker == nil {
		// Format message with total size if available
		var message string
		if totalSize > 0 {
			// Format total size in human-readable format
			totalSizeStr := formatBytes(totalSize)
			message = fmt.Sprintf("%-*.*s [%s]", maxNameLen-len(totalSizeStr)-3, maxNameLen-len(totalSizeStr)-3, name, totalSizeStr)
		} else {
			message = fmt.Sprintf("%-*.*s", maxNameLen, maxNameLen, name)
		}

		// Create a new tracker
		tracker = &progress.Tracker{
			Message: message,
			Total:   totalSize,
			Units:   progress.UnitsBytes, // Use bytes units for file downloads
		}

		// If total size is unknown, mark as indeterminate
		if totalSize <= 0 {
			tracker.Total = 0 // This makes it indeterminate in go-pretty
		}

		// Start the tracker
		tracker.Start()

		// Add to the shared progress writer
		initSharedWriter().AppendTracker(tracker)

		// Store in our map
		t.trackers[src] = tracker
	} else {
		// If tracker exists, maybe update total size if it was unknown initially
		if tracker.Value() == 0 && totalSize > 0 {
			tracker.UpdateTotal(totalSize)
		}
	}

	// Set initial progress if resuming
	if currentSize > 0 {
		tracker.SetValue(currentSize)
	}

	// Create a reader that updates the progress
	return &prettyProgressReader{
		ReadCloser:       stream,
		tracker:          tracker,
		trackerContainer: t,
		src:              src,
		initialTotalSize: totalSize, // Store initial total size
	}
}

// prettyProgressReader wraps a ReadCloser to update progress as data is read.
type prettyProgressReader struct {
	io.ReadCloser
	tracker          *progress.Tracker
	trackerContainer *PrettyProgressTracker
	src              string
	initialTotalSize int64 // Store the total size known at creation
}

// Read reads data from the underlying reader and updates the progress.
func (r *prettyProgressReader) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if n > 0 {
		r.tracker.Increment(int64(n))
	}
	return n, err
}

// Close closes the underlying reader and marks the tracker as done.
func (r *prettyProgressReader) Close() error {
	// Mark the tracker as done
	if !r.tracker.IsDone() {
		if r.initialTotalSize > 0 {
			// If we know the total size, ensure we're at 100%
			if r.tracker.Value() < r.initialTotalSize {
				r.tracker.SetValue(r.initialTotalSize)
			}
		} else {
			// If we don't know the total size, set the total to the current value
			r.tracker.UpdateTotal(r.tracker.Value())
		}

		// Mark as done
		r.tracker.MarkAsDone()
	}

	return r.ReadCloser.Close()
}

// Wait decrements the active count.
// The shared writer will auto-stop when all trackers are done.
func (t *PrettyProgressTracker) Wait() {
	// Decrement active count
	writerMu.Lock()
	activeCount--
	writerMu.Unlock()

	// No need to stop the writer manually, it will auto-stop when all trackers are done
}

// StopSharedWriter explicitly stops the shared writer.
// This should only be used in exceptional cases, as the writer will auto-stop
// when all trackers are done.
func StopSharedWriter() {
	writerMu.Lock()
	defer writerMu.Unlock()

	if sharedWriter != nil {
		sharedWriter.Stop()
		activeCount = 0
	}
}
