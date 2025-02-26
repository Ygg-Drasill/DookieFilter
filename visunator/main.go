package main

import (
	"github.com/Ygg-Drasill/DookieFilter/common/logger"
	"log/slog"
	"visunator/app"
)

func init() {
	slog.SetDefault(logger.New("visunator", "DEBUG"))
}

const dataPath = "./raw.jsonl" //TODO: Change me

func main() {
	a := app.NewFromReader(dataPath)
	a.Run()
}
