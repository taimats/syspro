package main

import (
	"context"
	_ "embed"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
)

//go:embed content.txt
var content []byte

func main() {
	path := filepath.Join(os.TempDir(), "unix-socket")
	defer os.Remove(path)
	l, err := net.Listen("unix", path)
	if err != nil {
		panic(err)
	}
	defer l.Close()

	fmt.Println("Server is listening on file", path)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-ctx.Done()
	fmt.Println("プログラムを停止します")
}
