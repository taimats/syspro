package internal

import (
	"fmt"
	"os"
)

const usage = `<Usage>
 ln registers a physical or simbolic name to a certain path.
 It behaves like Unix "ln" command, so does its usage in the following way:
   ln <options> target(=path) dest(=path)

<Options>
 -s: registers a simbolic name to dest
 -h, -help: shows info for this command usage`

type Command interface {
	Name() string
	Parse()
	Run() error
	Usage()
}

var _ Command = (*LnCMD)(nil)

type LnCMD struct {
	args   []string
	target string
	dest   string
	simbol bool
}

func NewLnCMD(args []string) *LnCMD {
	return &LnCMD{
		args:   args,
		simbol: false,
	}
}

func (ln *LnCMD) Name() string {
	return ln.args[0]
}

func (ln *LnCMD) Parse() {
	//a process when only a command name (plus alpha) is passed
	if len(ln.args[0:]) < 2 {
		ln.Usage()
		os.Exit(0)
	}
	//parsing flags
	flag := ln.args[1]
	if flag == "-h" || flag == "-help" {
		ln.Usage()
		os.Exit(0)
	}
	if flag == "-s" {
		if len(ln.args[2:]) != 2 {
			fmt.Println("ファイル名を指定ください")
			os.Exit(2)
		}
		ln.simbol = true
		ln.target = ln.args[2]
		ln.dest = ln.args[3]
		return
	}
	//in the case of no flags
	if len(ln.args[1:]) != 2 {
		fmt.Println("ファイル名を指定ください")
		os.Exit(2)
	}
	ln.target = ln.args[1]
	ln.dest = ln.args[2]
}

func (ln *LnCMD) Run() error {
	if ln.isSimbolic() {
		if err := os.Symlink(ln.target, ln.dest); err != nil {
			return err
		}
		return nil
	}
	if err := os.Link(ln.target, ln.dest); err != nil {
		return err
	}
	return nil
}

func (ln *LnCMD) Usage() {
	fmt.Println(usage)
}

func (ln *LnCMD) isSimbolic() bool {
	return ln.simbol
}
