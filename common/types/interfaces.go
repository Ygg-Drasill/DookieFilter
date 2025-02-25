package types

type FrameLoader[TData DataPlayer | DataSignal] interface {
	Next() *Frame[TData]
}
