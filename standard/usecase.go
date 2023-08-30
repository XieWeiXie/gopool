package standard

import "fmt"

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
