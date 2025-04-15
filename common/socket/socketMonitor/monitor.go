package socketMonitor

import (
	"fmt"
	zmq "github.com/pebbe/zmq4"
	"sync"
)

func Run(ctx *zmq.Context, address string, wg *sync.WaitGroup) {
	socket, err := ctx.NewSocket(zmq.PAIR)
	if err != nil {
		panic(err)
	}

	err = socket.Connect(address)
	if err != nil {
		panic(err)
	}
	for {
		a, b, c, err := socket.RecvEvent(0)
		if err != nil {
			panic(err)
		}
		fmt.Printf("---\n%s %s %d\n---\n", a.String(), b, c)
	}

	socket.Close()
	wg.Done()
}
