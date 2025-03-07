package filter

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// Position struct represents a player's position
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
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
	var PositionsArray []float64
	//var yPositionsArray []float64
	for _, pos := range data.Positions {
		PositionsArray = append(PositionsArray, pos.X, pos.Y, pos.Z)

	}

	//var sArray = PlayerPosition{xPositionsArray, yPositionsArray}

	return PositionsArray
}

func TestFilter(t *testing.T) {

	//Reads data from json file
	data := jsonReader("soccer_player_positions.json")

	// Apply Savitzky-Golay filter to smooth the trajectory
	smoothData := savitzkyGolayFilter(data)

	// labels for the graph files
	var labelSmooth = "PlayerPosition/player_positions_smooth.png"
	var labelUnsmooth = "PlayerPosition/player_positions_unsmooth.png"

	// Used to visualize the smoothed data.
	ShowGraph(smoothData, labelSmooth)
	ShowGraph(data, labelUnsmooth)

}
