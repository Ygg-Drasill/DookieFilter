package frameReader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"io"
	"log"
	"os"
)

type FrameReader struct {
	buff       *bufio.Reader
	file       *os.File
	prefixBuff bytes.Buffer
}

func New(path string) (*FrameReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	return &FrameReader{buff: bufio.NewReader(f), file: f}, nil
}

func (fr *FrameReader) Next() (*types.Frame[types.DataPlayer], error) {
	lBuff, isPrefix, err := fr.buff.ReadLine()
	if isPrefix {
		n, err := fr.prefixBuff.Write(lBuff)
		if err != nil {
			return nil, fmt.Errorf("writing prefix: %w", err)
		}
		log.Printf("buffering %d bytes as prefix", n)
		return fr.Next()
	}

	prefixSize := fr.prefixBuff.Len()
	if prefixSize > 0 {
		fr.prefixBuff.Write(lBuff)
		all := make([]byte, fr.prefixBuff.Len())
		_, err := fr.prefixBuff.Read(all)
		if err != nil {
			return nil, fmt.Errorf("reading prefix: %w", err)
		}

		fr.prefixBuff.Reset()
		lBuff = all
	}

	if err != nil {
		if err == io.EOF {
			cerr := fr.file.Close()
			if cerr != nil {
				return nil, fmt.Errorf("closing file: %w", cerr)
			}
			return nil, fmt.Errorf("end of file: %w", err)
		}
		return nil, fmt.Errorf("reading line: %w", err)
	}

	newFrame := new(types.Frame[types.DataPlayer])
	err = json.Unmarshal(lBuff, newFrame)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling frame: %w", err)
	}

	if len(newFrame.Data[0].Ball.Xyz) == 0 { //TODO: maybe return signal later :)
		fmt.Println("hello")
		return fr.Next()
	}

	if len(newFrame.Data[0].HomePlayers) == 0 {
		return fr.Next()
	}

	return newFrame, nil
}
