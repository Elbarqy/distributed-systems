package main

import "fmt"

func main() {
	pipeline()
}

func pipeline() {
	generator := func(
		done <-chan interface{},
		integers ...int,
	) <-chan int {
		intStream := make(chan int)
		go func() {
			defer close(intStream)
			for _, i := range integers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}
	multiply := func(
		done <-chan interface{},
		intStream <-chan int,
		multiplier int,
	) <-chan int {
		multiplyStream := make(chan int)
		go func() {
			defer close(multiplyStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multiplyStream <- multiplier * i:
				}
			}
		}()
		return multiplyStream
	}
	add := func(
		done <-chan interface{},
		intStream <-chan int,
		additive int,
	) <-chan int {
		addStream := make(chan int)
		go func() {
			defer close(addStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addStream <- i + additive:
				}
			}
		}()
		return addStream
	}

	done := make(chan interface{})
	defer close(done)
	intStream := generator(done, 1, 2, 3, 4)
	pipelineStream := multiply(done, add(done, intStream, 5), 2)
	for v := range pipelineStream {
		fmt.Println(v)
	}
}
