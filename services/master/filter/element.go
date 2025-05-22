package filter

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/filter"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
	"strconv"
	"strings"
)

type filterableFrame types.SmallFrame

func (f *filterableFrame) Update(key string, value float64) error {
	h, n, x := decodeKey(key)
	for i, player := range f.Players {
		if player.Home == h && player.PlayerNum == n {
			if x {
				f.Players[i].X = value
			} else {
				f.Players[i].Y = value
			}
		}
	}

	return filter.KeyNotFoundError{}
}

func (f *filterableFrame) Get(key string) (float64, error) {
	h, n, x := decodeKey(key)
	for _, player := range f.Players {
		if player.Home == h && player.PlayerNum == n {
			if x {
				return player.X, nil
			} else {
				return player.Y, nil
			}
		}
	}

	return 0, filter.KeyNotFoundError{}
}

type savGolFilter struct {
	filter.Filter[*filterableFrame]
}

func (f savGolFilter) Keys() []string {
	keys := make([]string, 0)
	for _, player := range f.Elements[0].Players {
		k := fmt.Sprintf("%t_%d", player.Home, player.PlayerNum)
		keys = append(keys, fmt.Sprintf("%s_%s", k, "x"))
		keys = append(keys, fmt.Sprintf("%s_%s", k, "y"))
	}

	return keys
}

func decodeKey(key string) (bool, int, bool) {
	parts := strings.Split(key, "_")
	home, _ := strconv.ParseBool(parts[0])
	number, _ := strconv.ParseInt(parts[1], 10, 64)
	axis := parts[2] == "x"
	return home, int(number), axis
}
