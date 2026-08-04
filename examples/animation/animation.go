package main

import (
	"math"
	"time"

	"github.com/adynascimento/plot/plotter"
	"gonum.org/v1/gonum/mat"
)

func main() {
	x := mat.NewDense(1, 300, plotter.Linspace(0, 2*math.Pi, 300))

	plt := plotter.NewPlot()
	plt.FigSize(11, 10)

	update := func(frame int) {
		plt.Clear()

		// recompute the wave every frame
		phase := float64(frame) * 0.10
		y1 := plotter.Apply(func(_, _ int, v float64) float64 { return math.Sin(2*v - phase) }, x)
		y2 := plotter.Apply(func(_, _ int, v float64) float64 { return math.Cos(2*v - phase) }, x)

		plt.Plot(x.RawMatrix().Data, y1.RawMatrix().Data,
			plotter.WithLineColor(plotter.Blue),
			plotter.WithLineWidth(2),
		)

		plt.Plot(x.RawMatrix().Data, y2.RawMatrix().Data,
			plotter.WithLineColor(plotter.Red),
			plotter.WithLineWidth(2),
		)

		plt.Title("animation example")
		plt.Legend("sin", "cos").Location(plotter.LowerLeft)
		plt.XLabel("xLabel")
		plt.YLabel("yLabel")
		plt.XLim(0, 2*math.Pi)
		plt.YLim(-1.2, 1.2)
		plt.Grid()
	}

	nFrames := 240
	plt.Animation(nFrames, update,
		plotter.WithAnimationInterval(40*time.Millisecond),
		plotter.WithAnimationLoop(true),
	)

	plt.Save("animation.gif")
	plt.Show()
}
