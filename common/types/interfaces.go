package types

type FrameLoader[TData Frame | Signal] interface {
	Next() (*TData, error)
	FrameCount() int64
	GoToFrame(int64) error
}
