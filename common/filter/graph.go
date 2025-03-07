package filter

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"log"
)

func ShowGraph(data []float64, label string) {

	TestXdata, TestYdata, _ := splitIntoThree(data)

	// Convert positions to plotter.XYs
	pts := make(plotter.XYs, len(TestYdata))
	for i, _ := range TestXdata {
		pts[i].X = TestYdata[i]
		pts[i].Y = TestYdata[i]
	}

	// Create a new plot
	p := plot.New()
	p.Title.Text = "Soccer Player Positions"
	p.X.Label.Text = "X Position"
	p.Y.Label.Text = "Y Position"

	// Create scatter plot
	s, err := plotter.NewScatter(pts)
	if err != nil {
		log.Fatal(err)
	}

	p.Add(s)

	// Save the plot to a PNG file
	if err := p.Save(6*vg.Inch, 6*vg.Inch, label); err != nil {
		log.Fatal(err)
	}
}
