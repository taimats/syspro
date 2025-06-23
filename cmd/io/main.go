package main

import "github.com/taimats/internal"

func main() {
	path := "./cmd/io/demo.png"
	internal.ParsePNG(path)
}
