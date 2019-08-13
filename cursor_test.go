package disruptor

import (
	"fmt"
	"testing"
)

func Test_newCursors(t *testing.T) {
	cs := newCursors(2, 10)
	for _, c := range cs.cs {
		fmt.Printf("r:%v,w:%v\n", *c.r, *c.w)
	}
}

func Test_cursorRead(t *testing.T) {
	r, w, s, e := int64(0), int64(1), int64(0), int64(2)
	c := _cursor{&r, &w, s, e}
	if !c.read() {
		t.Fatalf("error")
	}
}

func Test_cursorWrite(t *testing.T) {
	r, w, s, e := int64(1), int64(1), int64(0), int64(2)
	c := _cursor{&r, &w, s, e}
	if !c.write() {
		t.Fatalf("error")
	}
}
