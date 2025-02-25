package main

import (
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	zmq "github.com/pebbe/zmq4"
	"log"
	"log/slog"
	"os"
	"time"
)

type Endpoint struct {
	Protocol string
	Host     string
	Port     string
}

var e Endpoint

func init() {
	slog.SetDefault(logger.New("transmitter"))
	e = Endpoint{
		Protocol: os.Getenv("PROTOCOL"),
		Host:     os.Getenv("HOST"),
		Port:     os.Getenv("PORT"),
	}
	if e.Protocol == "" {
		e.Protocol = "tcp"
		slog.Info("PROTOCOL not set, using default", "protocol", e.Protocol)
	}
	if e.Host == "" {
		e.Host = "*"
		slog.Info("HOST not set, using default", "host", e.Host)
	}
	if e.Port == "" {
		e.Port = "5555"
		slog.Debug("PORT not set, using default", "port", e.Port)
	}
}

func main() {
	Transmitter()
}

func Transmitter() {
	endpoint := fmt.Sprintf("%s://%s:%s", e.Protocol, e.Host, e.Port)
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
