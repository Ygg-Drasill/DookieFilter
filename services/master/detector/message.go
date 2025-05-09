package detector

type holeMessage struct {
	FrameIdx     int  `json:"frameIdx"`
	PlayerNumber string `json:"playerNumber"`
	Home         bool `json:"home"`
}
