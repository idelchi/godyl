package progress

import (
	"fmt"
	"io"
	"os"
	"sync"
	"sync/atomic"
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

// SetManualDisplay customizes a tracker's presentation for sources without
// real byte progress (for example, go install). The supplied labels replace
// the rendered value and speed text while the message is updated to the
// provided string.
func (pt *Tracker) SetManualDisplay(src, message, valueLabel, speedLabel string) {
	pt.lock.Lock()

	tracker, ok := pt.trackers[src]
	pt.lock.Unlock()

	if !ok {
		return
	}

	tracker.UpdateMessage(message)

	var counter uint64

	tracker.Units.Notation = ""
	tracker.Units.NotationPosition = gpp.UnitsNotationPositionBefore
	tracker.Units.Formatter = func(int64) string {
		if speedLabel == "" {
			return valueLabel
		}

		if atomic.AddUint64(&counter, 1)&1 == 1 {
			return valueLabel
		}

		return speedLabel
	}
}

// StartSynthetic creates and drives a synthetic progress entry identified by key.
// It is useful for operations that lack byte-level progress (like go install) but
// should still appear in the global progress list. The bar advances gradually and
// stalls once it reaches the given fraction (0 < stallFraction ≤ 1) until the
// returned function is invoked, at which point it completes instantly.
func StartSynthetic(
	tracker *Tracker,
	key string,
	message string,
	valueLabel string,
	speedLabel string,
	stallFraction float64,
) func() {
	if tracker == nil {
		return func() {}
	}

	pr, pw := io.Pipe()

	const syntheticTotal = 1000

	tracked := tracker.TrackProgress(key, 0, syntheticTotal, pr)
	tracker.SetManualDisplay(key, message, valueLabel, speedLabel)

	wrapped, ok := tracked.(*readCloserWithProgress)
	if !ok {
		return func() {}
	}

	if stallFraction <= 0 || stallFraction > 1 {
		stallFraction = 1
	}

	plateau := int64(float64(syntheticTotal) * stallFraction)
	if plateau <= 0 {
		plateau = syntheticTotal
	}

	const stepSize = 50

	step := syntheticTotal / stepSize
	if step <= 0 {
		step = 1
	}

	stop := make(chan struct{})

	var once sync.Once

	const tickerInterval = 200 * time.Millisecond

	go func() {
		ticker := time.NewTicker(tickerInterval)
		defer ticker.Stop()

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				current := wrapped.Tracker.Value()
				if current >= plateau {
					continue
				}

				next := min(current+int64(step), plateau)

				wrapped.Tracker.SetValue(next)
			}
		}
	}()

	return func() {
		once.Do(func() {
			wrapped.Tracker.SetValue(syntheticTotal)
			wrapped.Tracker.MarkAsDone()
			close(stop)

			_ = tracked.Close()
			_ = pw.Close()
		})
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
