package plotter

import (
	"image/color"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type contourOptions struct {
	nLevels      int
	colormap     Colormap
	lineSettings lineSettings
	colorbar     colorbar
}

type lineSettings struct {
	show  bool
	width font.Length
	style []vg.Length
}

// parameters to contour plot
func (plt *plotParameters) Contour(x, y, z *mat.Dense, options ...func(*contourOptions)) {
	// default options
	plt.contour = contourOptions{
		nLevels: 10,
		lineSettings: lineSettings{
			style: Solid,
			width: vg.Points(1),
		},
	}

	// apply additional options
	for _, option := range options {
		option(&plt.contour)
	}
	plt.colorbar = plt.contour.colorbar

	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	var p palette.Palette
	if plt.contour.colormap != nil {
		// add colormap and make a contour plotter
		p = &colormapPalette{
			colorList: plt.contour.colormap.Colors(plt.contour.nLevels),
			colormap:  plt.contour.colormap,
			min:       mat.Min(z),
			max:       mat.Max(z),
		}
	}

	levels := Linspace(mat.Min(z), mat.Max(z), plt.contour.nLevels)
	c := plotter.NewContour(m, levels, p)
	c.LineStyles = []draw.LineStyle{{
		Color:  color.Black,
		Width:  plt.contour.lineSettings.width,
		Dashes: plt.contour.lineSettings.style,
	}}

	// add the plotters to the plot
	plt.plot.Add(c)

	if plt.colorbar.show {
		// get min and max values
		plt.colorbar.min = c.Min
		plt.colorbar.max = c.Max
	}
}

// parameters to contourf plot
func (plt *plotParameters) ContourF(x, y, z *mat.Dense, options ...func(*contourOptions)) {
	// default options
	plt.contour = contourOptions{
		nLevels:  10,
		colormap: Viridis,
		lineSettings: lineSettings{
			style: Solid,
			width: vg.Points(1),
		},
		colorbar: colorbar{
			colormap: Viridis,
		},
	}

	// apply additional options
	for _, option := range options {
		option(&plt.contour)
	}
	plt.colorbar = plt.contour.colorbar

	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	// add colormap and make a heatmap plotter
	p := colormapPalette{
		colorList: plt.contour.colormap.Colors(plt.contour.nLevels),
		colormap:  plt.contour.colormap,
		min:       mat.Min(z),
		max:       mat.Max(z),
	}
	raster := plotter.NewHeatMap(m, &p)
	raster.Rasterized = true

	// add the plotters to the plot
	plt.plot.Add(raster)

	if plt.colorbar.show {
		// get min and max values
		plt.colorbar.min = raster.Min
		plt.colorbar.max = raster.Max
	}

	if plt.contour.lineSettings.show {
		// add contour lines to contourf
		levels := Linspace(mat.Min(z), mat.Max(z), plt.contour.nLevels)
		c := plotter.NewContour(m, levels, nil)
		c.LineStyles = []draw.LineStyle{{
			Color:  Black,
			Width:  plt.contour.lineSettings.width,
			Dashes: plt.contour.lineSettings.style,
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

func WithContourColormap(cmap Colormap) func(*contourOptions) {
	return func(co *contourOptions) {
		co.colormap = cmap
		co.colorbar.colormap = cmap
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
		if co.colormap == nil {
			co.colorbar.show = false
		} else {
			co.colorbar.show = true
			co.colorbar.position = position
		}
	}
}
