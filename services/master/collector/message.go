package collector

type holeMessage struct {
	FrameIdx     int64 `json:"frameIdx"`
	PlayerNumber int   `json:"playerNumber"`
}
