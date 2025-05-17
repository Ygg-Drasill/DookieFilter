package filter

type Interface[TElement FilterableElement] interface {
	// Step increments the data through the filter window,
	// if the filter is full, the last element in the filter will be returned.
	Step(TElement) (*TElement, FilterError)
	// Size returns the size/length of the filter window.
	Size() int
	// Keys returns the keys for the data series to be passed through the filter.
	Keys() []string
}

// FilterableElement is an element with one or more data points, to be passed through a filter.
type FilterableElement interface {
	// Update mutates one of the data points in the element.
	Update(key string, value float64) error
	// Get returns the value associated with the given key.
	Get(key string) (float64, error)
}
