package detector

type holeMessage struct {
	FrameIdx     int  `json:"frameIdx"`
	PlayerNumber int  `json:"playerNumber"`
	Home         bool `json:"home"`
}
