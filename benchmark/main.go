package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/socket/endpoints"
	"github.com/Ygg-Drasill/DookieFilter/common/testutils"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	zmq "github.com/pebbe/zmq4"
	"os"
	"strings"
	"sync"
	"time"
)

const sendInterval = time.Second / 100

func main() {
	mutex := sync.Mutex{}
	timeStartMap := make(map[int]time.Time)

	raw, produced := os.Args[1], os.Args[2]
	rawStat, err := os.Stat(raw)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File %s does not exist\n", raw)
	}

	producedStat, err := os.Stat(produced)
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("File %s does not exist\n", produced)
	}

	fmt.Printf("Running benchmark on raw:%s produced:%s\n", rawStat.Name(), producedStat.Name())

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

	frame := testutils.RandomFrame(11, 11)
	frame.FrameIdx = 1
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for range 2048 {
			data, err := json.Marshal(frame)
			if err != nil {
				panic(err)
			}
			_, err = socketInput.SendMessage("frame", data)
			if err != nil {
				panic(err)
			}
			fmt.Println("Sent frame", frame.FrameIdx)
			mutex.Lock()
			timeStartMap[frame.FrameIdx] = time.Now()
			mutex.Unlock()

			frame = testutils.RandomNextFrame(frame)
			time.Sleep(sendInterval)
		}
		time.Sleep(1 * time.Second)
		wg.Done()
	}()

	durationSum := int64(0)
	frameCount := 0

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
		}
	}()

	wg.Wait()
	fmt.Printf("Results\nAverage time in pipeline: %d\nFrames processed: %d", durationSum/int64(frameCount), frameCount)
}
