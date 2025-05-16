package storage

import "github.com/Ygg-Drasill/DookieFilter/common/types"

type FrameNearest struct {
	Target types.PlayerPosition `json:"target"`
	Home   playerDistanceSorted `json:"home"`
	Away   playerDistanceSorted `json:"away"`
}

type FrameRangeNearestResponse []FrameNearest
