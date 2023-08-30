package standard

import (
	"fmt"
	"sync"
	"time"
)

// Notify 通知 信号传递
func Notify() {
	var c = make(chan int)
	go func() {
		c <- 100
	}()
	fmt.Println(<-c)
}

type Signal struct {
}

func Notify1Vs() {
	c := Spawn(func() {
		fmt.Println("Hello Spawn")
	})
	<-c
}

func Spawn(f func()) <-chan Signal {
	var c = make(chan Signal)
	go func() {
		f()
		c <- Signal{}
	}()
	return c
}

func Notify1VsN() {
	var g = make(chan Signal)
	c := Spawn1VsN(func() { fmt.Println("Hello World ", time.Now()) }, g)
	close(g)
	<-c
}

func Spawn1VsN(f func(), global <-chan Signal) <-chan Signal {
	var c = make(chan Signal)
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-global
			f()
		}()
	}

	go func() {
		wg.Wait()
		c <- Signal{}
	}()
	return c
}

type Counter struct {
	c chan int
	i int
}

func NewCounter() *Counter {
	var c Counter = Counter{
		c: make(chan int),
		i: 0,
	}
	go func() {
		for {
			c.i++
			c.c <- c.i
		}

	}()
	return &c
}

func (c *Counter) Increase() int {
	return <-c.c
}

func CounterExample() {
	var counter = NewCounter()
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			fmt.Println(counter.Increase())
		}()
	}
	wg.Wait()
}

var queue = make(chan int, 100)

func Producer() {
	for i := 0; i < 1000; i++ {
		queue <- i
	}
	close(queue)
}

func Consumer() {
	for {
		select {
		case c, ok := <-queue:
			if !ok {
				fmt.Println("All Done")
				return
			}
			fmt.Println(c)
		}
	}

}

func QueueExample() {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		Producer()
	}()
	go func() {
		defer wg.Done()
		Consumer()
	}()
	wg.Wait()
}
