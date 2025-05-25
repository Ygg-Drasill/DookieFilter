package main

import (
	"log/slog"
	"os"
	"visunator/app"

	"github.com/Ygg-Drasill/DookieFilter/common/logger"
)

func init() {
	slog.SetDefault(logger.New("visunator", "DEBUG"))
}

var dataPath = os.Args[1] //TODO: Change me

func main() {
	a := app.NewFromReader(dataPath)
	a.Run()
}
