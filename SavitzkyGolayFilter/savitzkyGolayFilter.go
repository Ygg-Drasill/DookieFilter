package main

import "log"

// Apply Savitzky-Golay filter with a fixed 5-point quadratic window
func savitzkyGolayFilter(data []float64) []float64 {
	n := len(data)
	if n < 5 {
		log.Fatal("Data size must be at least 5 for this filter")
	}

	// Precomputed coefficients for a 5-point quadratic window
	coeffs := []float64{-3, 12, 17, 12, -3}
	divisor := 35.0 // Sum of positive coefficients

	smoothed := make([]float64, n)

	// Apply filter to the middle of the data
	for i := 2; i < n-2; i++ {
		smoothed[i] = (coeffs[0]*data[i-2] +
			coeffs[1]*data[i-1] +
			coeffs[2]*data[i] +
			coeffs[3]*data[i+1] +
			coeffs[4]*data[i+2]) / divisor
	}

	// Copy the edges without modification
	smoothed[0] = data[0]
	smoothed[1] = data[1]
	smoothed[n-1] = data[n-1]
	smoothed[n-2] = data[n-2]

	return smoothed
}
