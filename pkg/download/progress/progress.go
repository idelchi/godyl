package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/dustin/go-humanize"
	gpp "github.com/jedib0t/go-pretty/v6/progress"
	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/idelchi/godyl/pkg/path/file"
)

// Trackable defines the interface for progress tracking implementations.
type Trackable interface {
	Start()
	Wait()
	TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser
}

// readCloserWithProgress wraps a reader to update a progress tracker as bytes are read.
type readCloserWithProgress struct {
	io.Reader
	io.Closer

	Tracker   *gpp.Tracker // associated tracker
	totalRead int64        // number of bytes read so far
}

// Read reads from the underlying Reader and updates the tracker with bytes read.
func (r *readCloserWithProgress) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 {
		r.totalRead += int64(n)
		r.Tracker.SetValue(r.totalRead)
	}

	return n, err
}

// Tracker tracks multiple concurrent progress bars and renders them.
type Tracker struct {
	pw       gpp.Writer
	trackers map[string]*gpp.Tracker
	wg       *sync.WaitGroup
	lock     sync.Mutex
}

// New initializes and returns a Tracker with default styles and settings.
func New() *Tracker {
	pw := gpp.NewWriter()

	const MessageLength = 60

	const TrackerLength = 45

	pw.SetMessageLength(MessageLength)
	pw.SetTrackerLength(TrackerLength)
	pw.SetStyle(gpp.StyleDefault)
	pw.SetTrackerPosition(gpp.PositionRight)
	pw.SetUpdateFrequency(100 * time.Millisecond) //nolint:mnd // 100ms update frequency is self-documenting

	// Configure visibility
	pw.Style().Visibility.Percentage = true
	pw.Style().Visibility.Time = false
	pw.Style().Visibility.Value = true
	pw.Style().Visibility.Speed = true
	pw.Style().Visibility.Tracker = true
	pw.Style().Visibility.ETA = false
	pw.Style().Visibility.TrackerOverall = false
	pw.Style().Visibility.SpeedOverall = false
	pw.Style().Visibility.ETAOverall = false
	pw.Style().Visibility.Pinned = false

	// Configure colors
	pw.Style().Colors.Message = text.Colors{text.FgWhite}
	pw.Style().Colors.Tracker = text.Colors{text.FgYellow}
	pw.Style().Colors.Value = text.Colors{text.FgCyan}
	pw.Style().Colors.Time = text.Colors{text.FgGreen}
	pw.Style().Colors.Percent = text.Colors{text.FgHiRed}
	pw.Style().Colors.Speed = text.Colors{text.FgMagenta}

	const TimeInProgressPrecision = 100 * time.Millisecond

	const TimeDonePrecision = 100 * time.Millisecond

	const ETAPrecision = 1 * time.Second

	// Configure timing precision
	pw.Style().Options.TimeInProgressPrecision = TimeInProgressPrecision
	pw.Style().Options.TimeDonePrecision = TimeDonePrecision
	pw.Style().Options.ETAPrecision = ETAPrecision

	pw.SetAutoStop(false)
	pw.SetOutputWriter(os.Stdout)

	return &Tracker{
		pw:       pw,
		trackers: make(map[string]*gpp.Tracker),
		wg:       &sync.WaitGroup{},
	}
}

// Start begins rendering progress bars in a separate goroutine.
func (pt *Tracker) Start() {
	go pt.pw.Render()
}

// Wait waits for all tracked readers to finish and stops rendering once done.
func (pt *Tracker) Wait() {
	const DelayTime = 50 * time.Millisecond

	pt.wg.Wait()

	for pt.pw.IsRenderInProgress() && pt.pw.LengthActive() > 0 {
		time.Sleep(DelayTime)
	}

	pt.pw.Stop()
}

// TrackProgress wraps the given ReadCloser to update the progress bar during read and on close.
func (pt *Tracker) TrackProgress(
	src string,
	currentSize, totalSize int64,
	stream io.ReadCloser,
) io.ReadCloser {
	pt.lock.Lock()
	defer pt.lock.Unlock()

	pt.wg.Add(1)

	srcPretty := file.New(src).Unescape()

	tracker, ok := pt.trackers[src]
	if !ok {
		var sizeStr string

		if totalSize >= 0 {
			sizeStr = humanize.Bytes(uint64(totalSize))
		} else {
			sizeStr = "unknown"
		}

		tracker = &gpp.Tracker{
			Message:            fmt.Sprintf("%-45s [%s]", srcPretty, sizeStr),
			Total:              totalSize,
			Units:              gpp.UnitsBytes,
			DeferStart:         false,
			RemoveOnCompletion: false,
		}

		if currentSize > 0 {
			tracker.SetValue(currentSize)
		}

		pt.trackers[src] = tracker
		pt.pw.AppendTracker(tracker)
	}

	return &readCloserWithProgress{
		Reader:  stream,
		Closer:  &closeWrapper{stream, pt.wg, tracker},
		Tracker: tracker,
	}
}

// closeWrapper wraps an io.Closer to mark a tracker as done and signal completion.
type closeWrapper struct {
	io.Closer

	wg      *sync.WaitGroup
	tracker *gpp.Tracker
}

// Close closes the wrapped resource, marks the tracker as done if needed, and signals via WaitGroup.
func (c *closeWrapper) Close() error {
	err := c.Closer.Close()

	if c.tracker.Value() < c.tracker.Total {
		c.tracker.MarkAsDone()
	}

	c.wg.Done()

	return err
}
