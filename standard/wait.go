package standard

import "fmt"

func Wait() {
	var done = make(chan bool)
	var results = make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			results <- i
		}
		close(results)
		done <- true
	}()
	var hit = make([]int, 0)
	for i := range results {
		hit = append(hit, i)
	}
	<-done
	fmt.Println("done", hit)

}
