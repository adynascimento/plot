package examples

import (
	"math"
	"math/rand"
	"plot/plot"

	ngo "plot/numeric"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
)

// lines plot
func PlotLines() {
	n := 300
	x := make([]float64, n)
	y1 := make([]float64, n)
	y2 := make([]float64, n)
	for i := range x {
		x[i] = float64(i) * (1.0 / float64(n-1))
		y1[i] = math.Sin(15. * x[i])
		y2[i] = 0.5 * math.Sin(15.*x[i])
	}

	plt := plot.NewPlot()
	plt.FigSize(11, 10)

	plt.Plot(x, y1)
	plt.Plot(x, y2)
	plt.Title("plot example")
	plt.XLabel("x_label")
	plt.YLabel("y_label")
	plt.Legend("line1", "line2")

	plt.Save("figures/lines.png")
}

// heatmap plot
func PlotHeatMap() {
	n := 300
	x := mat.NewDense(1, n, ngo.Linspace(-3.0, 3.0, n))
	y := mat.NewDense(1, n, ngo.Linspace(-3.0, 3.0, n))
	Z := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := math.Sin(x.At(0, j)) * math.Cos(y.At(0, i))
			Z.Set(i, j, v)
		}
	}

	plt := plot.NewHeatMap()
	plt.FigSize(10, 10)

	plt.HeatMap(x, y, Z, 12, colorgrad.Viridis())
	plt.Title("heatmap plot example")
	plt.XLabel("x_label")
	plt.YLabel("y_label")

	plt.Save("figures/heatmap.png")
}

// contour plot
func PlotContour() {
	n := 300
	x := mat.NewDense(1, n, ngo.Linspace(-3.0, 3.0, n))
	y := mat.NewDense(1, n, ngo.Linspace(-3.0, 3.0, n))
	Z := mat.NewDense(n, n, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			v := math.Sin(x.At(0, j)) * math.Cos(y.At(0, i))
			Z.Set(i, j, v)
		}
	}

	plt := plot.NewContour()
	plt.FigSize(10, 10)

	plt.Contour(x, y, Z, 12, colorgrad.Viridis())
	plt.Title("contour plot example")
	plt.XLabel("x_label")
	plt.YLabel("y_label")

	plt.Save("figures/contour.png")
}

// scatter plot
func PlotScatter() {
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

	plt := plot.NewScatter()
	plt.FigSize(10, 8)

	plt.Scatter(x, y, Z, colorgrad.Viridis())
	plt.Title("scatter plot example")
	plt.XLabel("x_label")
	plt.YLabel("y_label")

	plt.Save("figures/scatter.png")
}
