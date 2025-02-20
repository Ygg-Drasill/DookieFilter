package frame

type RawFrame struct {
	League    string `json:"league"`
	GameId    string `json:"gameId"`
	FeedName  string `json:"feedName"`
	MessageId string `json:"messageId"`
	Data      []struct {
		AwayPlayers []struct {
			Number   string    `json:"number"`
			OptaId   string    `json:"optaId"`
			PlayerId string    `json:"playerId"`
			Speed    float64   `json:"speed"`
			Xyz      []float64 `json:"xyz"`
		} `json:"awayPlayers"`
		Ball struct {
			Speed float64   `json:"speed"`
			Xyz   []float64 `json:"xyz"`
		} `json:"ball"`
		FrameIdx    int     `json:"frameIdx"`
		GameClock   float64 `json:"gameClock"`
		HomePlayers []struct {
			Number   string    `json:"number"`
			OptaId   string    `json:"optaId"`
			PlayerId string    `json:"playerId"`
			Speed    float64   `json:"speed"`
			Xyz      []float64 `json:"xyz"`
		} `json:"homePlayers"`
		Period    int   `json:"period"`
		WallClock int64 `json:"wallClock"`
	} `json:"data"`
}

type RawFrameSignal struct {
	League    string `json:"league"`
	GameId    string `json:"gameId"`
	FeedName  string `json:"feedName"`
	MessageId string `json:"messageId"`
	Data      []struct {
		EndFrameIdx    int   `json:"endFrameIdx"`
		EndWallClock   int64 `json:"endWallClock"`
		Number         int   `json:"number"`
		StartFrameIdx  int   `json:"startFrameIdx"`
		StartWallClock int64 `json:"startWallClock"`
	} `json:"data"`
}

type ProducedFrame struct {
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
