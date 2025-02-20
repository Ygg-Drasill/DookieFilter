package main

import (
	zmq "github.com/pebbe/zmq4"
	"log"
	"time"
)

const (
	endpoint = "tcp://localhost:5555"
)

func main() {
	Transmitter()
}

func Transmitter() {
	zmqContext, err := zmq.NewContext()
	if err != nil {
		log.Fatal(err)
	}

	socket, err := zmqContext.NewSocket(zmq.PUB)
	if err != nil {
		log.Fatal(err)
	}
	err = socket.Bind(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	defer func(socket *zmq.Socket) {
		err := socket.Disconnect(endpoint)
		if err != nil {
			log.Fatal(err)
		}
	}(socket)

	for {
		msg := "Hello"
		b, err := socket.Send(msg, 0)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(b)
		time.Sleep(2 * time.Second)
	}
}
