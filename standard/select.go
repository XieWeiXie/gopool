package standard

import (
	"fmt"
	"time"
)

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

func Select() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 20; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

func SelectV2() {
	c := make(chan int)
	quit := make(chan int)
	go fibonacci(c, quit)
	for i := 0; i < 20; i++ {
		fmt.Println(<-c)
	}
	quit <- 0
}

func Example() {
	var (
		c       = make(chan int)
		returnC = make(chan int)
		done    = make(chan struct{})
	)
	go func() {
		for {
			select {
			case id, ok := <-c:
				returnC <- id
				if !ok {
					done <- struct{}{}
					break
				}
			}
		}
	}()
	var ids = make([]int, 0)
	go func() {
		for id := range returnC {
			ids = append(ids, id)
		}
	}()
	for _, id := range []int{1, 2, 3, 4, 5, 6, 7, 8, 9} {
		c <- id
	}
	close(c)
	<-done

	fmt.Println(ids)
}

func Worker(jobs <-chan int, results chan<- int, done chan struct{}) {
	for {
		select {
		case r, more := <-jobs:
			if !more {
				done <- struct{}{}
				break
			}
			time.Sleep(1 * time.Second)
			results <- r * 2
		}
	}
}

func Example2() {
	job := make(chan int)
	results := make(chan int)
	done := make(chan struct{})

	for i := 0; i < 5; i++ {
		go Worker(job, results, done)
	}

	go func() {
		for {
			select {
			case r, ok := <-results:
				if !ok {
					break
				}
				fmt.Println(r)
			}
		}
	}()

	for i := 0; i <= 100; i++ {
		job <- i
	}
	close(job)
	<-done

}
