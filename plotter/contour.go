package plotter

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type contourOptions struct {
	nLevels      int
	gradient     colorgrad.Gradient
	lineSettings lineSettings
	colorBar     colorBar
}

type lineSettings struct {
	show  bool
	width font.Length
	style []vg.Length
}

type colorBar struct {
	show     bool
	gradient colorgrad.Gradient
	min, max float64
	position positionType
}

// parameters to contour plot
func (plt *plotParameters) Contour(x, y, z *mat.Dense, options ...func(*contourOptions)) {
	// default options
	plt.contourOptions = contourOptions{
		nLevels: 10,
		lineSettings: lineSettings{
			style: Solid,
			width: vg.Points(1),
		},
	}

	// apply additional options
	for _, option := range options {
		option(&plt.contourOptions)
	}
	plt.colorBar = plt.contourOptions.colorBar

	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	var p palette.Palette
	if plt.contourOptions.gradient != (colorgrad.Gradient{}) {
		// add colormap and make a contour plotter
		p = &colorsGradient{colorList: plt.contourOptions.gradient.Colors(uint(plt.contourOptions.nLevels))}
	}

	levels := Linspace(mat.Min(z), mat.Max(z), plt.contourOptions.nLevels)
	c := plotter.NewContour(m, levels, p)
	c.LineStyles = []draw.LineStyle{{
		Color:  color.Black,
		Width:  plt.contourOptions.lineSettings.width,
		Dashes: plt.contourOptions.lineSettings.style,
	}}

	// add the plotters to the plot
	plt.plot.Add(c)

	if plt.colorBar.show {
		// get min and max values
		plt.colorBar.min = c.Min
		plt.colorBar.max = c.Max
	}
}

// parameters to contourf plot
func (plt *plotParameters) ContourF(x, y, z *mat.Dense, options ...func(*contourOptions)) {
	// default options
	plt.contourOptions = contourOptions{
		nLevels:  10,
		gradient: colorgrad.Viridis(),
		lineSettings: lineSettings{
			style: Solid,
			width: vg.Points(1),
		},
		colorBar: colorBar{
			gradient: colorgrad.Viridis(),
		},
	}

	// apply additional options
	for _, option := range options {
		option(&plt.contourOptions)
	}
	plt.colorBar = plt.contourOptions.colorBar

	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	// add colormap and make a heatmap plotter
	p := colorsGradient{colorList: plt.contourOptions.gradient.Colors(uint(plt.contourOptions.nLevels))}
	raster := plotter.NewHeatMap(m, &p)
	raster.Rasterized = true

	// add the plotters to the plot
	plt.plot.Add(raster)

	if plt.colorBar.show {
		// get min and max values
		plt.colorBar.min = raster.Min
		plt.colorBar.max = raster.Max
	}

	if plt.contourOptions.lineSettings.show {
		// add contour lines to contourf
		levels := Linspace(mat.Min(z), mat.Max(z), plt.contourOptions.nLevels)
		c := plotter.NewContour(m, levels, nil)
		c.LineStyles = []draw.LineStyle{{
			Color:  Black,
			Width:  plt.contourOptions.lineSettings.width,
			Dashes: plt.contourOptions.lineSettings.style,
		}}

		// add the plotters to the plot
		plt.plot.Add(c)
	}
}

func WithContourLevels(levels int) func(*contourOptions) {
	return func(co *contourOptions) {
		co.nLevels = levels
	}
}

func WithContourGradient(gradient colorgrad.Gradient) func(*contourOptions) {
	return func(co *contourOptions) {
		co.gradient = gradient
		co.colorBar.gradient = gradient
	}
}

func WithContourLines() func(*contourOptions) {
	return func(co *contourOptions) {
		co.lineSettings.show = true
	}
}

func WithContourLineWidth(width float64) func(*contourOptions) {
	return func(co *contourOptions) {
		co.lineSettings.width = vg.Points(width)
	}
}

func WithContourLineStyle(style lineStyleType) func(*contourOptions) {
	return func(co *contourOptions) {
		co.lineSettings.style = style
	}
}

func WithContourColorbar(position positionType) func(*contourOptions) {
	return func(co *contourOptions) {
		if co.gradient == (colorgrad.Gradient{}) {
			co.colorBar.show = false
		} else {
			co.colorBar.show = true
			co.colorBar.position = position
		}
	}
}
