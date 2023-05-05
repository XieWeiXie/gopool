package gopool

import "time"

const (
	workTime = time.Second * 1
)

const (
	TimesMillion    = 1e6      // 一百万
	TimesTenMillion = 1e6 * 10 // 一千万
)

type Task struct {
	Message string
}

func DoWork(task Task) {
	time.Sleep(workTime)
}
