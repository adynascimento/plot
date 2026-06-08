package main

import (
	"math/rand/v2"

	"github.com/adynascimento/plot/plotter"
)

func main() {
	// histogram plot
	x := make([]float64, 1000)
	for i := range x {
		x[i] = rand.NormFloat64()
	}

	plt := plotter.NewPlot()
	plt.FigSize(10, 10)

	plt.Hist(x, 16,
		plotter.WithHistFillColor(plotter.Gray),
		plotter.WithHistNormalCurve(plotter.Red),
	)
	plt.Title("histogram plot example")
	plt.XLabel("x")

	plt.Show()
}
