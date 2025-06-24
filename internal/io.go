package internal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"slices"

	"github.com/fatih/color"
)

func ParsePNG(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	chunks, err := pngChunks(f)
	if err != nil {
		panic(err)
	}
	fmt.Println("==========================")
	fmt.Println("PNG画像の各チャンクを出力")
	fmt.Println("==========================")
	for _, c := range chunks {
		if err := showChunk(c); err != nil {
			panic(err)
		}
	}
	color.Green("PNG画像のパースが完了しました!!")
}

func ModifyPNG(path string) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	chunks, err := pngChunks(f)
	if err != nil {
		panic(err)
	}
	c := newChunk("teXt", []byte("inserted Text"))
	chunks = insertChunk(chunks, 2, c)

	nf := createPNGFile("new_demo.png", chunks)
	defer nf.Close()

	newChunks, err := pngChunks(nf)
	if err != nil {
		panic(err)
	}
	fmt.Println("==========================")
	fmt.Println("新PNG画像の各チャンクを出力")
	fmt.Println("==========================")
	for _, c := range newChunks {
		if err := showChunk(c); err != nil {
			panic(err)
		}
	}
	color.Green("PNG画像のパースが完了しました!!")
}

// PNG画像からチャンクを抽出して、そのReaderを返す。
// 特定の場所のみ読み取る必要があるため、内部ではSectionReaderを生成。
func pngChunks(f *os.File) ([]io.Reader, error) {
	//PNG画像の先頭8bitはシグネチャのため
	//読み取りはせずに進める。
	f.Seek(8, io.SeekStart)

	//チャンクごとにSectionReaderを生成。
	var offset int64 = 8
	var chunks []io.Reader
	for {
		//チャンクごとに入っているデータサイズは可変なので
		//dataSizeに格納する（それ以外は固定長のデータ）。
		var dataSize int32
		err := binary.Read(f, binary.BigEndian, &dataSize)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("binary.Read in PNGChunks returns an error:%w", err)
		}
		chunks = append(chunks, io.NewSectionReader(f, offset, 8+int64(dataSize)+4))
		offset, _ = f.Seek(4+int64(dataSize)+4, io.SeekCurrent)
	}
	return chunks, nil
}

func showChunk(chunk io.Reader) error {
	var dataSize int32
	if err := binary.Read(chunk, binary.BigEndian, &dataSize); err != nil {
		return fmt.Errorf("binary.Read in DumpChunk returns an error:%w", err)
	}
	buf := make([]byte, 4)
	chunk.Read(buf)
	fmt.Printf("{種類: %s, サイズ: %d Bytes}\n", buf, dataSize)
	fmt.Println()
	return nil
}

func newChunk(kind string, data []byte) io.Reader {
	var buf bytes.Buffer
	crc := crc32.NewIEEE()

	//先頭4byteに「長さ」を書き込み
	dataSize := len(data)
	binary.Write(&buf, binary.BigEndian, int32(dataSize))
	//後続の4byteに「種類」を書き込み
	io.WriteString(&buf, kind)

	//データ（ペイロード）を効率的に書き込みたいため、
	//同時書き込みができるwriterを用意。
	mw := io.MultiWriter(&buf, crc)
	mw.Write(data)

	//crcを末尾に書き込む
	binary.Write(&buf, binary.BigEndian, crc.Sum32())

	return &buf
}

func createPNGFile(fileName string, chunks []io.Reader) *os.File {
	f, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	//先頭8bitに固定長のシグネチャ
	io.WriteString(f, "\x89PNG\r\n\x1a\n")
	for _, c := range chunks {
		io.Copy(f, c)
	}
	color.Green("PNGファイルを作成しました! (ファイル名: %s)", fileName)
	return f
}

func insertChunk(chunks []io.Reader, pos int, c io.Reader) []io.Reader {
	return slices.Insert(chunks, pos, c)
}
