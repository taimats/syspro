package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/taimats/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: can get the size of a specified path")
		os.Exit(0)
	}
	if len(os.Args) > 2 {
		color.Red("only one file name is available")
		os.Exit(-1)
	}
	internal.Size(os.Args[1])
}
