package main

import (
	"math"

	"github.com/adynascimento/plot/plot"
	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// heatmap plot
	n := 300
	x := mat.NewDense(1, n, plot.Linspace(-3.0, 3.0, n))
	y := mat.NewDense(1, n, plot.Linspace(-3.0, 3.0, n))
	Z := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := math.Sin(x.At(0, j)) * math.Cos(y.At(0, i))
			Z.Set(i, j, v)
		}
	}

	plt := plot.NewPlot()
	plt.FigSize(10, 10)

	plt.HeatMap(x, y, Z, 12, colorgrad.Viridis())
	plt.Title("heatmap plot example")
	plt.XLabel("x_label")
	plt.YLabel("y_label")

	plt.Save("heatmap.png")
}
