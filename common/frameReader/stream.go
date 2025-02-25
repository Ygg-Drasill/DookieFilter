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
	"os"
	"sync"
)

const StartByte = 100_000_000

type FrameReader struct {
	buff        *bufio.Reader
	file        *os.File
	prefixBuff  bytes.Buffer
	frameStarts []int64
}

func New(path string) *FrameReader {
	f, err := os.Open(path)
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
	fr.goToNextFrameStart()
	fmt.Println(len(fr.frameStarts))
	return fr
}

const FrameLoaders = 2

func (fr *FrameReader) loadFrameBeginnings() {
	fr.frameStarts = make([]int64, 1)
	fr.frameStarts = append(fr.frameStarts, 0)
	info, err := fr.file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	chunkSize := info.Size() / FrameLoaders
	wg := sync.WaitGroup{}
	mu := sync.Mutex{}
	for i := 0; i < FrameLoaders; i++ {
		info, _ := fr.file.Stat()
		bar := progressbar.Default(info.Size())
		wg.Add(1)
		go func(offset, chunkSize int64) {
			f, err := os.Open(fr.file.Name())
			defer f.Close()
			if err != nil {
				log.Fatal(err)
			}

			buff := make([]byte, chunkSize)
			n, err := f.Read(buff)
			if err != nil || n == 0 {
				return
			}
			for k := int64(0); k < chunkSize; k += 1 {
				err = bar.Add(1)
				if err != nil {
					return
				}
				if buff[k] == '\n' {
					mu.Lock()
					fr.frameStarts = append(fr.frameStarts, k+1)
					mu.Unlock()
				}
			}
			wg.Done()
		}(int64(i)*chunkSize, chunkSize)
	}
	wg.Wait()
}

func (fr *FrameReader) goToNextFrameStart() {
	char := make([]byte, 1)
	for char[0] != '\n' {
		fr.file.Read(char)
	}
}

func (fr *FrameReader) Next() *types.Frame[types.DataPlayer] {
	lBuff, isPrefix, err := fr.buff.ReadLine()
	if isPrefix {
		_, err := fr.prefixBuff.Write(lBuff)
		if err != nil {
			log.Fatal(err)
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
			log.Fatal(err)
		}

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
