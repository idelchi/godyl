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
)

/* ------------------------------------------------------------------ */
/*                          Public interface                          */
/* ------------------------------------------------------------------ */

type ProgressReporter interface {
	Start()
	Wait()
	TrackProgress(src string, currentSize, totalSize int64, stream io.ReadCloser) io.ReadCloser
}

/* ------------------------------------------------------------------ */
/*                         Real implementation                        */
/* ------------------------------------------------------------------ */

type readCloserWithProgress struct {
	io.Reader
	io.Closer
	Tracker   *gpp.Tracker
	totalRead int64
}

func (r *readCloserWithProgress) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if n > 0 {
		r.totalRead += int64(n)
		r.Tracker.SetValue(r.totalRead)
	}

	return n, err
}

/* ------------------------------------------------------------------ */
/*                         ProgressTracker                            */
/* ------------------------------------------------------------------ */

type ProgressTracker struct {
	pw       gpp.Writer
	trackers map[string]*gpp.Tracker
	lock     sync.Mutex
	wg       *sync.WaitGroup
}

// NewProgressTracker creates a tracker; pass showETA=true to display ETA.
func New() *ProgressTracker {
	pw := gpp.NewWriter()

	// --- Styling (aligned with external reference) ------------------
	pw.SetMessageLength(45)
	pw.SetTrackerLength(60)
	pw.SetStyle(gpp.StyleDefault)
	pw.SetTrackerPosition(gpp.PositionRight)
	pw.SetUpdateFrequency(100 * time.Millisecond)

	// Visibility
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

	// Colours
	pw.Style().Colors.Message = text.Colors{text.FgWhite}
	pw.Style().Colors.Tracker = text.Colors{text.FgYellow}
	pw.Style().Colors.Value = text.Colors{text.FgCyan}
	pw.Style().Colors.Time = text.Colors{text.FgGreen}
	pw.Style().Colors.Percent = text.Colors{text.FgHiRed}
	pw.Style().Colors.Speed = text.Colors{text.FgMagenta}

	// Limit time precision to two decimals (requires go-pretty ≥ 6.4)
	pw.Style().Options.TimeInProgressPrecision = 100 * time.Millisecond
	pw.Style().Options.TimeDonePrecision = 100 * time.Millisecond
	pw.Style().Options.ETAPrecision = 1 * time.Second
	//-----------------------------------------------------------------

	pw.SetAutoStop(false)
	pw.SetOutputWriter(os.Stdout)

	return &ProgressTracker{
		pw:       pw,
		trackers: map[string]*gpp.Tracker{},
		wg:       &sync.WaitGroup{},
	}
}

func (pt *ProgressTracker) Start() { go pt.pw.Render() }

func (pt *ProgressTracker) Wait() {
	pt.wg.Wait()
	for pt.pw.IsRenderInProgress() {
		if pt.pw.LengthActive() == 0 {
			pt.pw.Stop()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (pt *ProgressTracker) TrackProgress(
	src string,
	currentSize, totalSize int64,
	stream io.ReadCloser,
) io.ReadCloser {
	pt.lock.Lock()
	defer pt.lock.Unlock()

	pt.wg.Add(1)

	tr, ok := pt.trackers[src]
	if !ok {
		tr = &gpp.Tracker{
			Message:            fmt.Sprintf("%-50s [%s]", src, humanize.Bytes(uint64(totalSize))),
			Total:              totalSize,
			Units:              gpp.UnitsBytes,
			DeferStart:         false,
			RemoveOnCompletion: false,
		}
		if currentSize > 0 {
			tr.SetValue(currentSize)
		}
		pt.trackers[src] = tr
		pt.pw.AppendTracker(tr)
	}

	return &readCloserWithProgress{
		Reader:  stream,
		Closer:  &closeWrapper{stream, pt.wg, tr},
		Tracker: tr,
	}
}

/* ------------------------------------------------------------------ */
/*                           helpers                                  */
/* ------------------------------------------------------------------ */

type closeWrapper struct {
	io.Closer
	wg      *sync.WaitGroup
	tracker *gpp.Tracker
}

func (c *closeWrapper) Close() error {
	err := c.Closer.Close()
	if c.tracker.Value() < c.tracker.Total {
		c.tracker.MarkAsDone()
	}
	c.wg.Done()
	return err
}
