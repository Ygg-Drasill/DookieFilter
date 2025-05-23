package storage

import "github.com/Ygg-Drasill/DookieFilter/common/types"

type playerDistance struct {
	types.PlayerPosition
	distance float64
}

type playerDistanceSorted []playerDistance

func (pd playerDistanceSorted) Len() int {
	return len(pd)
}

func (pd playerDistanceSorted) Less(i, j int) bool {
	return pd[i].distance < pd[j].distance
}

func (pd playerDistanceSorted) Swap(i, j int) {
	pd[i], pd[j] = pd[j], pd[i]
}
