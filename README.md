# go-pool


### 1、channel

多协程中通道通信。

- unbuffered channel 发送接收需要同时准备，否则阻塞
- buffered channel 空时，读取阻塞；满时，写入阻塞


```go
package standard

import (
	"fmt"
)

// UnbufferedEmptyRead 无缓冲通道 读阻塞
func UnbufferedEmptyRead() {
	var c = make(chan struct{})
	<-c
}

// UnbufferedEmptyWrite 无缓冲通道 写入阻塞
func UnbufferedEmptyWrite() {
	var c = make(chan struct{})
	c <- struct{}{}
}

// UnbufferedOK 无缓冲通道 读写在不同的协程中
func UnbufferedOK() {
	var c = make(chan struct{})
	go func() {
		c <- struct{}{}
	}()
	select {
	case ok := <-c:
		fmt.Println("OK", ok)
		break
	}
}

// BufferedEmptyRead 缓冲通道为空时读取 阻塞
func BufferedEmptyRead() {
	var c = make(chan struct{}, 1)
	<-c
}

// BufferedFullWrite 缓冲通道满时写入 阻塞
func BufferedFullWrite() {
	var c = make(chan struct{}, 1)
	c <- struct{}{}
	c <- struct{}{}
}

```


- nil channel 

读/写都阻塞

```go
package main

func NilChannel1() {
	var c chan struct{}
	<-c
}

func NilChannel2() {
	var c chan struct{}
	c <- struct{}{}
}

```

- 需不需要主动关闭 channel

不需要。

gc 会自动操作。

如果需要关闭 channel ，则在发送方主动关闭。
```go
package main

import (
	"fmt"
	"sync"
)

func CloseChannelOK() {
	var c = make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		i := 30
		for i > 0 {
			c <- i
			i--
		}
		close(c)
	}()
	go func() {
		defer wg.Done()
		for i := range c {
			fmt.Println("Receive ", i)
		}
	}()
	wg.Wait()
}


```


- panic 场景

```go
package main

import (
	"fmt"
)

func NilChannelClose() {
	var c chan struct{}
	close(c)
}

func SendValueToClosedChannel() {
	var c = make(chan struct{})
	go func() {
		c <- struct{}{}
		close(c)
	}()
	select {
	case ok := <-c:
		fmt.Println(ok)
		c <- struct{}{}
	}

}

func CloseClosedChannel() {
	var c = make(chan struct{})
	close(c)
	close(c)
}
```




### 通用用法

无缓冲通道：

- 信号通知 1 v 1; 1 v N
- 锁

```go

package standard

import (
	"fmt"
	"sync"
)

// Notify 通知 信号传递
func Notify() {
	var c = make(chan int)
	go func() {
		c <- 100
	}()
	fmt.Println(<-c)
}

type Signal struct {}

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


```

缓冲通道：

- 消息队列: 生产者不断的往队列中放数据，消费者不断的往队列中取数据

```go
package main

import (
	"fmt"
	"sync"
)

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
```