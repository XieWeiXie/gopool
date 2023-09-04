package main

import "fmt"

func main() {
	jobs := make(chan int)
	done := make(chan bool)

	go func() {
		for {
			select {
			case j, more := <-jobs:
				if !more {
					done <- true
					return
				}
				fmt.Println("received job", j)

			}
		}
	}()

	for j := 1; j <= 3; j++ {
		jobs <- j
		fmt.Println("sent job", j)
	}
	close(jobs)
	fmt.Println("sent all jobs")

	<-done
}
