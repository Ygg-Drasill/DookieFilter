package frameReader

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"visunator/frame"
)

type FrameReader struct {
	b *bufio.Reader
}

func New(path string) *FrameReader {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return &FrameReader{b: bufio.NewReader(f)}
}

func (fr *FrameReader) Next() *frame.Frame[frame.DataPlayer] {
	l, _, err := fr.b.ReadLine()
	if err != nil {
		log.Println(err)
		return nil
	}

	newFrame := new(frame.Frame[frame.DataPlayer])
	err = json.Unmarshal(l, newFrame)
	if err != nil {
		log.Println(err)
		return nil
	}

	return newFrame
}
