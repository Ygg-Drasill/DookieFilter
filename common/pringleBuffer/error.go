package pringleBuffer

type PringleBufferError struct {
	msg string
}

func (e PringleBufferError) Error() string {
	return e.msg
}
