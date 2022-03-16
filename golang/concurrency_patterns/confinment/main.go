package main

import (
	"fmt"
)

func main() {
}

func orchastrated_routine() {
	// <- marks the function as a recieve only
	dowork := func(
		done <-chan interface{},
	) <-chan int {
		// create a buffer channel of size 5
		// can use waitgroup instead
		results := make(chan int)
		go func() {
			defer fmt.Println("exited the routine")
			defer close(results)
			for {
				select {
				case results <- 1:
				case <-done:
					return
				}
			}
		}()
		return results
	}

	consumer := func(results <-chan int) {
		for i := 0; i <= 3; i++ {
			fmt.Printf("%d: %d\n", i, <-results)
		}
		fmt.Println("Done recieving")
	}

	done := make(chan interface{})
	defer close(done)

	results := dowork(done)
	consumer(results)
}
