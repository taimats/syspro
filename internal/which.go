package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

func Which(args []string) {
	if len(args) != 2 {
		color.Red("コマンド名を指定ください")
		os.Exit(1)
	}
	cmdName := args[1]
	if !strings.Contains(cmdName, ".exe") {
		cmdName = cmdName + ".exe"
	}

	paths := filepath.SplitList(os.Getenv("PATH"))
	for _, p := range paths {
		execpath := filepath.Join(p, cmdName)
		_, err := os.Stat(execpath)
		if err == nil {
			fmt.Println(execpath)
			return
		}
	}
	color.Green("一致する実行ファイルはありません")
}
