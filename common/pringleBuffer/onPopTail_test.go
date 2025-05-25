package pringleBuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_onPopTail(t *testing.T) {
	callbackRan := false
	buffer := New[testElement](3, WithOnPopTail(func(element testElement) {
		callbackRan = true
	}))
	assert.Equal(t, false, callbackRan, "Should not be called from initialization")
	buffer.Insert(newTestElement(1))
	assert.Equal(t, false, callbackRan, "Should not be called when not full yet")
	buffer.Insert(newTestElement(2))
	assert.Equal(t, false, callbackRan, "Should not be called when not full yet")
	buffer.Insert(newTestElement(3))
	assert.Equal(t, false, callbackRan, "Should not be called on last insertion before full")

	//Act
	buffer.Insert(newTestElement(4))
	assert.Equal(t, true, callbackRan, "Should be called on insertion when buffer is full ")
}
