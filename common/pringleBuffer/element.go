package pringleBuffer

type PringleElement[TElement PringleIndexable] struct {
	data TElement
	next *PringleElement[TElement]
	prev *PringleElement[TElement]
}

func (pe PringleElement[TElement]) Key() Key {
	return pe.data.Key()
}
