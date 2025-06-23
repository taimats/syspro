package main

import (
	"github.com/taimats/internal"
)

func main() {
	addr := "localhost:8080"
	internal.UDPRequest(addr)
}
