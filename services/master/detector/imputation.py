package detector

import (
    "github.com/Ygg-Drasill/DookieFilter/common/frameReader"
    zmq "github.com/pebbe/zmq4"
)

type Imputation struct {
    frameReader.FrameReader
}

func NewImputation(ctx *zmq.Context) *Imputation {
    return &Imputation{
        FrameReader: frameReader.NewFrameReader(ctx),
    }
}
func (i *Imputation) Save(frame *frameReader.Frame) {
    i.FrameReader.Save(frame)
}

