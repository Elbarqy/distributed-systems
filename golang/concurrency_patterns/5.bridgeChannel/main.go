package main

import (
	"fmt"
)

func main() {
	bridge_impl()
}

func bridge_impl() {
	done := make(chan interface{})
	defer close(done)
	bridge := func(
		done <-chan interface{},
		chanStream <-chan <-chan interface{},
	) <-chan interface{} {
		valStream := make(chan interface{})
		go func() {
			defer close(valStream)
			for {
				var stream <-chan interface{}
				select {
				case maybeStream, ok := <-chanStream:
					if ok == false {
						return
					}
					stream = maybeStream
				case <-done:
					return
				}
				for val := range stream {
					select {
					case valStream <- val:
					case <-done:
					}
				}
			}
		}()
		return valStream
	}
	genVals := func() <-chan <-chan interface{} {
		chanStream := make(chan (<-chan interface{}))
		go func() {
			defer close(chanStream)
			for i := 0; i < 10; i++ {
				stream := make(chan interface{})
				//a workaround stream := make(chan interface{},1)
				//make it so that stream <- i is non blocking_

				//Next line is blocking which means it won't move for the next line
				//till it's consumed
				chanStream <- stream
				stream <- i
				close(stream)
			}
		}()
		return chanStream
	}
	for v := range bridge(done, genVals()) {
		fmt.Printf("%v ", v)
	}
}
