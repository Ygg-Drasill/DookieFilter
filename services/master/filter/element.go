package filter

import (
	"fmt"
	"github.com/Ygg-Drasill/DookieFilter/common/filter"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
)

type filterableFrame types.SmallFrame

func (f filterableFrame) Update(key string, value float64) error {
	//TODO implement me
	panic("implement me")
}

func (f filterableFrame) Get(key string) (float64, error) {
	//TODO implement me
	panic("implement me")
}

type savGolFilter struct {
	filter.Filter[filterableFrame]
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
