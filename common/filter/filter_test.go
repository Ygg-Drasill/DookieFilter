package filter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// Position struct represents a player's position
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// Data struct represents the structure of the JSON file
type Data struct {
	Positions []Position `json:"positions"`
}

func jsonReader(jsonfile string) []float64 {
	// Open the JSON file
	file, err := os.Open(jsonfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file content
	byteValue, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	// Parse JSON data
	var data Data
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		log.Fatal(err)
	}

	// Convert to a float64 slice
	var positionsArray []float64
	for _, pos := range data.Positions {
		positionsArray = append(positionsArray, pos.X, pos.Y)
	}

	// Print the result
	fmt.Println("Test data: ", positionsArray)
	return positionsArray
}

func TestFilter(t *testing.T) {
	// Simulated noisy soccer player trajectory (X, Y coordinates)
	xData := []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	yData := []float64{0, 1.2, 1.8, 3.1, 3.9, 5.1, 6.2, 7.1, 8.5, 9.9, 11.0}

	data := jsonReader("soccer_player_positions.json")

	fmt.Println("Test data: ", data)

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
