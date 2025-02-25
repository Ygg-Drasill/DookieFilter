package main

import (
	"visunator/app"
)

const dataPath = "./raw.jsonl" //TODO: Change me

func main() {
	a := app.NewFromReader(dataPath)
	a.Run()
}
