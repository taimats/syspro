package main

import (
	"fmt"

	"github.com/taimats/internal"
)

func main() {
	done := make(chan struct{})
	defer close(done)

	intStream := internal.Generator(done, 1, 2, 3, 4)
	pipeline := internal.Multiply(done, internal.Add(done, internal.Multiply(done, intStream, 2), 1), 2)
	for v := range pipeline {
		fmt.Println(v)
	}
}
