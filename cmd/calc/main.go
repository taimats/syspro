package main

import (
	"runtime"
	"sync"

	"github.com/taimats/internal"
)

func main() {
	l := internal.NewLoan(1, 400000000, 0.011)
	years := make(chan int, 35)
	defer close(years)
	for i := 1; i < 36; i++ {
		years <- i
	}

	var wg sync.WaitGroup
	wg.Add(35)
	for range runtime.NumCPU() {
		go internal.Worker(l, years, &wg)
	}

	wg.Wait()
}
