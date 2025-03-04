package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"github.com/Ygg-Drasill/DookieFilter/services/master/collector"
	zmq "github.com/pebbe/zmq4"
	"log/slog"
	"sync"
)

func init() {
	slog.SetDefault(logger.New("master", "DEBUG"))
}

const (
	endpointDealSwap    = "deal-swap"
	endpointDealJump    = "deal-jump"
	endpointDealClean   = "deal-clean"
	endpointDealMissing = "deal-missing"
)

func main() {
	slog.Info("starting service")
	socketCtx, err := zmq.NewContext()
	if err != nil {
		slog.Error(err.Error())
	}
	wCollector := collector.New(collector.WithSocketContext(socketCtx))
	workerWaitGroup := sync.WaitGroup{}

	go wCollector.Run(&workerWaitGroup)
	workerWaitGroup.Wait()
}
