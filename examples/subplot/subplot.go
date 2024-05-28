package main

import (
	"math"

	"github.com/adynascimento/plot/plotter"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// lines plot
	x := mat.NewDense(1, 300, plotter.Linspace(0., 1., 300))

	applySin1 := func(_, _ int, v float64) float64 { return math.Sin(25. * v) }
	applySin2 := func(_, _ int, v float64) float64 { return 0.75 * math.Sin(25.*v) }
	func1 := plotter.Apply(applySin1, x)
	func2 := plotter.Apply(applySin2, x)

	applyTan1 := func(_, _ int, v float64) float64 { return math.Tan(15. * v) }
	applyTan2 := func(_, _ int, v float64) float64 { return 0.5 * math.Tan(15.*v) }
	func3 := plotter.Apply(applyTan1, x)
	func4 := plotter.Apply(applyTan2, x)

	plt := plotter.NewSubplot(1, 2)
	plt.FigSize(23, 10)

	subplt1 := plt.Subplot(0, 0)
	subplt1.Plot(x.RawMatrix().Data, func1.RawMatrix().Data)
	subplt1.Plot(x.RawMatrix().Data, func2.RawMatrix().Data)
	subplt1.Title("sin function")
	subplt1.XLabel("xLabel")
	subplt1.YLabel("yLabel")
	subplt1.Legend("sin1", "sin2")
	subplt1.Grid()

	subplt2 := plt.Subplot(0, 1)
	subplt2.Plot(x.RawMatrix().Data, func3.RawMatrix().Data)
	subplt2.Plot(x.RawMatrix().Data, func4.RawMatrix().Data)
	subplt2.Title("tan function")
	subplt2.XLabel("xLabel")
	subplt2.YLabel("yLabel")
	subplt2.Legend("tan1", "tan2")
	subplt2.Grid()

	plt.Save("subplot.png")
}
