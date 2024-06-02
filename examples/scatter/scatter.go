package main

import (
	"math"
	"math/rand"

	"github.com/adynascimento/plot/plotter"
	"github.com/mazznoer/colorgrad"
)

func main() {
	// scatter plot
	rnd := rand.New(rand.NewSource(1))

	n := 500
	theta := plotter.Linspace(0, 1, n)
	x := make([]float64, n)
	y := make([]float64, n)
	Z := make([]float64, n)
	for i := range x {
		x[i] = math.Exp(theta[i]) * math.Sin(100.*theta[i])
		y[i] = math.Exp(theta[i]) * math.Cos(100.*theta[i])
		Z[i] = math.Cos(30. * rnd.Float64())
	}

	plt := plotter.NewPlot()
	plt.FigSize(10, 9)

	plt.Scatter(x, y, Z,
		plotter.WithScatterGradient(colorgrad.Viridis()),
		plotter.WithScatterMarker(plotter.Circle),
		plotter.WithScatterColorbar(plotter.Vertical),
	)
	plt.Title("scatter plot example")
	plt.XLabel("xLabel")
	plt.YLabel("yLabel")
	plt.XLim(-3, 3)
	plt.YLim(-3, 3)
	plt.Grid()

	plt.Save("scatter.png")
}
