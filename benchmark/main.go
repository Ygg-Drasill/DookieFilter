package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

const sendInterval = time.Second / 25
const testLength = 25 * 60 * 10 //1

func main() {
	mutex := sync.Mutex{}
	timeStartMap := make(map[int]time.Time)

	rawFilePath, producedFilePath := os.Args[1], os.Args[2]
	rawStat, err := os.Stat(rawFilePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File %s does not exist\n", rawFilePath)
	}

	producedStat, err := os.Stat(producedFilePath)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File %s does not exist\n", producedFilePath)
	}

	fmt.Printf("Running benchmark on rawFilePath:%s produced:%s\n", rawStat.Name(), producedStat.Name())

	rawReader, _ := frameReader.New(rawFilePath)
	producedReader, _ := New(producedFilePath)

	atGameStart := false
	var startFrame *types.Frame
	for !atGameStart {
		startFrame, err = rawReader.Next()
		if err != nil {
			panic(err)
		}
		atGameStart = startFrame.GameClock != 0
	}

	pFrame := &producedFrame{}
	for {
		pFrame = producedReader.next()
		if pFrame.GameClock == startFrame.GameClock {
			break
		}
	}

	startFrame.FrameIdx = pFrame.FrameIdx

	ctx, err := zmq.NewContext()
	if err != nil {
		panic(err)
	}
	socketInput, err := ctx.NewSocket(zmq.PUSH)
	if err != nil {
		panic(err)
	}
	socketOutput, err := ctx.NewSocket(zmq.PULL)
	if err != nil {
		panic(err)
	}

	err = socketInput.Connect(endpoints.TcpEndpoint(endpoints.COLLECTOR))
	if err != nil {
		panic(err)
	}
	err = socketOutput.Connect(endpoints.TcpEndpoint(endpoints.FILTER_OUTPUT))
	if err != nil {
		panic(err)
	}

	rawFrame := startFrame
	rawFrameChan := make(chan *types.Frame, testLength)
	frameIndex := startFrame.FrameIdx
	wg := sync.WaitGroup{}
	go func() {
		for {
			rawFrame.FrameIdx = frameIndex
			rawFrameChan <- rawFrame
			data, err := json.Marshal(rawFrame)
			if err != nil {
				panic(err)
			}
			_, err = socketInput.SendMessage("frame", data)
			if err != nil {
				panic(err)
			}
			fmt.Println("Sent frame", frameIndex)
			mutex.Lock()
			timeStartMap[frameIndex] = time.Now()
			mutex.Unlock()

			frameIndex++
			rawFrame, err = rawReader.Next()
			time.Sleep(sendInterval)
		}
		time.Sleep(1 * time.Second)
	}()

	durationSum := int64(0)
	frameCount := 0

	totalRawErrorMm, totalFilteredErrorMm := float64(0), float64(0)
	missingRawPoints, missingFilteredPoints := 0, 0
	wg.Add(1)
	go func() {
		for range testLength {
			message, err := socketOutput.RecvMessage(0)
			if err != nil {
				panic(err)
			}
			filteredFrame := types.SmallFrame{}
			err = json.Unmarshal([]byte(strings.Join(message[1:], "")), &filteredFrame)
			if err != nil {
				panic(err)
			}
			fmt.Println("received", filteredFrame.FrameIdx)

			mutex.Lock()
			tStart := timeStartMap[filteredFrame.FrameIdx]
			mutex.Unlock()

			duration := time.Since(tStart)
			durationSum += duration.Milliseconds()
			frameCount++

			raw := <-rawFrameChan

			if pFrame.FrameIdx != raw.FrameIdx {
				fmt.Println(pFrame.FrameIdx, raw.FrameIdx)
				panic("Unequal frame index")
			}

			var compare = func(producedPlayers []struct {
				PlayerId string    `json:"playerId"`
				Number   int       `json:"number"`
				Xyz      []float64 `json:"xyz"`
				Speed    float64   `json:"speed"`
				OptaId   string    `json:"optaId"`
			}, rawPlayers []types.Player, home bool) {
				for _, player := range producedPlayers {
					rawPlayer, filteredPlayer := types.Position{}, types.Position{}
					for _, p := range filteredFrame.Players {
						if p.Home == home && p.PlayerNum == player.Number {
							filteredPlayer = p.Position
						}
					}

					for _, p := range rawPlayers {
						if p.Number == fmt.Sprintf("%d", player.Number) {
							rawPlayer = types.Position{
								X: p.Xyz[0],
								Y: p.Xyz[1],
							}
						}
					}
					if rawPlayer.X == 0 || rawPlayer.Y == 0 {
						missingRawPoints++
					} else {
						totalRawErrorMm = math.Sqrt(
							math.Pow(player.Xyz[0]-rawPlayer.X, 2)+
								math.Pow(player.Xyz[1]-rawPlayer.Y, 2)) * 100 * 100
					}
					if filteredPlayer.X == 0 || filteredPlayer.Y == 0 {
						missingFilteredPoints++
					} else {
						totalFilteredErrorMm = math.Sqrt(
							math.Pow(player.Xyz[0]-filteredPlayer.X, 2)+
								math.Pow(player.Xyz[1]-filteredPlayer.Y, 2)) * 100 * 100
					}
				}
			}

			compare(pFrame.HomePlayers, rawFrame.HomePlayers, true)
			compare(pFrame.AwayPlayers, rawFrame.AwayPlayers, false)
			//Done measuring error

			pFrame = producedReader.next()
		}
		wg.Done()
	}()

	wg.Wait()
	fmt.Printf("Results\nAverage time in pipeline: %d\nFrames processed: %d\n"+
		"Average positional deviation (filtered): %f mm\n"+
		"Missing datapoints (filtered): %d\n"+
		"Average positional deviation (raw): %f mm\n"+
		"Missing raw points (raw): %d\n",
		durationSum/int64(frameCount), frameCount,
		totalFilteredErrorMm/float64(frameCount),
		missingRawPoints,
		totalRawErrorMm/float64(frameCount),
		missingFilteredPoints)
}
