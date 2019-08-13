package disruptor

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	rb, err := New(10, 2048)
	if err != nil {
		t.Fatal(err)
	}
	go rb.Start()
	now := time.Now()
	for i := 0; i < 10000000; i++ {
		rb.Write(&JobTest{ID: i})
	}
	go func(now time.Time, rb *ringBuffer) {
		fmt.Printf("t:%v,len:%d\n", now.Sub(time.Now()), rb.buf)
		time.Sleep(time.Second)
	}(now, rb)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
LOOP:
	for {
		select {
		case <-signalChan:
			rb.Stop()
			break LOOP
		}
	}
}

type JobTest struct {
	ID int
}

func (j *JobTest) Do() error {
	// fmt.Println(j.ID)
	return nil
}
func (j *JobTest) CallBack(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
