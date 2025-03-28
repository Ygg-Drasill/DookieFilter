package pringleBuffer

import (
    "github.com/stretchr/testify/assert"
    "testing"
)

type testElement Key

func (te testElement) Key() Key {
    return Key(te)
}

func TestPringleBuffer_Ascending(t *testing.T) {
    aKey := Key(1)
    a := testElement(1)
    bKey := Key(2)
    b := testElement(2)
    cKey := Key(3)
    c := testElement(3)
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

func TestPringleBuffer_Descending(t *testing.T) {
    queue := New[testElement](3)

    queue.Insert(testElement(5))
    got, err := queue.Get(5)
    assert.NotNil(t, got)
    assert.Nil(t, err)
    queue.Insert(testElement(4))
    got, err = queue.Get(4)
    assert.NotNil(t, got)
    assert.Nil(t, err)
    queue.Insert(testElement(3))
    assert.Equal(t, 3, queue.Count(), "Expected count of 3 after full")
    queue.Insert(testElement(2))
    assert.Equal(t, 3, queue.Count(), "Expected count of 3 after full")
    got, err = queue.Get(2)
    assert.NotNil(t, err, "Expected error since element 2 should not exist")
    assert.Empty(t, got, "Expected empty since element 2 should not exist")
    queue.Insert(testElement(1))
    assert.Equal(t, 3, queue.Count(), "Expected count of 3 after full")
    got, err = queue.Get(1)
    assert.NotNil(t, err, "Expected error since element 2 should not exist")
    assert.Empty(t, got, "Expected empty since element 1 should not exist")
}

func TestPringleBuffer_Count(t *testing.T) {
    queue := New[testElement](4)

    queue.Insert(testElement(1))
    assert.Equal(t, 1, queue.Count(), "Expected count of 1 after first insertion")
    queue.Insert(testElement(3))
    assert.Equal(t, 2, queue.Count(), "Expected count of 2 after second insertion")
    queue.Insert(testElement(2))
    assert.Equal(t, 3, queue.Count(), "Expected count of 3 after third insertion")
    queue.Insert(testElement(5))
    assert.Equal(t, 4, queue.Count(), "Expected count of 4 after fourth")
    queue.Insert(testElement(4))
    assert.Equal(t, 4, queue.Count(), "Expected count of 4 after full")
}
