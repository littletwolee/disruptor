package disruptor

type Job interface {
	Do() error
	CallBack(error)
}
type jobs []Job
