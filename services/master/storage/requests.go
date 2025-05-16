package storage

import "github.com/Ygg-Drasill/DookieFilter/common/types"

type FrameRangeNearestRequest struct {
	StartIndex   int             `json:"startIdx"`
	EndIndex     int             `json:"endIdx"`
	NearestCount int             `json:"n"`
	TargetPlayer types.PlayerKey `json:"target"`
}
