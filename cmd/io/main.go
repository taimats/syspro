package main

import "github.com/taimats/internal"

func main() {
	path := "./cmd/io/demo.png"
	internal.ModifyPNG(path)
}
