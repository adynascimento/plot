package main

import (
	"math"

	"github.com/adynascimento/plot/plotter"
	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// contourf plot
	n := 300
	x := mat.NewDense(1, n, plotter.Linspace(-3.0, 3.0, n))
	y := mat.NewDense(1, n, plotter.Linspace(-3.0, 3.0, n))
	Z := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := math.Sin(x.At(0, j)) * math.Cos(y.At(0, i))
			Z.Set(i, j, v)
		}
	}

	plt := plotter.NewPlot()
	plt.FigSize(10, 10)

	plt.ContourF(x, y, Z,
		plotter.WithLevels(12),
		plotter.WithGradient(colorgrad.Viridis()),
		plotter.WithContourLines(),
		plotter.WithContourLineStyle(plotter.Dashed),
		plotter.WithColorbar(plotter.Vertical),
	)
	plt.Title("contourf plot example")
	plt.XLabel("xLabel")
	plt.YLabel("yLabel")

	plt.Save("contourf.png")
}
