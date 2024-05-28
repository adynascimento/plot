package plotter

import (
	"image/color"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

// default colors to lines plot
var colors = []color.Color{
	color.RGBA{000, 000, 000, 255}, // black
	color.RGBA{255, 000, 000, 255}, // red
	color.RGBA{000, 000, 255, 255}, // blue
	color.RGBA{000, 128, 000, 255}, // green
	color.RGBA{000, 255, 255, 255}, // cyan
	color.RGBA{255, 000, 255, 255}, // magenta
	color.RGBA{255, 165, 000, 255}, // orange
	color.RGBA{128, 000, 128, 255}, // purple
	color.RGBA{255, 255, 000, 255}, // yellow
}

var colormap = map[string]color.Color{
	"black":   color.RGBA{000, 000, 000, 255},
	"red":     color.RGBA{255, 000, 000, 255},
	"blue":    color.RGBA{000, 000, 255, 255},
	"green":   color.RGBA{000, 128, 000, 255},
	"cyan":    color.RGBA{000, 255, 255, 255},
	"magenta": color.RGBA{255, 000, 255, 255},
	"orange":  color.RGBA{255, 165, 000, 255},
	"purple":  color.RGBA{128, 000, 128, 255},
	"yellow":  color.RGBA{255, 255, 000, 255},
}

var linestyle = map[string][]vg.Length{
	"--": {vg.Points(5)},                                           // dashed line
	":":  {vg.Points(1)},                                           // dotted line
	"-.": {vg.Points(6), vg.Points(3), vg.Points(1), vg.Points(3)}, // dash-dotted line
}

var markers = map[string]draw.GlyphDrawer{
	"o": draw.CircleGlyph{},  // solid circle
	"s": draw.BoxGlyph{},     // filled square
	"p": draw.PyramidGlyph{}, // filled triangle
	"+": draw.PlusGlyph{},    // plus sign
	"x": draw.CrossGlyph{},   // x sign
}

// get the next color in the sequence
func (plt *PlotParameters) nextColor() color.Color {
	for {
		c := colors[plt.lineOptions.colorIndex]
		plt.lineOptions.colorIndex = (plt.lineOptions.colorIndex + 1) % len(colors)
		if !plt.lineOptions.usedColors[c] {
			plt.lineOptions.usedColors[c] = true
			return c
		}
	}
}

// add markers to line plotter
func (plt *PlotParameters) addMarkers(pts plotter.XYs) *plotter.Scatter {
	spacing := plt.lineOptions.markerSpacing
	spacedPts := make(plotter.XYs, (len(pts)+spacing-1)/spacing)
	for i := 0; i < len(spacedPts); i++ {
		spacedPts[i] = pts[i*spacing]
	}
	scatter, _ := plotter.NewScatter(spacedPts)
	scatter.GlyphStyle.Shape = plt.lineOptions.marker
	scatter.GlyphStyle.Radius = plt.lineOptions.markerSize
	scatter.Color = plt.lineOptions.color

	return scatter
}

type lineOptions struct {
	params
	colorIndex int
	usedColors map[color.Color]bool
}

type params struct {
	color         color.Color
	lineWidth     font.Length
	lineStyle     []vg.Length
	marker        draw.GlyphDrawer
	markerSize    font.Length
	markerSpacing int
}

func WithLineColor(color string) func(*lineOptions) {
	return func(lo *lineOptions) {
		if c, exists := colormap[color]; exists {
			lo.color = c
			lo.usedColors[c] = true
		}
	}
}

func WithLineWidth(width float64) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.lineWidth = vg.Points(width)
	}
}

func WithLineStyle(style string) func(*lineOptions) {
	return func(lo *lineOptions) {
		if l, exists := linestyle[style]; exists {
			lo.lineStyle = l
		}
	}
}

func WithMarker(marker string) func(*lineOptions) {
	return func(lo *lineOptions) {
		if m, exists := markers[marker]; exists {
			lo.marker = m
		}
	}
}

func WithMarkerSize(size float64) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.markerSize = vg.Points(size)
	}
}

func WithMarkerSpacing(spacing int) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.markerSpacing = spacing
	}
}
