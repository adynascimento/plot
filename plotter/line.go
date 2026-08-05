package plotter

import (
	"image/color"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type (
	colorType     color.Color
	lineStyleType []vg.Length
	markerType    draw.GlyphDrawer
)

var (
	Black   colorType = color.RGBA{000, 000, 000, 255}
	Red     colorType = color.RGBA{255, 000, 000, 255}
	Blue    colorType = color.RGBA{000, 000, 255, 255}
	Green   colorType = color.RGBA{000, 128, 000, 255}
	Cyan    colorType = color.RGBA{000, 255, 255, 255}
	Magenta colorType = color.RGBA{255, 000, 255, 255}
	Orange  colorType = color.RGBA{255, 165, 000, 255}
	Purple  colorType = color.RGBA{128, 000, 128, 255}
	Yellow  colorType = color.RGBA{255, 255, 000, 255}
	Gray    colorType = color.RGBA{128, 128, 128, 255}
)

// default colors to line plot
var colors = []color.Color{
	Black, Red, Blue,
	Green, Cyan, Magenta,
	Orange, Purple, Yellow,
}

var (
	Solid      lineStyleType = []vg.Length{}             // "-" solid line
	Dashed     lineStyleType = []vg.Length{vg.Points(5)} // "--" dashed line
	Dotted     lineStyleType = []vg.Length{vg.Points(1)} // ":" dotted line
	DashDotted lineStyleType = []vg.Length{vg.Points(6), vg.Points(3),
		vg.Points(1), vg.Points(3)} // "-." dash-dotted line
)

var (
	Circle    markerType = draw.CircleGlyph{}  // filled circle
	Square    markerType = draw.BoxGlyph{}     // filled square
	Triangle  markerType = draw.PyramidGlyph{} // filled triangle
	PlusSign  markerType = draw.PlusGlyph{}    // plus sign
	CrossSign markerType = draw.CrossGlyph{}   // x sign
)

type lineOptions struct {
	params
	colorIndex int
	usedColors map[color.Color]bool
}

type params struct {
	color         color.Color
	lineStyle     []vg.Length
	lineWidth     font.Length
	marker        draw.GlyphDrawer
	markerSize    font.Length
	markerSpacing int
}

// parameters to line plots
func (plt *plotParameters) Plot(x, y []float64, options ...func(*lineOptions)) {
	var thumbs []plot.Thumbnailer
	var plotters []plot.Plotter

	// default options
	plt.line.params = params{
		lineStyle:     Solid,
		lineWidth:     vg.Points(1.5),
		markerSize:    vg.Points(3),
		markerSpacing: 1,
	}

	// apply additional options
	for _, option := range options {
		option(&plt.line)
	}

	// automatic color assignment
	if plt.line.color == nil {
		plt.line.color = plt.nextColor()
	}

	// various plots to the figure
	pts := make(plotter.XYs, len(x))
	for j := range pts {
		pts[j].X = x[j]
		pts[j].Y = y[j]
	}

	// make a line plotter and set its style.
	line, _ := plotter.NewLine(pts)
	line.Color = plt.line.color
	line.LineStyle.Width = plt.line.lineWidth
	line.LineStyle.Dashes = plt.line.lineStyle

	thumbs = append(thumbs, line)
	plotters = append(plotters, line)

	// add markers to line plotter
	if plt.line.marker != nil {
		scatter := plt.addMarkers(pts)

		thumbs = append(thumbs, scatter)
		plotters = append(plotters, scatter)
	}

	// thumbs for the legends
	plt.legends = append(plt.legends, thumbs)

	// add the plotters to the plot
	plt.plot.Add(plotters...)
}

func WithLineColor(color colorType) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.color = color
		lo.usedColors[color] = true
	}
}

func WithLineWidth(width float64) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.lineWidth = vg.Points(width)
	}
}

func WithLineStyle(style lineStyleType) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.lineStyle = style
	}
}

func WithLineMarker(marker markerType) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.marker = marker
	}
}

func WithLineMarkerSize(size float64) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.markerSize = vg.Points(size)
	}
}

func WithLineMarkerSpacing(spacing int) func(*lineOptions) {
	return func(lo *lineOptions) {
		lo.markerSpacing = spacing
	}
}

// get the next color in the sequence
func (plt *plotParameters) nextColor() color.Color {
	for {
		c := colors[plt.line.colorIndex]
		plt.line.colorIndex = (plt.line.colorIndex + 1) % len(colors)
		if !plt.line.usedColors[c] {
			plt.line.usedColors[c] = true
			return c
		}
	}
}

// add markers to line plotter
func (plt *plotParameters) addMarkers(pts plotter.XYs) *plotter.Scatter {
	spacing := plt.line.markerSpacing
	spacedPts := make(plotter.XYs, (len(pts)+spacing-1)/spacing)
	for i := 0; i < len(spacedPts); i++ {
		spacedPts[i] = pts[i*spacing]
	}
	scatter, _ := plotter.NewScatter(spacedPts)
	scatter.GlyphStyle.Shape = plt.line.marker
	scatter.GlyphStyle.Radius = plt.line.markerSize
	scatter.Color = plt.line.color

	return scatter
}
