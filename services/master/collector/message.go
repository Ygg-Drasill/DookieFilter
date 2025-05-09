package collector

type holeMessage struct {
	FrameIdx     int  `json:"frameIdx"`
	PlayerNumber int  `json:"playerNumber"`
	HomePlayer   bool `json:"homePlayer"`
}
