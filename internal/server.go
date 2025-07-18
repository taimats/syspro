package internal

import (
	"bufio"
	"bytes"
	"compress/gzip"
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"slices"
	"strings"
	"time"
)

// クライアントからのDialをセッションとしてとらえ、
// １コネクションを1セッションとする。1セッション内に
// 複数回リクエストが来れば、それぞれhandleRequestで対応。
// responseはwriteToConnを共通出口として処理。
func HandleSession(conn net.Conn, content []byte) {
	defer conn.Close()
	fmt.Printf("Accept %v\n", conn.RemoteAddr())

	queue := make(chan *http.Response, 50)
	defer close(queue)
	go writeToConn(conn, queue)

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

		go handleRequest(req, conn, queue, content)
	}
}

// キューのなかに複数のレスポンスを内包し、順にレスポンスを返す。
func writeToConn(conn net.Conn, queue chan *http.Response) {
	defer conn.Close()
	for response := range queue {
		response.Write(conn)
	}
}

// リクエスト内容を読み込み、レスポンスをキューに送信。
// レスポンス処理は以下に対応している。
// ・gzipでの圧縮
// ・チャンク形式の送付
func handleRequest(req *http.Request, conn net.Conn, queue chan *http.Response, content []byte) {
	if isChunkedTransferOK(req) {
		//http.ResponseはContent-Length指定がないと、Connection Closeをクライアントに送る。
		//だが、チャンク形式の送付ではBodyを複数回に分けて送付するため、現段階ではサイズが未定である。
		//よって、暫定的にBodyとContent-Length以外の文字列を直接connに書き込んでいる。
		res := strings.Join([]string{
			"HTTP/1.1 200 OK",
			"Content-Type: text/plain",
			"Transfer-Encoding: chunked",
			"", "", //Bodyまでは2行分空白を作る必要がある
		}, "\r\n")
		//この場合は、キューを通していない。
		fmt.Fprint(conn, res)

		r := bytes.NewReader(content)
		sc := bufio.NewScanner(r)
		for sc.Scan() {
			fmt.Fprintf(conn, "%s\r\n", sc.Bytes())
		}
		fmt.Fprint(conn, io.EOF)
		return
	}
	res := &http.Response{
		StatusCode:    http.StatusOK,
		ProtoMajor:    1,
		ProtoMinor:    1,
		ContentLength: int64(len(content)),
		Body:          io.NopCloser(bytes.NewReader(content)),
		Header:        make(http.Header),
	}
	if isGzipOK(req) {
		var buf bytes.Buffer
		w := gzip.NewWriter(&buf)
		w.Write(content)
		w.Close()

		res.Body = io.NopCloser(&buf)
		res.ContentLength = int64(buf.Len())
		res.Header.Set("Content-Encoding", "gzip")
	}
	queue <- res
}

func isGzipOK(req *http.Request) bool {
	hs := req.Header["Accept-Encoding"]
	return slices.Contains(hs, "gzip")
}

func isChunkedTransferOK(req *http.Request) bool {
	return req.Header.Get("transfer-encoding-type") == "chunked"
}

// UDPプロトコルで通信を行うお試し関数。
func UDPTrial(address string) {
	conn, err := net.ListenPacket("udp", address)
	if err != nil {
		log.Fatalf("net.ListenPacket failed: (error: %v)", err)
	}
	fmt.Println("Server is listening on", address)
	handleUDPSession(conn)
}

func handleUDPSession(conn net.PacketConn) {
	defer conn.Close()
	buf := make([]byte, 1500)
	for {
		length, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatalf("conn.ReadFrom failed: (error: %v)", err)
		}
		fmt.Printf("{\ncontent: %s\naddress: %s\n}\n", buf[:length], addr)

		res := `Hello, client`
		_, err = conn.WriteTo([]byte(res), addr)
		if err != nil {
			log.Fatalf("conn.WriteTo failed: (error: %v)", err)
		}
	}
}
