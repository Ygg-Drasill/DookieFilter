package pringleBuffer

type Key int

type PringleIndexable interface {
	Key() Key
}

type PringleBuffer[TElement PringleIndexable] struct {
	head  *PringleElement[TElement]
	tail  *PringleElement[TElement]
	Size  int
	count int
	dirty bool
}

func New[TElement PringleIndexable](size int) *PringleBuffer[TElement] {
	return &PringleBuffer[TElement]{
		Size:  size,
		count: 0,
	}
}

func (pb *PringleBuffer[TElement]) Count() int {
	if !pb.dirty {
		return pb.count
	}

	count := 0
	element := pb.head
	for element.next != nil {
		count++
	}
	pb.count = count
	pb.dirty = false
	return count
}

func (pb *PringleBuffer[TElement]) Insert(data TElement) {
	pb.dirty = true
	newElement := &PringleElement[TElement]{
		data: data,
	}
	if pb.head == nil {
		pb.head = newElement
		pb.count++
		return
	}

	next := pb.head
	var prev *PringleElement[TElement]
	for next.Key() > newElement.Key() {
		prev = next
		next = next.next
	}

	full := pb.Count() < pb.Size

	if !full && prev != nil {
		prev.next = newElement
	}

	if next == nil && !full {
		pb.tail = newElement
	} else if !full {
		newElement.next = next
	}

	pb.count++
}

func (pb *PringleBuffer[TElement]) Get(key Key) (PringleIndexable, error) {
	var element *PringleElement[TElement]
	current := pb.head
	for current.Key() > key {
		next := current.next
		if next.Key() == key {
			return next, nil
		}
		if next.Key() < key {
			break
		}
		current = next
	}
	if element == nil {
		return nil, PringleBufferError{msg: "element does not exist"}
	}
	return element.data, nil
}
