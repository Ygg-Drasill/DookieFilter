package pringleBuffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testElement struct {
	key  Key
	data string
}

func newTestElement(key Key) testElement {
	return testElement{
		key:  key,
		data: "",
	}
}

func (te testElement) Key() Key {
	return Key(te.key)
}

func Test_PringleBuffer_Ascending(t *testing.T) {
	aKey := Key(1)
	a := newTestElement(1)
	bKey := Key(2)
	b := newTestElement(2)
	cKey := Key(3)
	c := newTestElement(3)
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
	assert.Empty(t, element, "Expected empty since element a should not exist")
	assert.NotEqualf(t, a, queue.tail.data, "Element a should not still be tail")

	err = nil
	element, err = queue.Get(bKey)
	assert.NotNil(t, element, "Key b should exist")
	assert.Equal(t, b, element)

	element, err = queue.Get(cKey)
	assert.NotNil(t, element, "Key c should exist")
	assert.Equal(t, c, element)
}

func Test_PringleBuffer_Descending(t *testing.T) {
	queue := New[testElement](3)

	queue.Insert(newTestElement(5))
	got, err := queue.Get(5)
	assert.NotNil(t, got)
	assert.Nil(t, err)
	queue.Insert(newTestElement(4))
	got, err = queue.Get(4)
	assert.NotNil(t, got)
	assert.Nil(t, err)
	queue.Insert(newTestElement(3))
	assert.Equal(t, 3, queue.Count(), "Expected count of 3 after full")
	queue.Insert(newTestElement(2))
	assert.Equal(t, 3, queue.Count(), "Expected count of 3 after full")
	got, err = queue.Get(2)
	assert.NotNil(t, err, "Expected error since element 2 should not exist")
	assert.Empty(t, got, "Expected empty since element 2 should not exist")
	queue.Insert(newTestElement(1))
	assert.Equal(t, 3, queue.Count(), "Expected count of 3 after full")
	got, err = queue.Get(1)
	assert.NotNil(t, err, "Expected error since element 2 should not exist")
	assert.Empty(t, got, "Expected empty since element 1 should not exist")
}

func Test_PringleBuffer_Count(t *testing.T) {
	queue := New[testElement](4)

	queue.Insert(newTestElement(1))
	assert.Equal(t, 1, queue.Count(), "Expected count of 1 after first insertion")
	queue.Insert(newTestElement(3))
	assert.Equal(t, 2, queue.Count(), "Expected count of 2 after second insertion")
	queue.Insert(newTestElement(2))
	assert.Equal(t, 3, queue.Count(), "Expected count of 3 after third insertion")
	queue.Insert(newTestElement(5))
	assert.Equal(t, 4, queue.Count(), "Expected count of 4 after fourth")
	queue.Insert(newTestElement(4))
	assert.Equal(t, 4, queue.Count(), "Expected count of 4 after full")
}

type overwriteCase struct {
	key    Key
	before string
	after  string
}

func Test_PringleBuffer_Overwrite(t *testing.T) {
	var cases = []overwriteCase{
		{0, "a", "e"},
		{1, "b", "f"},
		{2, "c", "g"},
		{3, "d", "h"},
	}

	queue := New[testElement](4)
	for _, testCase := range cases {
		queue.Insert(testElement{key: testCase.key, data: testCase.before})
	}

	for _, testCase := range cases {
		e, err := queue.Get(testCase.key)
		assert.NoError(t, err)
		assert.Equal(t, testCase.before, e.data)
		queue.Insert(testElement{key: testCase.key, data: testCase.after})
		e, err = queue.Get(testCase.key)
		assert.NoError(t, err)
		assert.Equal(t, testCase.after, e.data)
	}

	for _, testCase := range cases {
		e, err := queue.Get(testCase.key)
		assert.NoError(t, err)
		assert.Equal(t, testCase.after, e.data)
	}
}
