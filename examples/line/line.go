package main

import (
	"math"

	"github.com/adynascimento/plot/plotter"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// lines plot
	x := mat.NewDense(1, 300, plotter.Linspace(0., 1., 300))

	applySin1 := func(_, _ int, v float64) float64 { return math.Sin(15. * v) }
	applySin2 := func(_, _ int, v float64) float64 { return 0.75 * math.Sin(15.*v) }
	func1 := plotter.Apply(applySin1, x)
	func2 := plotter.Apply(applySin2, x)

	plt := plotter.NewPlot()
	plt.FigSize(11, 10)

	plt.Plot(x.RawMatrix().Data, func1.RawMatrix().Data,
		plotter.WithLineColor(plotter.Blue),
		plotter.WithLineStyle(plotter.DashDotted),
		plotter.WithMarker(plotter.Circle),
		plotter.WithMarkerSpacing(8),
	)

	plt.Plot(x.RawMatrix().Data, func2.RawMatrix().Data,
		plotter.WithLineColor(plotter.Red),
		plotter.WithMarker(plotter.Square),
		plotter.WithMarkerSpacing(8),
	)

	plt.Title("plot example")
	plt.XLabel("xLabel")
	plt.YLabel("yLabel")
	plt.Legend("line1", "line2")
	plt.Grid()

	plt.Save("line.png")
}
