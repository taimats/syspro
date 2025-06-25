package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/taimats/internal"
)

func main() {
	ln := internal.NewLnCMD(os.Args)
	ln.Parse()
	if err := ln.Run(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	color.Green("ファイル名を変更しました!!")
}
