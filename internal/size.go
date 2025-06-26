package internal

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Size(path string) {
	info, err := os.Stat(path)
	if err != nil {
		color.Red("ファイルパスが不正です")
		os.Exit(-1)
	}
	fmt.Printf("サイズ: %d bytes\n", info.Size())
}
