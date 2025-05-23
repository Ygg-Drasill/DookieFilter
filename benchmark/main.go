package main

import (
	"errors"
	"fmt"
	zmq "github.com/pebbe/zmq4"
	zmq4 "github.com/pebbe/zmq4/draft"
	"os"
)

func main() {
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

	socketInput.Connect("tcp://" + )
	socketOutput
}
