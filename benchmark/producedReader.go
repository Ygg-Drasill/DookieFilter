package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type producedFrameReader struct {
	buff       *bufio.Reader
	file       *os.File
	prefixBuff bytes.Buffer
}

func New(path string) (*producedFrameReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	return &producedFrameReader{
		buff: bufio.NewReader(f),
		file: f,
	}, nil
}

func (reader *producedFrameReader) next() *producedFrame {
	line, isPrefix, err := reader.buff.ReadLine()
	if err != nil {
		return nil
	}

	if !isPrefix {
		frame := producedFrame{}
		err = json.Unmarshal(line, &frame)
		if err != nil {
			return nil
		}

		reader.prefixBuff.Reset()
		return &frame
	}

	reader.prefixBuff.Write(line)
	return reader.next()
}
