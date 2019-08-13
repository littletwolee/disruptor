package disruptor

import (
	"context"
	"fmt"
	"runtime"
)

type ringBuffer struct {
	buf            jobs
	cs             *cursors
	stop, stopBack chan bool
}

func New(chunk int, l int64) (*ringBuffer, error) {
	if int64(chunk) >= l {
		return nil, fmt.Errorf("2 * chunk <= buffer len")
	}
	return &ringBuffer{
		buf:      make(jobs, l),
		cs:       newCursors(chunk, l),
		stop:     make(chan bool),
		stopBack: make(chan bool),
	}, nil
}
func (rb *ringBuffer) Len() int64 {
	return rb.cs.len()
}
func (rb *ringBuffer) Write(js ...Job) {
	for i := 0; i < len(js); {
		c := rb.cs.get()
		if c.write() {
			rb.buf[c.wCursor()] = js[i]
			c.wCursorCompleted()
			i++
		}
		runtime.Gosched()
	}
}
func (rb *ringBuffer) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	for _, c := range rb.cs.cs {
		go rb.read(ctx, c)
	}
	<-rb.stop
	cancel()
	rb.stopBack <- true
}
func (rb *ringBuffer) read(ctx context.Context, c *_cursor) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if c.read() {
				j := rb.buf[c.rCursor()]
				if j != nil {
					j.CallBack(j.Do())
				}
				rb.buf[c.rCursor()] = nil
				c.rCursorCompleted()
				continue
			}
			runtime.Gosched()

		}
	}
}
func (rb *ringBuffer) Stop() {
	rb.stop <- true
	<-rb.stopBack
}
