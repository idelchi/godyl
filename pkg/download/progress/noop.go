package progress

import (
	"io"
	"sync"
)

type NoopTracker struct{ wg *sync.WaitGroup }

func NewNoop() *NoopTracker { return &NoopTracker{wg: &sync.WaitGroup{}} }

func (n *NoopTracker) Start() {}

func (n *NoopTracker) Wait() { n.wg.Wait() }

func (n *NoopTracker) TrackProgress(_ string, _,
	_ int64, rc io.ReadCloser,
) io.ReadCloser {
	n.wg.Add(1)
	return &noopRC{ReadCloser: rc, wg: n.wg}
}

type noopRC struct {
	io.ReadCloser
	wg *sync.WaitGroup
}

func (n *noopRC) Close() error {
	err := n.ReadCloser.Close()
	n.wg.Done()
	return err
}
