package main

import (
	"bufio"
	"log"
	"net"
	"net/http"

	"github.com/taimats/internal"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	req, err := internal.NewRequest(http.MethodGet, "http://localhost:8080", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("transfer-encoding-type", "chunked")
	if err := req.Write(conn); err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(conn)
	res, err := http.ReadResponse(r, req)
	if err != nil {
		log.Fatalf("http.ReadResponseでエラー:%v", err)
	}
	defer res.Body.Close()

	internal.HandleResponse(res, conn)
	log.Println("リクエスト処理が完了！！")
}
