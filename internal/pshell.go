package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/google/shlex"
)

func PShell() {
	sc := bufio.NewScanner(os.Stdin)
	for {
		prompt()
		cmd := parseCMD(sc)
		parseShell(cmd)
	}
}

func prompt() {
	color.Blue("コマンドを入力ください")
	fmt.Print(">>> ")
}

func parseCMD(sc *bufio.Scanner) (cmd string) {
	if next := sc.Scan(); !next {
		color.Green("終了します")
	}
	cmd = sc.Text()
	if cmd == "exit" || cmd == "q" || cmd == "quit" {
		color.Green("終了します")
		os.Exit(0)
	}
	return cmd
}

func parseShell(cmd string) {
	l := shlex.NewLexer(strings.NewReader(cmd))
	for {
		token, err := l.Next()
		if err != nil {
			break
		}
		fmt.Println("取得したtoken:", token)
	}
}
