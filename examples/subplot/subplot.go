package main

import (
	"math"

	"github.com/adynascimento/plot/plot"
)

func main() {
	// lines plot
	n := 300
	x := plot.Linspace(0.0, 1.0, n)
	func1 := make([]float64, n)
	func2 := make([]float64, n)
	func3 := make([]float64, n)
	func4 := make([]float64, n)
	for i := range x {
		func1[i] = math.Sin(25. * x[i])
		func2[i] = 0.75 * math.Sin(25.*x[i])
		func3[i] = math.Tan(15. * x[i])
		func4[i] = 0.5 * math.Tan(15.*x[i])
	}

	plt := plot.NewSubplot(1, 2)
	plt.FigSize(23, 10)

	subplt1 := plt.Subplot(0, 0)
	subplt1.Plot(x, func1)
	subplt1.Plot(x, func2)
	subplt1.Title("sin function")
	subplt1.XLabel("xLabel")
	subplt1.YLabel("yLabel")
	subplt1.Legend("sin1", "sin2")

	subplt2 := plt.Subplot(0, 1)
	subplt2.Plot(x, func3)
	subplt2.Plot(x, func4)
	subplt2.Title("tan function")
	subplt2.XLabel("xLabel")
	subplt2.YLabel("yLabel")
	subplt2.Legend("tan1", "tan2")

	plt.Save("subplot.png")
}
