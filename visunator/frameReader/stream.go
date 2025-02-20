package frameReader

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"visunator/types"
)

type FrameReader struct {
	buff *bufio.Reader
	file *os.File
}

func New(path string) *FrameReader {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return &FrameReader{buff: bufio.NewReader(f), file: f}
}

func (fr *FrameReader) Next() *types.Frame[types.DataPlayer] {
	l, _, err := fr.buff.ReadLine()
	if err != nil {
		if err == io.EOF {
			cerr := fr.file.Close()
			if cerr != nil {
				log.Fatal(cerr)
			}
			return nil
		}
		log.Println(err)
		return nil
	}

	newFrame := new(types.Frame[types.DataPlayer])
	err = json.Unmarshal(l, newFrame)
	if err != nil {
		log.Println(err)
		return nil
	}

	if len(newFrame.Data[0].Ball.Xyz) == 0 { //TODO: maybe return signal later :)
		fmt.Println("hello")
		return fr.Next()
	}

	if len(newFrame.Data[0].HomePlayers) == 0 {
		return fr.Next()
	}

	return newFrame
}
