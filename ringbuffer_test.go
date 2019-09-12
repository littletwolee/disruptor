package disruptor

import (
	"fmt"
	"testing"
	"time"
)

func Test_New(t *testing.T) {
	rb, err := New(4, 2048)
	if err != nil {
		t.Fatal(err)
	}
	go rb.Start()
	now := time.Now()
	rb.Writing()
	for i := 0; i < 10000000; i++ {
		// fmt.Println(i)
		rb.Write(&JobTest{ID: i})
	}
	fmt.Println("write completed")
	rb.Completed()
	rb.Wait()
	fmt.Println(time.Now().Sub(now))
}

type JobTest struct {
	ID int
}

func (j *JobTest) Do() error {
	// log.Println(j.ID)
	return nil
}
func (j *JobTest) CallBack(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
