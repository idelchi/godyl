// Package processor handles the processing of tool installations and management.
package processor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	"github.com/jedib0t/go-pretty/v6/text"

	"github.com/idelchi/godyl/internal/cache"
	"github.com/idelchi/godyl/internal/config"
	"github.com/idelchi/godyl/internal/tools"
	"github.com/idelchi/godyl/internal/tools/tags"
	"github.com/idelchi/godyl/internal/tools/tool"
	"github.com/idelchi/godyl/pkg/download/progress"
	"github.com/idelchi/godyl/pkg/logger"
	"github.com/idelchi/godyl/pkg/path/file"
)

// ErrToolsFailedToInstall is returned when one or more tools failed to install.
var ErrToolsFailedToInstall = errors.New("tools failed to install")

// Processor handles tool installation and management.
type Processor struct {
	log         *logger.Logger
	cache       *cache.Cache
	progressMgr *progress.Trackable
	tools       tools.Tools
	results     tools.Tools
	Options     []tool.ResolveOption
	config      config.Config
	mu          sync.Mutex
	hasErrors   bool
	NoDownload  bool
}

// New creates a new Processor.
func New(toolsList tools.Tools, cfg config.Config, log *logger.Logger) *Processor {
	return &Processor{
		tools:   toolsList,
		config:  cfg,
		log:     log,
		results: make(tools.Tools, 0, len(toolsList)),
	}
}

// Process installs and manages tools with the given tags.
func (p *Processor) Process(tags tags.IncludeTags) error {
	if !p.config.Root.Cache.Disabled {
		if err := p.Cache(); err != nil {
			return err
		}
	}

	if err := p.processTools(tags); err != nil {
		return err
	}

	p.logFinalResults()

	if p.hasErrors {
		return fmt.Errorf("one or more %w", ErrToolsFailedToInstall)
	}

	return nil
}

// Cache initializes the cache for the processor.
func (p *Processor) Cache() error {
	file, _ := cache.File(p.config.Root.Cache.Dir)
	p.cache = cache.New(file)

	if !p.config.Root.Cache.Disabled {
		return p.cache.Load()
	}

	return nil
}

// processTools handles the concurrent processing of all tools.
func (p *Processor) processTools(tags tags.IncludeTags) error {
	results := make(chan *tool.Tool, len(p.tools)) // buffered: no blocking, no mutex

	// select the progress reporter once, upfront
	var tracker progress.Trackable
	if p.config.Root.NoProgress {
		tracker = progress.NewNoop()
	} else {
		tracker = progress.New()
		tracker.Start()
	}

	g, _ := errgroup.WithContext(context.Background())
	if par := p.config.Root.Parallel; par > 0 {
		g.SetLimit(par)
		p.log.Debugf("running with %d parallel downloads", par)
	}

	// spawn a worker per tool
	for _, t := range p.tools {
		g.Go(func() error {
			p.processOneTool(t, tags, results, tracker)

			return nil
		})
	}

	// close the results channel once all workers are done
	go func() {
		_ = g.Wait()

		close(results)
	}()

	// collect results (single goroutine ⇒ no mutex)
	for r := range results {
		p.results = append(p.results, r)
		if r.Result().AsError() != nil {
			p.hasErrors = true
		}
	}

	// propagate worker error (if any)
	if err := g.Wait(); err != nil {
		return err
	}

	tracker.Wait() // wait for progress bars to finish

	return nil
}

// logFinalResults iterates over collected results and logs them as a table.
func (p *Processor) logFinalResults() {
	const maxWidth = 100
	// Initialize result table
	table := NewResultTable(
		HeaderConfig{Name: "Tool", WidthMax: maxWidth},
		HeaderConfig{Name: "Version", WidthMax: maxWidth},
		HeaderConfig{Name: "Output Path", WidthMax: maxWidth},
		HeaderConfig{Name: "OS/ARCH", WidthMax: maxWidth},
		HeaderConfig{Name: "File", WidthMax: maxWidth},
		HeaderConfig{Name: "Status", WidthMax: maxWidth, Bold: true},
	)

	// Add all result rows to the table
	for _, r := range p.results {
		message := r.Result().Message

		// Determine the appropriate color provider based on the status
		var color text.Colors

		switch {
		case r.Result().IsFailed():
			color = ErrorColors
			message = "failed, see below for details"
		case r.Result().IsSkipped():
			color = InfoColors

			if r.URL == "" {
				r.URL = "Not Applicable"
			} else {
				r.URL = file.New(r.URL).Base()
			}
		case r.Result().IsOK():
			p.UpdateCache(r)

			color = SuccessColors
		default:
			color = DefaultColors
		}

		table.AddResult(r, color, message)
	}

	if table.writer.Length() == 0 {
		p.log.Info("Nothing of interest to show")

		return
	}

	// Render the table
	if p.config.Root.Verbose > 0 {
		p.log.Info("")
		p.log.Info("Installation Summary:")
		p.log.Info(table.Render())
	} else {
		p.log.Info("Done!")
	}

	const wrap = 120

	messages := errorMessages{}

	for _, r := range p.results {
		if r.Result().IsFailed() {
			messages = append(messages, errorMessage{
				Tool:    r.Name,
				Message: r.Result().Error(),
			})
		}
	}

	errOutput := p.config.Root.ErrorFile

	if errOutput.Path() == "" {
		p.log.Error(messages.Dump())
	} else {
		bytes, err := messages.ToJsonBytes()
		if err != nil {
			p.log.Errorf("failed to marshal error messages: %v", err)
		}

		if err := errOutput.Write(bytes); err != nil {
			p.log.Errorf("failed to write error output to %q: %v", errOutput.Path(), err)
		}

		p.log.Errorf("See error file %q for details", errOutput.Path())
	}
}

type errorMessage struct {
	Tool    string
	Message string
}

type errorMessages []errorMessage

func (e errorMessages) ToJsonBytes() ([]byte, error) {
	return json.MarshalIndent(e, "  ", "  ")
}

func (e errorMessages) Dump() string {
	var sb strings.Builder

	for _, m := range e {
		sb.WriteString("\n")
		sb.WriteString(fmt.Sprintf("%q", m.Tool))
		sb.WriteString("\n")
		sb.WriteString(m.Message)
	}

	return sb.String()
}
