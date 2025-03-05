package frameReader

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"github.com/schollz/progressbar/v3"
	"io"
	"log"
	"log/slog"
	"os"
)

const StartByte = 100_000_000

type FrameReader struct {
	buff        *bufio.Reader
	file        *os.File
	prefixBuff  bytes.Buffer
	lineStarts  []int64
	lineCount   int64
	frameBuffer []types.Frame
}

func New(path string) (*FrameReader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	_, err = f.Seek(StartByte, 0)
	if err != nil {
		log.Fatal(err)
	}
	_, err = f.Seek(StartByte, 0)
	if err != nil {
		log.Fatal(err)
	}
	fr := &FrameReader{
		buff: bufio.NewReader(f),
		file: f,
	}
	fr.loadFrameBeginnings()
	fr.file.Seek(0, 0)
	slog.Debug(fmt.Sprintf("%d frames loaded", len(fr.lineStarts)))
	return fr, nil
}

func (fr *FrameReader) loadFrameBeginnings() {
	fr.lineStarts = make([]int64, 1)
	fr.file.Seek(0, 0)
	info, _ := fr.file.Stat()

	buff := make([]byte, 64*1024)
	filePosition := int64(0)
	bar := progressbar.DefaultBytes(int64(info.Size()), "loading frame positions")

	for {
		n, readErr := fr.file.Read(buff)
		if readErr == io.EOF || n == 0 {
			break
		} else if readErr != nil {
			log.Fatal(readErr)
		}
		bar.Add(n)

		for i := 0; i < n; i++ {
			if buff[i] == '\n' {
				fr.lineStarts = append(fr.lineStarts, filePosition+int64(i)+1)
			}
		}
		filePosition += int64(n)
	}
	fr.lineCount = int64(len(fr.lineStarts))
}

func (fr *FrameReader) goToNextFrameStart() {
	char := make([]byte, 1)
	for char[0] != '\n' {
		fr.file.Read(char)
	}
}

func (fr *FrameReader) Next() (*types.Frame, error) {
	if len(fr.frameBuffer) > 0 {
		frames := fr.frameBuffer[0]
		fr.frameBuffer = fr.frameBuffer[1:]
		return &frames, nil
	}
	lBuff, isPrefix, err := fr.buff.ReadLine()
	if isPrefix {
		_, err := fr.prefixBuff.Write(lBuff)
		if err != nil {
			return nil, fmt.Errorf("writing prefix: %w", err)
		}
		//log.Printf("buffering %d bytes as prefix", n)
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
			return nil, err
		}
		return nil, fmt.Errorf("reading line: %w", err)
	}

	newFrame := new(types.GamePacket[types.Frame])
	err = json.Unmarshal(lBuff, newFrame)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling frame: %w", err)
	}

	if len(newFrame.Data) == 0 {
		return fr.Next()
	}
	fr.frameBuffer = newFrame.Data[1:]

	return &newFrame.Data[0], nil
}

func (fr *FrameReader) GoToFrame(frameIndex int64) error {
	if frameIndex > fr.lineCount {
		return fmt.Errorf("index %d out of frame range", frameIndex)
	}
	_, err := fr.file.Seek(fr.lineStarts[frameIndex], 0)
	if err != nil {
		slog.Error("failed to seek in file", "error", err)
	}
	fr.buff.Reset(fr.file)
	return nil
}

func (fr *FrameReader) FrameCount() int64 {
	return fr.lineCount
}
