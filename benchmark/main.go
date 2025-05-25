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
const testLength = 25 * 120

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
	wg.Add(1)
	go func() {
		for range testLength {
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
		wg.Done()
	}()

	durationSum := int64(0)
	frameCount := 0

	totalRawError, totalFilteredError := float64(0), float64(0)
	go func() {
		for {
			message, err := socketOutput.RecvMessage(0)
			if err != nil {
				panic(err)
			}
			filteredFrame := types.SmallFrame{}
			fmt.Println("received", filteredFrame.FrameIdx)
			err = json.Unmarshal([]byte(strings.Join(message[1:], "")), &filteredFrame)
			if err != nil {
				panic(err)
			}

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

			//Measure error

			for _, filteredPlayer := range filteredFrame.Players {
				rawPlayer, producedPlayer := types.Position{}, types.Position{}
				var producedTeam []struct {
					PlayerId string    `json:"playerId"`
					Number   int       `json:"number"`
					Xyz      []float64 `json:"xyz"`
					Speed    float64   `json:"speed"`
					OptaId   string    `json:"optaId"`
				}
				rawTeam := []types.Player{}

				if filteredPlayer.Home {
					rawTeam = raw.HomePlayers
					producedTeam = pFrame.HomePlayers
				} else {
					rawTeam = raw.AwayPlayers
					producedTeam = pFrame.AwayPlayers
				}

				for _, rp := range rawTeam {
					if rp.Number == fmt.Sprintf("%d", filteredPlayer.PlayerNum) {
						rawPlayer = types.Position{
							X: rp.Xyz[0],
							Y: rp.Xyz[1],
						}
					}
				}

				for _, pp := range producedTeam {
					if pp.Number == filteredPlayer.PlayerNum {
						rawPlayer = types.Position{
							X: pp.Xyz[0],
							Y: pp.Xyz[1],
						}
					}
				}

				totalFilteredError = math.Sqrt(
					math.Pow(producedPlayer.X-filteredPlayer.X, 2) +
						math.Pow(producedPlayer.Y-filteredPlayer.Y, 2))
				totalRawError = math.Sqrt(
					math.Pow(producedPlayer.X-rawPlayer.X, 2) +
						math.Pow(producedPlayer.Y-rawPlayer.Y, 2))
			}

			//Done measuring error

			pFrame = producedReader.next()
		}
	}()

	wg.Wait()
	fmt.Printf("Results\nAverage time in pipeline: %d\nFrames processed: %d\nAverage positional deviation (filtered): %f\nAverage positional deviation (Raw): %f\n",
		durationSum/int64(frameCount), frameCount, totalFilteredError/float64(frameCount), totalRawError/float64(frameCount))
}
