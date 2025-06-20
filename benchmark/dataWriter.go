package main

import (
	"encoding/json"
	"os"

	"github.com/Ygg-Drasill/DookieFilter/common/types"
)

type dataWriter struct {
	file       *os.File
}	

func newDataWriter(path string) *dataWriter {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	dw := &dataWriter{
		file: f,
	}

	return dw
}

func (dw *dataWriter) writeFrame(sFrame types.SmallFrame) {
	frame := types.Frame{
		HomePlayers: []types.Player{},
		AwayPlayers: []types.Player{},
		Ball: struct {
			Speed float64 "json:\"speed\""
			Xyz []float64 "json:\"xyz\""
		} {
			Speed: 0.0,
			Xyz: []float64{sFrame.Ball.X, sFrame.Ball.Y, 0.0},
		},
		FrameIdx: sFrame.FrameIdx,
	}
	bytes, err := json.Marshal(frame)
	if err != nil {
		panic(err)
	}
	dw.file.Write(bytes)
	n, err := dw.file.WriteString(",\n")
	if n == 0 {
		panic("failed to write newline")
	}
	if err != nil {
		panic(err)


	}
}
