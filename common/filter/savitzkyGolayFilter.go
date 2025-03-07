package filter

import "log"

type FilterConfig struct {
	Datasise int
	Coeffs   []float64
	divise   float64
}

type PlayerPosition struct {
	xData []float64
	yData []float64
	zData []float64
}

func Filtersettings() *FilterConfig {

	var filterData = new(FilterConfig)

	filterData.Datasise = 5

	// Precomputed coefficients for a 5-point quadratic window
	filterData.Coeffs = []float64{-3, 12, 17, 12, -3}

	filterData.divise = 35.0 // Sum of positive coefficients

	return filterData
}

func splitIntoThree(input []float64) ([]float64, []float64, []float64) {
	var first, second, third []float64

	for i, v := range input {
		switch i % 3 {
		case 0:
			first = append(first, v)
		case 1:
			second = append(second, v)
		case 2:
			third = append(third, v)
		}
	}

	return first, second, third
}

func mergeThreeSlices(a, b, c []float64) []float64 {
	var result []float64
	length := max(len(a), len(b), len(c))

	for i := 0; i < length; i++ {
		if i < len(a) {
			result = append(result, a[i])
		}
		if i < len(b) {
			result = append(result, b[i])
		}
		if i < len(c) {
			result = append(result, c[i])
		}
	}

	return result
}

func max(x, y, z int) int {
	if x > y {
		if x > z {
			return x
		}
		return z
	}
	if y > z {
		return y
	}
	return z
}

// Apply Savitzky-Golay filter with a fixed 5-point quadratic window
func savitzkyGolayFilter(data []float64) []float64 {

	var DataPosition = PlayerPosition{}

	DataPosition.xData, DataPosition.yData, DataPosition.zData = splitIntoThree(data)

	filterData := Filtersettings()

	n := len(DataPosition.xData)
	m := len(DataPosition.yData)
	if n < filterData.Datasise && m < filterData.Datasise {
		log.Fatal("Data size must be at least 5 for this filter")
	}

	smoothedXdata := make([]float64, n)
	smoothedYdata := make([]float64, n)

	//This is for the x coordinates
	// Apply filter to the middle of the data
	for i := 2; i < n-2; i++ {
		smoothedXdata[i] = (filterData.Coeffs[0]*DataPosition.xData[i-2] +
			filterData.Coeffs[1]*DataPosition.xData[i-1] +
			filterData.Coeffs[2]*DataPosition.xData[i] +
			filterData.Coeffs[3]*DataPosition.xData[i+1] +
			filterData.Coeffs[4]*DataPosition.xData[i+2]) / filterData.divise
	}
	// Copy the edges without modification
	smoothedXdata[0] = DataPosition.xData[0]
	smoothedXdata[1] = DataPosition.xData[1]
	smoothedXdata[n-1] = DataPosition.xData[n-1]
	smoothedXdata[n-2] = DataPosition.xData[n-2]

	//This is for the y coordinates
	for i := 2; i < n-2; i++ {
		smoothedYdata[i] = (filterData.Coeffs[0]*DataPosition.yData[i-2] +
			filterData.Coeffs[1]*DataPosition.yData[i-1] +
			filterData.Coeffs[2]*DataPosition.yData[i] +
			filterData.Coeffs[3]*DataPosition.yData[i+1] +
			filterData.Coeffs[4]*DataPosition.yData[i+2]) / filterData.divise
	}
	// Copy the edges without modification
	smoothedYdata[0] = DataPosition.yData[0]
	smoothedYdata[1] = DataPosition.yData[1]
	smoothedYdata[n-1] = DataPosition.yData[n-1]
	smoothedYdata[n-2] = DataPosition.yData[n-2]

	smoothedCoordinates := mergeThreeSlices(smoothedXdata, smoothedYdata, DataPosition.zData)

	return smoothedCoordinates
}
