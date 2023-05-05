package gopool

import (
	"runtime"
	"sync"
	"testing"
	"time"
)

const (
	_   = 1 << (10 * iota)
	KiB // 1024
	MiB // 1048576
)

var curMem uint64

func TestWorker(t *testing.T) {
	// 不计成本的创建 go 协程
	start := time.Now()
	var wg sync.WaitGroup
	for j := 0; j < TimesMillion; j++ {
		wg.Add(1)
		go func() {
			DoWork(Task{Message: "123"})
			wg.Done()
		}()
	}
	wg.Wait()
	mem := runtime.MemStats{}
	runtime.ReadMemStats(&mem)
	curMem := mem.TotalAlloc/MiB - curMem
	t.Logf("memory usage:%d MB", curMem)                                     // 472 MB
	t.Logf("cost time: %d millSecond", time.Now().Sub(start).Milliseconds()) // 2138 millSeconds
}
