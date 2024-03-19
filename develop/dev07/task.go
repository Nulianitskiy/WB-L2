package main

import (
	"fmt"
	"sync"
	"time"
)

func or(channels ...<-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	end := make(chan struct{})

	once := &sync.Once{}

	for _, c := range channels {
		go func(c <-chan interface{}) {
			v := <-c
			once.Do(func() {
				out <- v
				end <- struct{}{}
			})
		}(c)
	}

	go func() {
		<-end
		close(out)
	}()

	return out

}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

func main() {

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(10*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	fmt.Printf("fone after %v", time.Since(start))
}
