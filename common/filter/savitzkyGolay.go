package filter

import "log/slog"

const divisor = 35.0
const length = 5

// Order of 2
var coefficients = []float64{-3, 12, 17, 12, -3}

func SavGolFilter[TElement FilterableElement]() FilterFunction[TElement] {
	return func(f Interface[TElement], elements []TElement) []TElement {
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
				slog.Error("Error updating filtered element", "error", err)
				return elements
			}
			//End of filter logic

		}

		return elements
	}
}
