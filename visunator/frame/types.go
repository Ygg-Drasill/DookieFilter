package frame

type Player struct {
	Number   string    `json:"number"`
	OptaId   string    `json:"optaId"`
	PlayerId string    `json:"playerId"`
	Speed    float64   `json:"speed"`
	Xyz      []float64 `json:"xyz"`
}

type Frame[DType DataPlayer | DataSignal] struct {
	League    string  `json:"league"`
	GameId    string  `json:"gameId"`
	FeedName  string  `json:"feedName"`
	MessageId string  `json:"messageId"`
	Data      []DType `json:"data"`
}

type DataPlayer struct {
	AwayPlayers []Player `json:"awayPlayers"`
	Ball        struct {
		Speed float64   `json:"speed"`
		Xyz   []float64 `json:"xyz"`
	} `json:"ball"`
	FrameIdx    int      `json:"frameIdx"`
	GameClock   float64  `json:"gameClock"`
	HomePlayers []Player `json:"homePlayers"`
	Period      int      `json:"period"`
	WallClock   int64    `json:"wallClock"`
}

type DataSignal struct {
	EndFrameIdx    int   `json:"endFrameIdx"`
	EndWallClock   int64 `json:"endWallClock"`
	Number         int   `json:"number"`
	StartFrameIdx  int   `json:"startFrameIdx"`
	StartWallClock int64 `json:"startWallClock"`
}

type ProducedFrame struct {
	Period      int      `json:"period"`
	FrameIdx    int      `json:"frameIdx"`
	GameClock   float64  `json:"gameClock"`
	WallClock   int64    `json:"wallClock"`
	HomePlayers []Player `json:"homePlayers"`
	AwayPlayers []Player `json:"awayPlayers"`
	Ball        struct {
		Xyz   []float64 `json:"xyz"`
		Speed float64   `json:"speed"`
	} `json:"ball"`
	Live      bool   `json:"live"`
	LastTouch string `json:"lastTouch"`
}
