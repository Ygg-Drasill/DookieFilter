package main

import (
	"encoding/csv"
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"log"
	"os"
	"strconv"
)

// MockWorker is a simplified version of the detector worker for testing
type MockWorker struct {
	stateBuffer *pringleBuffer.PringleBuffer[types.SmallFrame]
}

func NewMockWorker() *MockWorker {
	return &MockWorker{
		stateBuffer: pringleBuffer.New[types.SmallFrame](10),
	}
}

func (w *MockWorker) detectHoles(frame types.SmallFrame) {
	// Get previous frames from the buffer
	prevFrames := make([]types.SmallFrame, 0)
	for i := 1; i <= 10; i++ {
		prevFrame, err := w.stateBuffer.Get(pringleBuffer.Key(frame.FrameIdx - i))
		if err == nil {
			prevFrames = append(prevFrames, prevFrame)
		}
	}

	if len(prevFrames) == 0 {
		return
	}

	// Track players that appear in current frame
	currentPlayers := make(map[string]bool)
	for _, player := range frame.Players {
		currentPlayers[player.PlayerId] = true
	}

	// Check each previous frame for missing players
	for _, prevFrame := range prevFrames {
		for _, player := range prevFrame.Players {
			if !currentPlayers[player.PlayerId] {
				log.Printf("Hole detected! Player %s missing in frame %d (last seen in frame %d)",
					player.PlayerId, frame.FrameIdx, prevFrame.FrameIdx)
			}
		}
	}
}

func readHoleCSV(filePath string) ([]types.SmallFrame, error) {
	log.Printf("Attempting to read file: %s", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	log.Printf("Read %d records from CSV", len(records))
	if len(records) == 0 {
		return nil, nil
	}

	// Get header row to identify columns
	headers := records[0]
	log.Printf("CSV Headers: %v", headers)
	
	// Find frame_index column
	frameIdxCol := -1
	for i, header := range headers {
		if header == "frame_index" {
			frameIdxCol = i
			break
		}
	}
	if frameIdxCol == -1 {
		return nil, nil
	}

	playerColumns := make(map[string]struct{})
	for _, header := range headers {
		if header == "frame_index" {
			continue
		}
		// Extract player ID from column name (e.g., "h_10_x" -> "h_10")
		playerID := header[:len(header)-2] // Remove "_x" or "_y"
		playerColumns[playerID] = struct{}{}
	}

	log.Printf("Found player columns: %v", playerColumns)

	frames := make([]types.SmallFrame, 0)
	for _, record := range records[1:] { // Skip header row
		// Parse frame index from CSV
		frameIdx, err := strconv.Atoi(record[frameIdxCol])
		if err != nil {
			log.Printf("Error parsing frame index: %v", err)
			continue
		}

		frame := types.SmallFrame{
			FrameIdx: frameIdx,
			Players:  make([]types.PlayerPosition, 0),
		}

		// Process each player's position
		for playerID := range playerColumns {
			xCol := playerID + "_x"
			yCol := playerID + "_y"
			
			xIdx := -1
			yIdx := -1
			for j, header := range headers {
				if header == xCol {
					xIdx = j
				}
				if header == yCol {
					yIdx = j
				}
			}

			if xIdx != -1 && yIdx != -1 {
				x, err := strconv.ParseFloat(record[xIdx], 64)
				if err != nil {
					log.Printf("Error parsing x coordinate for player %s in frame %d: %v", playerID, frameIdx, err)
					continue
				}
				y, err := strconv.ParseFloat(record[yIdx], 64)
				if err != nil {
					log.Printf("Error parsing y coordinate for player %s in frame %d: %v", playerID, frameIdx, err)
					continue
				}
				
				frame.Players = append(frame.Players, types.PlayerPosition{
					PlayerId: playerID,
					Position: types.Position{X: x, Y: y},
				})
			}
		}

		if len(frame.Players) > 0 {
			frames = append(frames, frame)
		}
	}

	log.Printf("Successfully created %d frames", len(frames))
	return frames, nil
}

func main() {
	log.Println("Starting hole detection test...")
	worker := NewMockWorker()
	
	// Read frames from hole.csv
	frames, err := readHoleCSV("../../../gym/data/hole.csv")
	if err != nil {
		log.Fatalf("Error reading hole.csv: %v", err)
	}
	
	if len(frames) == 0 {
		log.Fatal("No frames were created from the CSV file")
	}

	log.Printf("Processing %d frames...", len(frames))
	// Process frames
	for _, frame := range frames {
		worker.stateBuffer.Insert(frame)
		worker.detectHoles(frame)
	}
	log.Println("Test completed")
} 