// Package filter provides an interface for filters and for elements that can pass through a filter.
//
// It contains a default implementation for the Savitzky-Golay filter.
package filter

// FilterFunction is a function that applies a filter to a slice of type FilterableElement.
type FilterFunction[TElement FilterableElement] func(Interface[TElement], []TElement) []TElement
type KeysFunction[TElement FilterableElement] func(f *Filter[TElement]) []string

type Filter[TElement FilterableElement] struct {
	FilterFunction FilterFunction[TElement]
	KeysFunction   func(f *Filter[TElement]) []string
	Elements       []TElement
	size           int
	full           bool
}

func (f *Filter[TElement]) Size() int {
	return f.size
}

func (f *Filter[TElement]) Keys() []string {
	return f.KeysFunction(f)
}

func (f *Filter[TElement]) Step(element TElement) (*TElement, FilterError) {
	f.Elements = append(f.Elements, element)
	if !f.full && len(f.Elements) < f.size {
		return nil, NotFullError{}
	}
	f.full = true

	f.Elements = f.FilterFunction(f, f.Elements)
	poppedElement := &(f.Elements[0])
	f.Elements = f.Elements[1:]
	return poppedElement, nil
}

func New[TElement FilterableElement](fFilter FilterFunction[TElement], fKeys KeysFunction[TElement]) Filter[TElement] {
	return Filter[TElement]{
		FilterFunction: fFilter,
		KeysFunction:   fKeys,
		Elements:       make([]TElement, 0),
		size:           0,
		full:           false,
	}
}
