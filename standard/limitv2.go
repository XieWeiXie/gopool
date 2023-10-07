package standard

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

func LimitV2Worker(ch chan int, wg *sync.WaitGroup) {
	for c := range ch {
		fmt.Println(c, runtime.NumGoroutine())
		wg.Done()
	}
}
func LimitV2() {
	var c = make(chan int)
	var num = 3
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		go LimitV2Worker(c, &wg)
	}
	for i := 0; i < math.MaxInt; i++ {
		wg.Add(1)
		c <- i
	}
	wg.Wait()
}
