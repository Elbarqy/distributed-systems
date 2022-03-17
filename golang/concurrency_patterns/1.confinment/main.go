package main

import (
	"fmt"
)

func main() {
}

func orChannel() {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		case 0:
			return nil
		case 1:
			return channels[0]
		}
		orDone := make(chan interface{})
		//Creates a go routine each 3 channels
		//which results in more overhead
		go func() {
			defer close(orDone)
			switch len(channels) {
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}
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
