package main

import (
	"fmt"
)

func main() {

	// Simulated noisy soccer player trajectory (X, Y coordinates)
	xData := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	yData := []float64{0, 1.2, 1.8, 3.1, 3.9, 5.1, 6.2, 7.1, 8.5, 9.9, 11.0}

	// Apply Savitzky-Golay filter to smooth the trajectory
	smoothX := savitzkyGolayFilter(xData)
	smoothY := savitzkyGolayFilter(yData)

	// Printed noisy soccer player trajectory (X, Y coordinates)
	fmt.Println("Noisy X:", xData)
	fmt.Println("Noisy Y:", yData)

	// Print smoothed values
	fmt.Println("Smoothed X:", smoothX)
	fmt.Println("Smoothed Y:", smoothY)

	//Enables the display of the graph
	//http.HandleFunc("/", httpserver)
	//http.ListenAndServe(":8081", nil)

}
