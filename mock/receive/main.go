package main

import (
	zmq "github.com/pebbe/zmq4"
	"log"
)

const (
	endpoint = "tcp://localhost:5555"
)

func main() {
	Receiver()
}

func Receiver() {
	context, err := zmq.NewContext()
	if err != nil {
		log.Fatal(err)
	}

	socket, err := context.NewSocket(zmq.SUB)
	if err != nil {
		log.Fatal(err)
	}

	err = socket.Connect(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	err = socket.SetSubscribe("")
	if err != nil {
		log.Fatal(err)
	}

	for {
		m, err := socket.Recv(0)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(m)
	}
}
