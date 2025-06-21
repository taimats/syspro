package main

import (
	_ "embed"
	"fmt"
	"log"
	"net"

	"github.com/taimats/internal"
)

//go:embed content.txt
var content []byte

// tcpを使った自作のhttp1.1対応のサーバー。
// encodingではgzipを許可している。
func main() {
	address := "localhost:8080"
	l, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to open Listener:(error: %v)", err)
	}
	fmt.Printf("Server is listening on %v\n", address)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go internal.HandleRequest(conn, content)
	}
}
