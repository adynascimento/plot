package plotter

import (
	"image/color"
	"log"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type scatterOptions struct {
	color      color.Color
	gradient   colorgrad.Gradient
	marker     draw.GlyphDrawer
	markerSize font.Length
	colorBar   colorBar
}

// parameters to scatter plot
func (plt *plotParameters) Scatter(x, y, z []float64, options ...func(*scatterOptions)) {
	// default options
	plt.scatter = scatterOptions{
		color:      Blue,
		marker:     Circle,
		markerSize: vg.Points(3),
	}

	// apply additional options
	for _, option := range options {
		option(&plt.scatter)
	}
	plt.colorBar = plt.scatter.colorBar

	// prepare data to plot
	var xys plotter.XYer
	if len(z) == 0 {
		pts := make(plotter.XYs, len(x))
		for i := range pts {
			pts[i].X = x[i]
			pts[i].Y = y[i]
		}
		xys = pts
	} else {
		pts := make(plotter.XYZs, len(x))
		for i := range pts {
			pts[i].X = x[i]
			pts[i].Y = y[i]
			pts[i].Z = z[i]
		}
		xys = pts
	}

	// add colormap and make a scatter plotter
	sc, err := plotter.NewScatter(xys)
	if err != nil {
		log.Panic(err)
	}
	sc.GlyphStyle = draw.GlyphStyle{
		Color:  plt.scatter.color,
		Radius: plt.scatter.markerSize,
		Shape:  plt.scatter.marker,
	}

	if plt.scatter.gradient != (colorgrad.Gradient{}) {
		// specify style and color for individual points.
		sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
			colors := plt.scatter.gradient.Colors(uint(len(z)))
			return draw.GlyphStyle{Color: colors[i], Radius: plt.scatter.markerSize,
				Shape: plt.scatter.marker}
		}
	}

	// add the plotters to the plot
	plt.plot.Add(sc)

	if plt.colorBar.show {
		// get min and max values
		plt.colorBar.min = floats.Min(z)
		plt.colorBar.max = floats.Max(z)
	}
}

func WithScatterMarkerColor(color colorType) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.color = color
	}
}

func WithScatterGradient(gradient colorgrad.Gradient) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.gradient = gradient
		so.colorBar.gradient = gradient
	}
}

func WithScatterMarker(marker markerType) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.marker = marker
	}
}

func WithScatterMarkerSize(size float64) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.markerSize = vg.Points(size)
	}
}

func WithScatterColorbar(position positionType) func(*scatterOptions) {
	return func(so *scatterOptions) {
		if so.gradient == (colorgrad.Gradient{}) {
			so.colorBar.show = false
		} else {
			so.colorBar.show = true
			so.colorBar.position = position
		}
	}
}
