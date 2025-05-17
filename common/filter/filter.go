// Package filter provides an interface for filters and for elements that can pass through a filter.
//
// It contains a default implementation for the Savitzky-Golay filter.
package filter

// FilterFunction is a function that applies a filter to a slice of type FilterableElement.
type FilterFunction[TElement FilterableElement] func(Interface[TElement], []TElement) []TElement

type filter[TElement FilterableElement] struct {
	FilterFunction FilterFunction[TElement]
	Elements       []TElement
	size           int
	full           bool
}

func (f *filter[TElement]) Keys() []string {
	return []string{"x"}
}

func (f *filter[TElement]) Size() int {
	return f.size
}

func (f *filter[TElement]) Step(element TElement) (*TElement, FilterError) {
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
