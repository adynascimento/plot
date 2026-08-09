package plotter

import (
	"image/color"
	"log"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type scatterOptions struct {
	color      color.Color
	colormap   Colormap
	marker     draw.GlyphDrawer
	markerSize font.Length
	colorbar   colorbar
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
	plt.colorbar = plt.scatter.colorbar

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

	if plt.scatter.colormap != nil && len(z) > 0 {
		min := floats.Min(z)
		max := floats.Max(z)
		colors := plt.scatter.colormap.Colors(len(z))

		// specify style and color for individual points.
		// normalize to determine its percentage/position in the color palette
		sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
			var percent float64
			if max == min {
				percent = 0.5
			} else {
				percent = (z[i] - min) / (max - min)
			}
			colorIdx := int(percent * float64(len(z)-1))
			return draw.GlyphStyle{
				Color:  colors[colorIdx],
				Radius: plt.scatter.markerSize,
				Shape:  plt.scatter.marker,
			}
		}
	}

	// add the plotters to the plot
	plt.plot.Add(sc)

	if plt.colorbar.show && len(z) > 0 {
		// get min and max values
		plt.colorbar.min = floats.Min(z)
		plt.colorbar.max = floats.Max(z)
	}
}

func WithScatterMarkerColor(color colorType) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.color = color
	}
}

func WithScatterColorMap(cmap Colormap) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.colormap = cmap
		so.colorbar.colormap = cmap
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
		if so.colormap == nil {
			so.colorbar.show = false
		} else {
			so.colorbar.show = true
			so.colorbar.position = position
		}
	}
}
