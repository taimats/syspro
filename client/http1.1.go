package main

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	messages := []string{
		"apple",
		"banana",
		"orange",
	}
	count := 0

	for {
		if count == len(messages) {
			log.Println("すべてのメッセージの送信が完了!!")
			break
		}
		conn, err := net.Dial("tcp", "localhost:8080")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		req, err := http.NewRequest(
			http.MethodPost,
			"http://localhost:8080",
			strings.NewReader(messages[0]),
		)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Accept-Encoding", "gzip")
		if err = req.Write(conn); err != nil {
			log.Fatal(err)
		}

		res, err := http.ReadResponse(bufio.NewReader(conn), req)
		if err != nil {
			fmt.Println("Retry...")
			continue
		}
		defer res.Body.Close()
		handleResponse(res)

		count++
	}
}

func handleResponse(res *http.Response) {
	if isGzipRequired(res.Header.Get("Content-Encoding")) {
		r, err := gzip.NewReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		io.Copy(os.Stdout, r)
		return
	}
	io.Copy(os.Stdout, res.Body)
}

func isGzipRequired(value string) bool {
	return value == "gzip"
}
