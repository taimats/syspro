package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"slices"
	"strconv"
	"strings"
	"time"
)

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
		go handleRequest(conn)
	}
}

// リクエスト内容を読み込み、レスポンスに書き込む。
func handleRequest(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Accept %v\n", conn.RemoteAddr())

	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		req, err := http.ReadRequest(bufio.NewReader(conn))
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			neterr, ok := err.(net.Error)
			if ok && neterr.Timeout() {
				log.Println("time out")
				break
			}
			log.Println(err)
			break
		}
		dump, err := httputil.DumpRequest(req, true)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(dump))

		content := "Hello world\n"
		res := &http.Response{
			Status:        strconv.Itoa(http.StatusOK),
			ProtoMajor:    1,
			ProtoMinor:    1,
			ContentLength: int64(len(content)),
			Body:          io.NopCloser(strings.NewReader(content)),
			Header:        make(http.Header),
		}
		if isGzipAcceptable(req) {
			content = "Hello world (Gzipped)\n"

			var buf bytes.Buffer
			w := gzip.NewWriter(&buf)
			io.WriteString(w, content)
			w.Close()

			res.Body = io.NopCloser(&buf)
			res.ContentLength = int64(buf.Len())
			res.Header.Set("Content-Encoding", "gzip")
		}
		res.Write(conn)
	}
}

func isGzipAcceptable(req *http.Request) bool {
	hs := req.Header["Accept-Encoding"]
	return slices.Contains(hs, "gzip")
}
