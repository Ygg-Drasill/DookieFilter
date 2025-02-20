package frameReader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"visunator/types"
)

type FrameReader struct {
	buff       *bufio.Reader
	file       *os.File
	prefixBuff bytes.Buffer
}

func New(path string) *FrameReader {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	return &FrameReader{buff: bufio.NewReader(f), file: f}
}

func (fr *FrameReader) Next() *types.Frame[types.DataPlayer] {
	lBuff, isPrefix, err := fr.buff.ReadLine()
	if isPrefix {
		n, err := fr.prefixBuff.Write(lBuff)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("buffering %d bytes as prefix", n)
		return fr.Next()
	}

	prefixSize := fr.prefixBuff.Available()
	if prefixSize > 0 {
		fr.prefixBuff.Write(lBuff)
		all := make([]byte, fr.prefixBuff.Available())
		_, err := fr.prefixBuff.Read(all)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(all))
		fr.prefixBuff.Reset()
		lBuff = all
	}

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
	err = json.Unmarshal(lBuff, newFrame)
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
