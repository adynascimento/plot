package main

import (
	"math"
	"github.com/adynascimento/plot/plot"
)

func main() {
	n := 300
	x := plot.Linspace(0.0, 1.0, n)
	func1 := make([]float64, n)
	func2 := make([]float64, n)
	for i := range x {
		func1[i] = math.Sin(15. * x[i])
		func2[i] = 0.5 * math.Sin(15.*x[i])
	}

	plt := plot.NewPlot()
	plt.FigSize(11, 10)

	plt.Plot(x, func1)
	plt.Plot(x, func2)
	plt.Title("plot example")
	plt.XLabel("x_label")
	plt.YLabel("y_label")
	plt.Legend("line1", "line2")

	plt.Save("figures/lines.png")
}
