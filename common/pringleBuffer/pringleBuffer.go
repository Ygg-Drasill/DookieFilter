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

	prev := pb.head
	next := pb.head.next
	//traverse
	for prev.Key() > newElement.Key() {
		prev, next = next, next.next
	}

	//insert not full
	if !full {
		pb.count++
		if prev == pb.head { //is new head
			newElement.next = pb.head
			newElement.next.prev = newElement
			pb.head = newElement
			return
		}

		if next == nil { //is new tail
			prev.next = newElement
			newElement.prev = prev
			pb.tail = newElement
			return
		}

		prev.next = newElement
		next.prev = newElement
		newElement.next = next
		newElement.prev = prev
		return
	}

	//if full - only insert into middle -> trim tail
	if next == nil { //would be new tail, but buffer is full so don't insert
		return
	}

	if prev == pb.head { //is new head
		newElement.next = pb.head
		newElement.next.prev = newElement
		pb.head = newElement
	} else {
		prev.next = newElement
		next.prev = newElement
		newElement.next = next
		newElement.prev = prev
	}

	//trim tail
	tail := pb.tail
	pb.tail = tail.prev
	pb.tail.next = nil
}

func (pb *PringleBuffer[TElement]) Get(key Key) (PringleIndexable, error) {
	var element *PringleElement[TElement]
	current := pb.head
	for current.Key() != key {
		current = current.next
		if current == nil {
			return nil, PringleBufferError{msg: "element does not exist"}
		}
	}
	element = current
	return element.data, nil
}
