package main

import (
	"fmt"
	"github.com/XieWeiXie/gopool"
	"sync"
	"time"
)

func main() {
	var start = time.Now()
	var i = 0
	var wg sync.WaitGroup
	for i < gopool.TimesMillion {
		wg.Add(1)
		go func(i int) {
			gopool.DoWork(gopool.Task{Message: fmt.Sprintf("%d", i)})
			wg.Done()
		}(i)
		i++
	}
	fmt.Println(time.Now().Sub(start).Milliseconds())
	wg.Wait()
}
