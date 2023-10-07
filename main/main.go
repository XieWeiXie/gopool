package main

import (
	"fmt"
	"math"
	"runtime"
	"sync"
)

var wg sync.WaitGroup

func a(c chan bool, i int) {
	defer wg.Done()
	fmt.Println(fmt.Sprintf("i %d num goroutine %d", i, runtime.NumGoroutine()))
	<-c

}

func main() {
	task := math.MaxInt
	ch := make(chan bool, 3)
	for i := 0; i < task; i++ {
		wg.Add(1)
		ch <- true
		go a(ch, i)
	}
	wg.Wait()
}
