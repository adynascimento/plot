package main

import (
	"math/rand/v2"

	"github.com/adynascimento/plot/plotter"
)

func main() {
	// histogram plot
	x := make([]float64, 5000)
	for i := range x {
		if rand.Float64() < 0.65 {
			x[i] = 10 + 0.7*rand.NormFloat64()
		} else {
			x[i] = 12.2 + 0.6*rand.NormFloat64()
		}
	}

	plt := plotter.NewPlot()
	plt.FigSize(10, 10)

	plt.Hist(x, 16,
		plotter.WithHistFillColor(plotter.Gray),
		plotter.WithHistKDECurve(plotter.Blue),
		plotter.WithHistNormalCurve(plotter.Red),
	)
	plt.Title("histogram plot example")
	plt.Legend("kde curve", "normal curve")
	plt.XLabel("x")

	plt.Show()
}
