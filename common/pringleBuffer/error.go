package pringleBuffer

import "fmt"

type PringleBufferError struct {
	msg string
}

func (e PringleBufferError) Error() string {
	return e.msg
}

type EmptyError struct {
}

func (e EmptyError) Error() string {
	return "Buffer is empty"
}

type NotFoundError struct {
	key Key
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("Key %v not found", e.key)
}
