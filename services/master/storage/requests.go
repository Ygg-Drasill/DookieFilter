package storage

type FrameRangeRequest struct {
	StartIndex   int `json:"startIdx"`
	EndIndex     int `json:"endIdx"`
	NearestCount int `json:"n"`
	TargetPlayer int `json:"target"`
}
