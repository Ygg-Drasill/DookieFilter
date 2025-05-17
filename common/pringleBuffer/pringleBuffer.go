package pringleBuffer

type Key int

type PringleIndexable interface {
	Key() Key
}

// PringleBuffer
// A ring buffer, where elements are sorted (prioritized),
// descending from head to tail
type PringleBuffer[TElement PringleIndexable] struct {
	head  *PringleElement[TElement]
	tail  *PringleElement[TElement]
	Size  int
	count int
}

func New[TElement PringleIndexable](size int) *PringleBuffer[TElement] {
	return &PringleBuffer[TElement]{
		Size:  size,
		count: 0,
	}
}

func (pb *PringleBuffer[TElement]) Count() int {
	return pb.count
}

func (pb *PringleBuffer[TElement]) Insert(data TElement) {
	newElement := &PringleElement[TElement]{data: data}
	if pb.head == nil {
		pb.head = newElement
		pb.tail = newElement
		pb.count++
		return
	}
	if pb.tail == nil {
		pb.tail = newElement
	}
	full := pb.Count() >= pb.Size

	var prev *PringleElement[TElement]
	var next = pb.head
	//traverse
	for next != nil && next.Key() > newElement.Key() {
		prev, next = next, next.next
	}

	if next != nil && next.Key() == newElement.Key() {
		next.data = data
		return
	}

	//insert not full
	if !full {
		pb.count++
		pb.insertBetween(prev, next, newElement)
		return
	}

	//if full - only insert into middle -> trim tail
	if next == nil { //would be new tail, but buffer is full so don't insert
		return
	}

	pb.insertBetween(prev, next, newElement)
	pb.trimTail()
}

func (pb *PringleBuffer[TElement]) Get(key Key) (TElement, error) {
	var empty TElement
	var element *PringleElement[TElement]
	current := pb.head
	if current == nil {
		return empty, EmptyError{}
	}
	for current.Key() != key {
		current = current.next
		if current == nil {
			return empty, NotFoundError{key: key}
		}
	}
	element = current
	return element.data, nil
}

func (pb *PringleBuffer[TElement]) insertBetween(prev, next, element *PringleElement[TElement]) {
	if prev == nil { //is new head
		element.next = pb.head
		pb.head.prev = element
		pb.head = element
		return
	}

	if next == nil { //is new tail
		element.prev = pb.tail
		pb.tail.next = element
		pb.tail = element
		return
	}

	prev.next = element
	element.next = next
	next.prev = element
	element.prev = prev
}

func (pb *PringleBuffer[TElement]) trimTail() {
	tail := pb.tail
	pb.tail = tail.prev
	pb.tail.next = nil
}
