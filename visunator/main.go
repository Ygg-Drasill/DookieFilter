package main

import (
	"log/slog"
	"visunator/app"

	"github.com/Ygg-Drasill/DookieFilter/common/logger"
)

func init() {
	slog.SetDefault(logger.New("visunator", "DEBUG"))
}

const dataPath = "./raw.jsonl" //TODO: Change me

func main() {
	a := app.NewFromReader(dataPath)
	a.Run()
}
