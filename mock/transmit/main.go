package main

import (
	"encoding/json"
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	zmq "github.com/pebbe/zmq4"
	"log"
	"log/slog"
	"os"
	"time"
)

var endpoint string

func init() {
	endpoint = os.Getenv("ENDPOINT")
}

func main() {
	Transmitter()
}

func Transmitter() {
	if endpoint == "" {
		endpoint = "tcp://localhost:5555"
		slog.Info("ENDPOINT not set, using default", "endpoint", endpoint)
	}
	s := frameReader.New("raw.jsonl")
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
		startTime := time.Now()
		msg := s.Next()
		bytes, err := json.Marshal(msg)
		if err != nil {
			log.Fatal(err)
		}
		b, err := socket.Send(string(bytes), 0)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(b)
		for {
			if time.Now().Sub(startTime) > time.Second/25 {
				break
			}
		}
	}
}
