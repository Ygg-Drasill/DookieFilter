package main

type producedFrame struct {
	Period      int     `json:"period"`
	FrameIdx    int     `json:"frameIdx"`
	GameClock   float64 `json:"gameClock"`
	WallClock   int64   `json:"wallClock"`
	HomePlayers []struct {
		PlayerId string    `json:"playerId"`
		Number   int       `json:"number"`
		Xyz      []float64 `json:"xyz"`
		Speed    float64   `json:"speed"`
		OptaId   string    `json:"optaId"`
	} `json:"homePlayers"`
	AwayPlayers []struct {
		PlayerId string    `json:"playerId"`
		Number   int       `json:"number"`
		Xyz      []float64 `json:"xyz"`
		Speed    float64   `json:"speed"`
		OptaId   string    `json:"optaId"`
	} `json:"awayPlayers"`
	Ball struct {
		Xyz   []float64 `json:"xyz"`
		Speed float64   `json:"speed"`
	} `json:"ball"`
	Live      bool   `json:"live"`
	LastTouch string `json:"lastTouch"`
}
