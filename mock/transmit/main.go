package main

import (
	"encoding/json"
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/frameReader"
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	zmq "github.com/pebbe/zmq4"
	"log/slog"
	"os"
	"time"
)

type TransmitterData struct {
	endpoint Endpoint
	filepath string
}

type Endpoint struct {
	protocol string
	host     string
	port     string
}

var t TransmitterData

func init() {
	slog.SetDefault(logger.New("transmitter", "DEBUG"))
	t = TransmitterData{
		endpoint: Endpoint{
			protocol: os.Getenv("PROTOCOL"),
			host:     os.Getenv("HOST"),
			port:     os.Getenv("PORT"),
		},
		filepath: os.Getenv("FILEPATH"),
	}
	if t.endpoint.protocol == "" {
		t.endpoint.protocol = "tcp"
		slog.Warn("PROTOCOL not set, using default", "protocol", t.endpoint.protocol)
	}
	if t.endpoint.host == "" {
		t.endpoint.host = "*"
		slog.Warn("HOST not set, using default", "host", t.endpoint.host)
	}
	if t.endpoint.port == "" {
		t.endpoint.port = "5555"
		slog.Warn("PORT not set, using default", "port", t.endpoint.port)
	}
	if t.filepath == "" {
		t.filepath = "raw.jsonl"
		slog.Warn("FILEPATH not set, using default", "filepath", t.filepath)
	}
}

func main() {
	Transmitter()
}

func Transmitter() {
	endpoint := fmt.Sprintf("%s://%s:%s", t.endpoint.protocol, t.endpoint.host, t.endpoint.port)
	slog.Debug("starting transmitter", "endpoint", endpoint)
	s, err := frameReader.New(t.filepath)
	if err != nil {
		slog.Error("creating frame reader", "error", err)
		return
	}
	slog.Debug("created frame reader", "filepath", t.filepath)
	zmqContext, err := zmq.NewContext()
	if err != nil {
		slog.Error("creating zmq context", "error", err)
		return
	}

	socket, err := zmqContext.NewSocket(zmq.PUB)
	if err != nil {
		slog.Error("creating zmq socket", "error", err)
		return
	}
	slog.Debug("created zmq socket", "type", zmq.PUB)
	err = socket.Bind(endpoint)
	if err != nil {
		slog.Error("binding socket", "error", err, "endpoint", endpoint)
		return
	}
	slog.Debug("bound socket", "endpoint", endpoint)

	defer func(socket *zmq.Socket) {
		err := socket.Disconnect(endpoint)
		if err != nil {
			slog.Error("disconnecting socket", "error", err)
		}
	}(socket)

	for {
		startTime := time.Now()
		msg, err := s.Next()
		if err != nil {
			slog.Error("reading next frame", "error", err)
		}
		bytes, err := json.Marshal(msg)
		if err != nil {
			slog.Error("marshalling frame", "error", err)
		}
		_, err = socket.Send(string(bytes), 0)
		if err != nil {
			slog.Error("sending frame", "error", err)
		}
		for {
			if time.Now().Sub(startTime) > time.Second/25 {
				break
			}
		}
	}
}
