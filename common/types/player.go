package types

import "github.com/Ygg-Drasill/DookieFilter/common/pringleBuffer"

type Player struct {
	Number   string    `json:"number"`
	OptaId   string    `json:"optaId"`
	PlayerId string    `json:"playerId"`
	Speed    float64   `json:"speed"`
	Xyz      []float64 `json:"xyz"`
}

type PlayerPosition struct {
	FrameIdx int     `json:"frameIdx"`
	X        float64 `json:"x"`
	Y        float64 `json:"y"`
}

func (pp PlayerPosition) Key() pringleBuffer.Key {
	return pringleBuffer.Key(pp.FrameIdx)
}
