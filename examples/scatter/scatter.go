package main

import (
	"math/rand"

	"github.com/adynascimento/plot/plotter"
	"github.com/mazznoer/colorgrad"
)

func main() {
	// scatter plot
	rnd := rand.New(rand.NewSource(1))

	n := 15
	x := make([]float64, n)
	y := make([]float64, n)
	Z := make([]float64, n)
	for i := range x {
		x[i] = rnd.Float64()
		y[i] = rnd.Float64()
		Z[i] = 30.0 * rnd.Float64()
	}

	plt := plotter.NewPlot()
	plt.FigSize(10, 8)

	plt.Scatter(x, y, Z, colorgrad.Viridis())
	plt.Title("scatter plot example")
	plt.XLabel("xLabel")
	plt.YLabel("yLabel")
	plt.Grid()

	plt.Save("scatter.png")
}
