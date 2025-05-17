package filter

const divisor = 35.0

var coefficients = []float64{-3, 12, 17, 12, -3}

func NewSavitzkyGolayFilter[TElement FilterableElement](length int, order int) filter[TElement] {
	return filter[TElement]{
		FilterFunction: func(f Interface[TElement], elements []TElement) []TElement {
			keys := f.Keys()
			updateIndex := f.Size() / 2

			for _, k := range keys {

				//Per series filter logic
				sum := 0.0
				for i := range f.Size() {
					rawValue, err := elements[i].Get(k)
					if err != nil {
						return elements
					}
					sum += coefficients[i] * rawValue
				}
				filteredValue := sum / divisor
				err := elements[updateIndex].Update(k, filteredValue)
				if err != nil {
					return elements
				}
				//End of filter logic

			}

			return elements
		},
		full: false,
		size: length,
	}
}
