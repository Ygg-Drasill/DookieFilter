package types

import "sync"

type FrameLoader[TData DataPlayer | DataSignal] interface {
	Next() (*Frame[TData], error)
	FrameCount() int64
	GoToFrame(int64) error
}

type Worker interface {
	Run(wg *sync.WaitGroup)
}
