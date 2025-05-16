package filter

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"math/rand"
	"testing"
)

type testElement struct {
	x float64
}

func (t *testElement) Keys() []string {
	return []string{"x"}
}

func (t *testElement) Update(key string, val float64) error {
	if key == "x" {
		t.x = val
		return nil
	}
	return KeyNotFoundError{}
}

func (t *testElement) Get(key string) (float64, error) {
	if key == "x" {
		return t.x, nil
	}
	return 0, KeyNotFoundError{}
}

//func TestGolayFilter(t *testing.T) {
//	f := testFilter{
//		elements: make([]FilterableElement, 0),
//	}
//}

func TestPlot(t *testing.T) {
	const points = 50
	const scale float64 = 2
	noiseSine := make(plotter.XYs, points)
	sine := make(plotter.XYs, points)
	filtered := make(plotter.XYs, points)
	for i := range noiseSine {
		x := float64(i)
		y := math.Sin(x/float64(points)*math.Pi*4+math.Pi/2) * 2
		sine[i].X = x
		sine[i].Y = y

		noiseSine[i].X = x
		noiseSine[i].Y = y + (rand.Float64()-0.5)*scale
	}
	filter := NewSavitzkyGolayFilter[*testElement](5, 0)
	j := 0
	for i := range noiseSine {
		y, err := filter.Step(&testElement{x: noiseSine[i].Y})
		if errors.Is(err, NotFullError{}) {
			continue
		}
		x := float64(j)
		filtered[j].X = x
		filtered[j].Y = (*y).x
		j++
	}

	for j < points {
		x := float64(j)
		filtered[j].X = x
		filtered[j].Y = noiseSine[j].Y
		j++
	}

	var noiseDiff, filterDiff float64
	for i := range points {
		sine := sine[i].Y
		noiseDiff += math.Abs(sine - noiseSine[i].Y)
		filterDiff += math.Abs(sine - filtered[i].Y)
	}

	assert.Less(t, filterDiff, noiseDiff, "Filtered difference from clean signal should be less than compared to noisy signal")

	p := plot.New()
	p.Title.Text = "Sinus noise filter"
	err := plotutil.AddLinePoints(p,
		"Sine", sine,
		"Noise", noiseSine,
		"Filter", filtered)
	assert.NoError(t, err)
	err = p.Save(4*vg.Inch, 4*vg.Inch, "points.png")
	assert.NoError(t, err)
}
