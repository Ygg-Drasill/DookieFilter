package main

import (
	"encoding/csv"
	// "fmt" // Import fmt for formatted printing - Removed as it's unused
	"github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"log"
	"os"
	"strconv"
)

// MockWorker is a simplified version of the detector worker for testing
type MockWorker struct {
	stateBuffer        *pringleBuffer.PringleBuffer[types.SmallFrame]
	missingPlayers     map[string]int // playerID -> startFrameIdx of missing period
	lastProcessedFrame *types.SmallFrame
	holeCount          int // Counter for total holes detected
}

func NewMockWorker() *MockWorker {
	return &MockWorker{
		stateBuffer:        pringleBuffer.New[types.SmallFrame](10), // Keep buffer for potential future use
		missingPlayers:     make(map[string]int),
		lastProcessedFrame: nil,
		holeCount:          0,
	}
}

func (w *MockWorker) detectHoles(currentFrame types.SmallFrame) {
	if w.lastProcessedFrame == nil {
		// Cannot compare with a previous frame yet
		return
	}

	prevFrame := *w.lastProcessedFrame

	// Create sets for efficient lookup
	currentPlayers := make(map[string]bool)
	for _, player := range currentFrame.Players {
		currentPlayers[player.PlayerId] = true
	}

	prevPlayers := make(map[string]bool)
	for _, player := range prevFrame.Players {
		prevPlayers[player.PlayerId] = true
	}

	// Check for players who were present before but are missing now
	for playerID := range prevPlayers {
		if !currentPlayers[playerID] {
			// Player is missing in the current frame
			if _, alreadyMissing := w.missingPlayers[playerID]; !alreadyMissing {
				// Player just went missing, record start frame
				w.missingPlayers[playerID] = currentFrame.FrameIdx
				// log.Printf("Debug: Player %s started missing at frame %d", playerID, currentFrame.FrameIdx) // Optional debug log
			}
		}
	}

	// Check for players who were missing but have returned
	for playerID := range currentPlayers {
		if startFrame, wasMissing := w.missingPlayers[playerID]; wasMissing {
			// Player was missing and has now returned
			endFrame := currentFrame.FrameIdx - 1
			if startFrame <= endFrame { // Ensure start is not after end (can happen if missing for only one frame instant)
				log.Printf("Hole detected: Player %s missing from frame %d to %d",
					playerID, startFrame, endFrame)
				w.holeCount++ // Increment hole count
			} else {
				// This case (startFrame > endFrame) implies the player reappeared immediately in the next frame.
				// Depending on requirements, you might log this differently or ignore it.
				// For now, we can log it as a brief disappearance or skip logging.
				// log.Printf("Info: Player %s briefly disappeared and reappeared between frame %d and %d", playerID, prevFrame.FrameIdx, currentFrame.FrameIdx)
			}
			delete(w.missingPlayers, playerID) // Remove from missing map as they've returned
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
	frames, err := readHoleCSV("../../../gym/data/chunk_0.csv")
	if err != nil {
		log.Fatalf("Error reading chunk_0.csv: %v", err)
	}

	if len(frames) == 0 {
		log.Fatal("No frames were created from the CSV file")
	}

	log.Printf("Processing %d frames...", len(frames))
	lastFrameIdx := 0
	// Process frames
	for _, frame := range frames {
		// We might still want the buffer for other potential analyses, so let's insert.
		worker.stateBuffer.Insert(frame)
		worker.detectHoles(frame)
		// Update the last processed frame *after* detecting holes based on the previous one
		worker.lastProcessedFrame = &frame
		if frame.FrameIdx > lastFrameIdx {
			lastFrameIdx = frame.FrameIdx
		}
	}

	// After processing all frames, check for players still missing
	if len(worker.missingPlayers) > 0 {
		log.Println("Players still missing at the end of the data:")
		for playerID, startFrame := range worker.missingPlayers {
			log.Printf("  - Player %s missing from frame %d to end (last frame %d)",
				playerID, startFrame, lastFrameIdx)
			worker.holeCount++ // Increment hole count for players missing till the end
		}
	}

	log.Println("Test completed")
	log.Printf("Total holes detected: %d", worker.holeCount)
}
