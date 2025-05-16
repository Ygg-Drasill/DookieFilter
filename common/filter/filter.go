package filter

type FilterFunction[TElement FilterableElement] func(Interface[TElement], []TElement, string) []TElement

type Filter[TElement FilterableElement] struct {
	FilterFunction FilterFunction[TElement]
	Elements       []TElement
	size           int
	full           bool
}

func (f *Filter[TElement]) Size() int {
	return f.size
}

func (f *Filter[TElement]) Step(element TElement) (*TElement, FilterError) {
	f.Elements = append(f.Elements, element)
	if !f.full && len(f.Elements) < f.size {
		return nil, NotFullError{}
	}
	f.full = true

	f.Elements = f.FilterFunction(f, f.Elements, "x")
	poppedElement := &(f.Elements[0])
	f.Elements = f.Elements[1:]
	return poppedElement, nil
}
