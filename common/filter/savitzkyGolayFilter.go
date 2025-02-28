package filter

import "log"

type FilterConfig struct {
	Datasise int
	Coeffs   []float64
	divise   float64
}

func Filtersettings() *FilterConfig {

	var filterData = new(FilterConfig)

	filterData.Datasise = 5

	// Precomputed coefficients for a 5-point quadratic window
	filterData.Coeffs = []float64{-3, 12, 17, 12, -3}

	filterData.divise = 35.0 // Sum of positive coefficients

	return filterData
}

// Apply Savitzky-Golay filter with a fixed 5-point quadratic window
func savitzkyGolayFilter(data []float64) []float64 {

	filterData := Filtersettings()

	n := len(data)
	if n < filterData.Datasise {
		log.Fatal("Data size must be at least 5 for this filter")
	}

	smoothed := make([]float64, n)

	// Apply filter to the middle of the data
	for i := 2; i < n-2; i++ {
		smoothed[i] = (filterData.Coeffs[0]*data[i-2] +
			filterData.Coeffs[1]*data[i-1] +
			filterData.Coeffs[2]*data[i] +
			filterData.Coeffs[3]*data[i+1] +
			filterData.Coeffs[4]*data[i+2]) / filterData.divise
	}

	// Copy the edges without modification
	smoothed[0] = data[0]
	smoothed[1] = data[1]
	smoothed[n-1] = data[n-1]
	smoothed[n-2] = data[n-2]

	return smoothed
}
