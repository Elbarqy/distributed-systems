package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	randomSleeper := func(
		done <-chan interface{},
		id int,
	) <-chan interface{} {
		sleeper := make(chan interface{})
		go func() {
			for {
				select {
				case <-done:
					return
				case sleeper <- id:
					period := rand.Intn(2000)
					fmt.Printf("routine of %d has started sleeping for %d ms \n", id, period)
					time.Sleep(time.Duration(period) * time.Millisecond)
				}
			}
		}()
		return sleeper
	}

	done := make(chan interface{})
	defer close(done)
	workerCount := runtime.NumCPU()
	workerStreams := make([]<-chan interface{}, workerCount)
	for i := 0; i < workerCount; i++ {
		workerStreams[i] = randomSleeper(done, i)
	}

	for id := range fanIn(done, workerStreams...) {
		fmt.Println(id)
	}
	fmt.Printf("Search took: %v", time.Since(start))

}

//This only works well whe
//Collectioning all streams along one bus FAN IN

func fanIn(
	done <-chan interface{},
	channels ...<-chan interface{},
) <-chan interface{} {
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}
