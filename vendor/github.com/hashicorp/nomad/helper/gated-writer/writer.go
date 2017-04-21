package gatedwriter

import (
	"io"
	"sync"
)

// Writer is an io.Writer implementation that buffers all of its
// data into an internal buffer until it is told to let data through.
type Writer struct {
	Writer io.Writer

	buf   [][]byte
	flush bool
	lock  sync.RWMutex
}

// Flush tells the Writer to flush any buffered data and to stop
// buffering.
func (w *Writer) Flush() {
	w.lock.Lock()
	w.flush = true
	w.lock.Unlock()

	for _, p := range w.buf {
		w.Write(p)
	}
	w.buf = nil
}

func (w *Writer) Write(p []byte) (n int, err error) {
	w.lock.RLock()
	defer w.lock.RUnlock()

	if w.flush {
		return w.Writer.Write(p)
	}

	p2 := make([]byte, len(p))
	copy(p2, p)
	w.buf = append(w.buf, p2)
	return len(p), nil
}
