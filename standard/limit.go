package standard

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

func limitGoroutineWorker(c chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("Num goroutine ", runtime.NumGoroutine())
	<-c
}

func LimitGoroutine() {
	var buffered = make(chan struct{}, 3)
	var wg sync.WaitGroup
	for i := 0; i < math.MaxInt; i++ {
		wg.Add(1)
		buffered <- struct{}{}
		go limitGoroutineWorker(buffered, &wg)
	}
	wg.Wait()
}
