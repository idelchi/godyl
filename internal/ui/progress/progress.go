package progress

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/idelchi/godyl/internal/tools"
)

const (
	defaultMaxNameLen = 35
)

// Options defines configuration for progress tracking
type Options struct {
	MessageLength   int
	TrackerLength   int
	UpdateFrequency time.Duration
	Style           progress.Style
	TrackerPosition progress.Position
	ShowPercentage  bool
	ShowTime        bool
	ShowTracker     bool
	ShowValue       bool
	ShowETA         bool
	ShowSpeed       bool
	MessageColor    text.Colors
	TrackerColor    text.Colors
	ValueColor      text.Colors
	TimeColor       text.Colors
	PercentColor    text.Colors
	SpeedColor      text.Colors
}

// DefaultOptions provides sensible default configuration
func DefaultOptions() Options {
	return Options{
		MessageLength:   35,
		TrackerLength:   40,
		UpdateFrequency: 100 * time.Millisecond,
		Style:           progress.StyleDefault,
		TrackerPosition: progress.PositionRight,
		ShowPercentage:  true,
		ShowTime:        true,
		ShowTracker:     true,
		ShowValue:       true,
		ShowETA:         true,
		ShowSpeed:       true,
		MessageColor:    text.Colors{text.FgWhite},
		TrackerColor:    text.Colors{text.FgYellow},
		ValueColor:      text.Colors{text.FgCyan},
		TimeColor:       text.Colors{text.FgGreen},
		PercentColor:    text.Colors{text.FgHiRed},
		SpeedColor:      text.Colors{text.FgMagenta},
	}
}

// ProgressManager handles progress tracking with configurable options
type ProgressManager struct {
	writer progress.Writer
	mu     sync.Mutex
	count  int
	opts   Options
}

// NewProgressManager creates a new manager with the given options
func NewProgressManager(opts Options) *ProgressManager {
	pw := progress.NewWriter()

	// Apply writer options
	pw.SetAutoStop(true)
	pw.SetMessageLength(opts.MessageLength)
	pw.SetTrackerLength(opts.TrackerLength)
	pw.SetStyle(opts.Style)
	pw.SetTrackerPosition(opts.TrackerPosition)
	pw.SetUpdateFrequency(opts.UpdateFrequency)

	// Configure visibility options
	pw.Style().Visibility.Percentage = opts.ShowPercentage
	pw.Style().Visibility.Time = opts.ShowTime
	pw.Style().Visibility.Tracker = opts.ShowTracker
	pw.Style().Visibility.Value = opts.ShowValue
	pw.Style().Visibility.ETA = opts.ShowETA
	pw.Style().Visibility.Speed = opts.ShowSpeed

	// Configure colors
	pw.Style().Colors.Message = opts.MessageColor
	pw.Style().Colors.Tracker = opts.TrackerColor
	pw.Style().Colors.Value = opts.ValueColor
	pw.Style().Colors.Time = opts.TimeColor
	pw.Style().Colors.Percent = opts.PercentColor
	pw.Style().Colors.Speed = opts.SpeedColor

	// Start rendering in a goroutine
	go pw.Render()

	return &ProgressManager{
		writer: pw,
		opts:   opts,
	}
}

// NewTracker creates a progress tracker for a specific tool
func (pm *ProgressManager) NewTracker(tool *tools.Tool) *PrettyProgressTracker {
	if tool == nil {
		panic("cannot create progress tracker with nil tool")
	}

	pm.mu.Lock()
	pm.count++
	pm.writer.SetNumTrackersExpected(pm.count)
	pm.mu.Unlock()

	return &PrettyProgressTracker{
		trackers: make(map[string]*progress.Tracker),
		tool:     tool,
		manager:  pm,
	}
}

// DecrementCount decreases the active tracker count
func (pm *ProgressManager) DecrementCount() {
	pm.mu.Lock()
	if pm.count > 0 {
		pm.count--
	}
	pm.mu.Unlock()
}

// Stop stops the progress writer explicitly
func (pm *ProgressManager) Stop() {
	pm.mu.Lock()
	pm.writer.Stop()
	pm.count = 0
	pm.mu.Unlock()
}

// PrettyProgressTracker implements getter.ProgressTracker using go-pretty for visual feedback
type PrettyProgressTracker struct {
	trackers map[string]*progress.Tracker
	tool     *tools.Tool
	manager  *ProgressManager
	mu       sync.Mutex
}

// TrackProgress is called by go-getter to monitor a download stream
func (t *PrettyProgressTracker) TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser {
	if stream == nil {
		return nil
	}

	// Only lock while accessing/modifying the trackers map
	t.mu.Lock()
	tracker, exists := t.trackers[src]
	if !exists {
		tracker = t.createNewTracker(src, totalSize)
	} else if tracker.Value() == 0 && totalSize > 0 {
		tracker.UpdateTotal(totalSize)
	}
	t.mu.Unlock()

	// Set initial progress if resuming (doesn't need lock)
	if currentSize > 0 {
		tracker.SetValue(currentSize)
	}

	// Create a reader that updates the progress
	return &prettyProgressReader{
		ReadCloser:       stream,
		tracker:          tracker,
		initialTotalSize: totalSize,
	}
}

// createNewTracker creates a new progress tracker for a given source
// Caller must hold the mutex
func (t *PrettyProgressTracker) createNewTracker(src string, totalSize int64) *progress.Tracker {
	// Format tool information into a message
	message := t.formatTrackerMessage(totalSize)

	// Create a new tracker
	tracker := &progress.Tracker{
		Message: message,
		Total:   totalSize,
		Units:   progress.UnitsBytes,
	}

	// If total size is unknown, mark as indeterminate
	if totalSize <= 0 {
		tracker.Total = 0
	}

	// Start the tracker and add to the writer
	tracker.Start()
	t.manager.writer.AppendTracker(tracker)

	// Store in our map
	t.trackers[src] = tracker

	return tracker
}

// formatTrackerMessage creates a formatted message for the progress tracker
func (t *PrettyProgressTracker) formatTrackerMessage(totalSize int64) string {
	// Get base name
	name := t.tool.Name
	if t.tool.Exe.Name != "" {
		name = t.tool.Exe.Name
	}
	if t.tool.Version.Version != "" {
		name = fmt.Sprintf("%s %s", name, t.tool.Version.Version)
	}

	// Truncate name if too long
	maxNameLen := defaultMaxNameLen
	if len(name) > maxNameLen {
		name = "..." + name[len(name)-maxNameLen+3:]
	}

	// Format with size information if available
	if totalSize > 0 {
		totalSizeStr := humanize.Bytes(uint64(totalSize))
		availableSpace := maxNameLen - len(totalSizeStr) - 3
		return fmt.Sprintf("%-*.*s [%s]", availableSpace, availableSpace, name, totalSizeStr)
	}

	return fmt.Sprintf("%-*.*s", maxNameLen, maxNameLen, name)
}

// Wait decrements the active count
func (t *PrettyProgressTracker) Wait() {
	t.manager.DecrementCount()
}

// Clear removes all trackers from this progress tracker
func (t *PrettyProgressTracker) Clear() {
	t.mu.Lock()
	t.trackers = make(map[string]*progress.Tracker)
	t.mu.Unlock()
}

// prettyProgressReader wraps a ReadCloser to update progress as data is read
type prettyProgressReader struct {
	io.ReadCloser
	tracker          *progress.Tracker
	initialTotalSize int64
}

// Read reads data from the underlying reader and updates the progress
func (r *prettyProgressReader) Read(p []byte) (int, error) {
	n, err := r.ReadCloser.Read(p)
	if n > 0 {
		r.tracker.Increment(int64(n))
	}
	return n, err
}

// Close closes the underlying reader and marks the tracker as done
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
