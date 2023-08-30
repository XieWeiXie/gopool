package standard

import (
	"fmt"
	"sync"
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

func NilChannel1() {
	var c chan struct{}
	<-c
}

func NilChannel2() {
	var c chan struct{}
	c <- struct{}{}
}

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
