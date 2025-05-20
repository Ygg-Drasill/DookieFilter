package filter

import (
	"github.com/Ygg-Drasill/DookieFilter/common/filter"
	"github.com/Ygg-Drasill/DookieFilter/common/types"
)

type savGolFilter[TElement filter.FilterableElement] struct {
	filter.Filter[TElement]
}

type filterPlayerPosition types.PlayerPosition

func (f *filterPlayerPosition) Update(key string, value float64) error {
	switch key {
	case "x":
		f.X = value
		return nil
	case "y":
		f.Y = value
		return nil
	default:

		return filter.KeyNotFoundError{}
	}
}

func (f *filterPlayerPosition) Get(key string) (float64, error) {
	switch key {
	case "x":
		return f.X, nil
	case "y":
		return f.Y, nil
	default:

		return 0, filter.KeyNotFoundError{}
	}
}

func (f *savGolFilter[TElement]) Keys() []string {
	return []string{
		"x",
		"y",
	}
}
