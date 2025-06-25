package internal

import (
	"fmt"
	"os"
)

const usage = `<Usage>
 ln registers a physical or simbolic name to a certain path.
 It behaves like Unix "ln" command, so does its usage in the following way:
   ln <options> target(=path) dist(=path)

<Options>
 -s: registers a simbolic name to dist
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
	dist   string
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
	if len(ln.args[1:]) < 2 {
		ln.Usage()
		os.Exit(2)
	}

	flag := ln.args[1]
	if flag == "-h" || flag == "-help" {
		ln.Usage()
		os.Exit(0)
	}
	if flag == "-s" {
		if len(ln.args[2:]) != 2 {
			ln.Usage()
			os.Exit(2)
		}
		ln.simbol = true
		ln.target = ln.args[2]
		ln.dist = ln.args[3]
		return
	}

	ln.target = ln.args[1]
	ln.dist = ln.args[2]
}

func (ln *LnCMD) Run() error {
	if ln.isSimbolic() {
		if err := os.Symlink(ln.target, ln.dist); err != nil {
			return err
		}
		return nil
	}
	if err := os.Link(ln.target, ln.dist); err != nil {
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
