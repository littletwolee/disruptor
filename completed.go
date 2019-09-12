package disruptor

import (
	"sync/atomic"
	"unsafe"
)

type wCompleted int32

var (
	_EMPTY     wCompleted = 0
	_COMPLETED wCompleted = 1
	_WRITING   wCompleted = 2
)

func (w *wCompleted) intPoint() *int32 {
	return (*int32)(unsafe.Pointer(w))
}

func (w *wCompleted) load() int32 {
	return atomic.LoadInt32(w.intPoint())
}
func (w *wCompleted) writing() {
	atomic.StoreInt32(w.intPoint(), _WRITING.load())
}
func (w *wCompleted) completed() {
	atomic.StoreInt32(w.intPoint(), _COMPLETED.load())
}
func (w *wCompleted) empty() {
	atomic.StoreInt32(w.intPoint(), _EMPTY.load())
}
func (w *wCompleted) isCompleted() bool {
	return w.load() == _COMPLETED.load()
}
func (w *wCompleted) isWriting() bool {
	return w.load() == _WRITING.load()
}
