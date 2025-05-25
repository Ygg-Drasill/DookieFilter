package pringleBuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_onPopTail(t *testing.T) {
	callbackRan := false
	buffer := New(3, func(buffer *PringleBuffer[testElement]) {
		callbackRan = true
	})

	buffer.Insert(newTestElement(1))
	buffer.Insert(newTestElement(2))
	buffer.Insert(newTestElement(3))

	buffer.Insert(newTestElement(4))

	assert.Equal(t, true, callbackRan)
}
