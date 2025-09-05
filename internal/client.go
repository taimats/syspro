package internal

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"time"
)

func NewRequest(method string, url string, body io.Reader, encodings []string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if !(len(encodings) == 0) {
		for _, e := range encodings {
			req.Header.Set("Accept-Encoding", e)
		}
	}
	return req, nil
}

func HandleResponse(res *http.Response, conn net.Conn) {
	if isGzipRequired(res) {
		r, err := gzip.NewReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("====================")
		fmt.Println("gzipのデコードを実施 ")
		fmt.Println("====================")
		io.Copy(os.Stdout, r)
		fmt.Println()
		return
	}
	if isChunkedTransfer(res) {
		r := bufio.NewScanner(conn)
		fmt.Println("=======================")
		fmt.Println("chunkedのデコードを実施 ")
		fmt.Println("=======================")
		for r.Scan() {
			line := r.Bytes()
			fmt.Printf("{\nサイズ:%dbyte\n取得:%s\n}\n", len(line), string(line))
		}
		return
	}
	io.Copy(os.Stdout, res.Body)
	fmt.Println()
}

func isChunkedTransfer(res *http.Response) bool {
	return len(res.TransferEncoding) >= 1 && res.TransferEncoding[0] == "chunked"
}

func isGzipRequired(res *http.Response) bool {
	return res.Header.Get("Content-Encoding") == "gzip"
}

func UDPRequest(address string) {
	conn, err := net.Dial("udp4", address)
	if err != nil {
		log.Fatalf("net.Dial failed: (error: %v)", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(`Hello, server`))
	if err != nil {
		log.Fatalf("conn.Write in UDPRequest failed: (error: %v)", err)
	}
	fmt.Println("サーバーにメッセージを送信")

	buf := make([]byte, 1500)
	length, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("conn.Read in UDPRequest failed: (error: %v)", err)
	}
	fmt.Printf("{\ncontent: %s\n}\n", buf[:length])
}

type ExtendedTransport struct {
	transport http.RoundTripper

	maxRetryCount     int
	currentRetryCount int

	maxReqCount    int
	perMilliSecond int64
	window         Window
}

func NewExtendedTransport(transport http.RoundTripper, maxRetryCount int, maxReqCount int, perMilliSecond int64) *ExtendedTransport {
	return &ExtendedTransport{
		transport:         transport,
		maxRetryCount:     maxRetryCount,
		currentRetryCount: 0,
		maxReqCount:       maxReqCount,
		perMilliSecond:    perMilliSecond,
		window: Window{
			dueTime:  int64(0),
			reqCount: 0,
		},
	}
}

func (e *ExtendedTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	for {
		now := time.Now().UnixMilli()
		wait := e.window.dueTime - now
		//if wait is a minus number, now should be over the due time.
		if wait < 0 {
			//a new time span starting here
			e.window = Window{
				dueTime:  now + e.perMilliSecond,
				reqCount: 0,
			}
			break
		}
		if e.window.reqCount < e.maxReqCount {
			break
		}
		time.Sleep(time.Duration(wait) * time.Millisecond)
	}
	e.window.reqCount++

	var res *http.Response
	var err error
	for {
		res, err = e.transport.RoundTrip(r)
		if res != nil && res.StatusCode < http.StatusInternalServerError {
			break
		}
		e.currentRetryCount++
		if e.currentRetryCount > e.maxRetryCount {
			break
		}
		//Exponential backoff
		time.Sleep(time.Second * time.Duration(math.Pow(2, float64(e.currentRetryCount))))
	}
	//need to be initialized so that a conn can be used multiple times
	e.currentRetryCount = 0
	return res, err
}

// Fixed window counter algorithm
type Window struct {
	dueTime  int64
	reqCount int
}
