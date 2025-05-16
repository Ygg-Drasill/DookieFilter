package filter

func NewSavitzkyGolayFilter[TElement FilterableElement](length int, order int) Filter[TElement] {
	return Filter[TElement]{
		FilterFunction: func(f Interface[TElement], elements []TElement, key string) []TElement {
			updateIndex := f.Size() / 2
			sum := 0.0
			for i := range f.Size() {
				val, err := elements[i].Get(key)
				if err != nil {
					return elements
				}
				sum += val
				if err != nil {
					return elements
				}
			}
			err := elements[updateIndex].Update(key, sum/float64(f.Size()))
			if err != nil {
				return elements
			}
			return elements
		},
		full: false,
		size: length,
	}
}
