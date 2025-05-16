package filter

type Interface[TElement FilterableElement] interface {
	Step(TElement) (*TElement, FilterError)
	Size() int
}

type FilterableElement interface {
	Keys() []string
	Update(key string, value float64) error
	Get(key string) (float64, error)
}
