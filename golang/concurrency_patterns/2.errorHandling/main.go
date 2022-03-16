package main

import (
	"fmt"
	"net/http"
)

type Result struct {
	Error    error
	Response *http.Response
}

func main() {
	callServerOnUrls()
}

func callServerOnUrls() {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan Result {
		response := make(chan Result)
		go func() {
			defer close(response)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{
					Error:    err,
					Response: resp,
				}
				select {
				case <-done:
					return
				case response <- result:
				}
			}
		}()
		return response
	}

	done := make(chan interface{})
	defer close(done)
	urls := []string{
		"https://www.google.com",
		"https://badhost",
	}
	for response := range checkStatus(done, urls...) {
		if response.Error != nil {
			fmt.Println(response.Error)
			continue
		}
		fmt.Printf("Response: %v\n", response.Response.Status)
	}
}
