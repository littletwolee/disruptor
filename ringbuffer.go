package disruptor

import (
	"context"
	"fmt"
	"runtime"
)

type ringBuffer struct {
	chunk          int64
	buf            jobs
	cs             *cursors
	stop, stopBack chan bool
	wCheck         *wCompleted
	wait           chan bool
}

func New(chunk int, l int64) (*ringBuffer, error) {
	if int64(chunk) >= l {
		return nil, fmt.Errorf("2 * chunk <= buffer len")
	}
	return &ringBuffer{
		chunk:    int64(chunk),
		buf:      make(jobs, l),
		cs:       newCursors(chunk, l),
		stop:     make(chan bool),
		stopBack: make(chan bool),
		wCheck:   &_EMPTY,
		wait:     make(chan bool),
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
			continue
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
			// fmt.Printf("isc:%v,len:%d\n", rb.wCheck.isCompleted(), c.len())
			if rb.wCheck.isCompleted() && rb.Len() == 0 {
				// fmt.Println(rb.Len())
				rb.wCheck.empty()
				rb.wait <- true
			}
			if c.read() {
				j := rb.buf[c.rNext()]
				if j != nil {
					j.CallBack(j.Do())
					rb.buf[c.rNext()] = nil
					c.rCursorCompleted()
				}
			}
			runtime.Gosched()
		}
	}
}
func (rb *ringBuffer) Wait() {
	<-rb.wait
}

func (rb *ringBuffer) Stop() {
	rb.stop <- true
	<-rb.stopBack
}
func (rb *ringBuffer) Writing() {
	rb.wCheck.writing()
}
func (rb *ringBuffer) Completed() {
	rb.wCheck.completed()
}
