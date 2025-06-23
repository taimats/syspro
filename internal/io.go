package internal

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"

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
