package pringleBuffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testElement Key

func (te testElement) Key() Key {
	return Key(te)
}

func TestPringleBuffer(t *testing.T) {
	aKey := Key(1)
	a := testElement(aKey)
	bKey := Key(2)
	b := testElement(bKey)
	cKey := Key(3)
	c := testElement(cKey)
	size := 2

	queue := New[testElement](size)
	queue.Insert(a)
	assert.Equal(t, 1, queue.Count(), "Expected count of 1 after first insertion")
	queue.Insert(b)
	assert.Equal(t, 2, queue.Count(), "Expected count of 2 after second insertion")
	queue.Insert(c)
	assert.Equal(t, size, queue.Count(), "Expected count to still be 2 after third insertion")

	element, err := queue.Get(aKey)
	assert.NotNil(t, err, "Expected error since element a should not exist")
	assert.Nil(t, element, "Element a should not exist")
	assert.NotEqualf(t, a, queue.tail.data, "Element a should not still be tail")

	err = nil
	element, err = queue.Get(bKey)
	assert.NotNil(t, element, "Key b should exist")
	assert.Equal(t, b, element)

	element, err = queue.Get(cKey)
	assert.NotNil(t, element, "Key c should exist")
	assert.Equal(t, c, element)
}
