package disruptor

import (
	"sync/atomic"
)

type _cursors []*_cursor

func (_cs *_cursors) len() int64 {
	var l int64
	for _, c := range *_cs {
		l += c.len()
	}
	return l
}

type cursors struct {
	index *int64
	cs    _cursors
	chunk int64
}

func newCursors(chunk int, bufLen int64) *cursors {
	var _cs _cursors
	chunkInt64 := int64(chunk)
	for i := 0; i < chunk; i++ {
		s := int64(i) * bufLen / chunkInt64
		e := int64(i+1) * bufLen / chunkInt64
		r, w := s, s+1
		_cs = append(_cs, &_cursor{r: &r, w: &w, s: s, e: e})
	}
	return &cursors{
		index: new(int64),
		cs:    _cs,
		chunk: int64(chunk),
	}
}
func (cs *cursors) get() *_cursor {
	i := atomic.AddInt64(cs.index, 1)
	return cs.cs[i%cs.chunk]
}
func (cs *cursors) len() int64 {
	return cs.cs.len()
}

type _cursor struct {
	r, w *int64
	s, e int64
}

func (c *_cursor) rCursor() int64 {
	return atomic.LoadInt64(c.r)
}
func (c *_cursor) rCursorCompleted() {
	if !atomic.CompareAndSwapInt64(c.r, c.e-1, c.s) {
		atomic.AddInt64(c.r, 1)
	}
}
func (c *_cursor) wCursor() int64 {
	return atomic.LoadInt64(c.w)
}
func (c *_cursor) wCursorCompleted() {
	if !atomic.CompareAndSwapInt64(c.w, c.e-1, c.s) {
		atomic.AddInt64(c.w, 1)
	}
}
func (c *_cursor) read() bool {
	r, w := atomic.LoadInt64(c.r), atomic.LoadInt64(c.w)
	return (r <= w && w < c.e) || (w < r && r < c.e)
}
func (c *_cursor) write() bool {
	r, w := atomic.LoadInt64(c.r), atomic.LoadInt64(c.w)
	return (r < w && w < c.e) || (w < r && r < c.e)
}
func (c *_cursor) len() int64 {
	r, w := atomic.LoadInt64(c.r), atomic.LoadInt64(c.w)
	if r <= w {
		return w - r
	}
	return c.e - r + w - c.s
}
