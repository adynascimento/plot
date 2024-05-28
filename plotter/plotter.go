package plotter

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type PlotterInterface interface {
	Plot(x, y []float64, options ...func(*lineOptions))
	Contour(x, y, z *mat.Dense, options ...func(*contourOptions))
	ContourF(x, y, z *mat.Dense, options ...func(*contourOptions))
	Scatter(x, y, z []float64, params ...interface{})
	ImShow(image image.Image)
	Title(str string)
	XLabel(xlabel string)
	YLabel(ylabel string)
	Legend(str ...string)
	XLim(xmin, xmax float64)
	YLim(ymin, ymax float64)
	Grid()
}

type Plot interface {
	PlotterInterface
	Save(name string)
	FigSize(xwidth, ywidth int)
}

func NewPlot() Plot {
	return &PlotParameters{
		plot: plot.New(),
		lineOptions: lineOptions{
			usedColors: make(map[color.Color]bool),
		},
	}
}

// parameters to lines plots
func (plt *PlotParameters) Plot(x, y []float64, options ...func(*lineOptions)) {
	var thumbs []plot.Thumbnailer
	var plotters []plot.Plotter

	// default options
	plt.lineOptions.params = params{
		lineWidth:     vg.Points(1.5),
		lineStyle:     []vg.Length{},
		markerSize:    vg.Points(3),
		markerSpacing: 1,
	}

	// apply additional options
	for _, option := range options {
		option(&plt.lineOptions)
	}

	// automatic color assignment
	if plt.lineOptions.color == nil {
		plt.lineOptions.color = plt.nextColor()
	}

	// various plots to the figure
	pts := make(plotter.XYs, len(x))
	for j := range pts {
		pts[j].X = x[j]
		pts[j].Y = y[j]
	}

	// make a line plotter and set its style.
	line, _ := plotter.NewLine(pts)
	line.Color = plt.lineOptions.color
	line.LineStyle.Width = plt.lineOptions.lineWidth
	line.LineStyle.Dashes = plt.lineOptions.lineStyle

	thumbs = append(thumbs, line)
	plotters = append(plotters, line)

	// add markers to line plotter
	if plt.lineOptions.marker != nil {
		scatter := plt.addMarkers(pts)

		thumbs = append(thumbs, scatter)
		plotters = append(plotters, scatter)
	}

	// thumbs for the legends
	plt.legends = append(plt.legends, thumbs)

	// add the plotters to the plot
	plt.plot.Add(plotters...)
}

// parameters to contour plot
func (plt *PlotParameters) Contour(x, y, z *mat.Dense, options ...func(*contourOptions)) {
	// default options
	plt.contourOptions = contourOptions{
		nLevels:   10,
		gradient:  colorgrad.Viridis(),
		lineWidth: vg.Points(1),
		lineStyle: []vg.Length{},
	}

	// apply additional options
	for _, option := range options {
		option(&plt.contourOptions)
	}

	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	// add colormap and make a contour plotter
	pal := colorsGradient{
		ColorList: plt.contourOptions.gradient.Colors(uint(plt.contourOptions.nLevels)),
	}
	levels := Linspace(mat.Min(z), mat.Max(z), plt.contourOptions.nLevels)
	c := plotter.NewContour(m, levels, pal)
	c.LineStyles = []draw.LineStyle{{
		Color:  color.Black,
		Width:  plt.contourOptions.lineWidth,
		Dashes: plt.contourOptions.lineStyle,
	}}

	// add the plotters to the plot
	plt.plot.Add(c)
}

// parameters to contourf plot
func (plt *PlotParameters) ContourF(x, y, z *mat.Dense, options ...func(*contourOptions)) {
	// default options
	plt.contourOptions = contourOptions{
		nLevels:  10,
		gradient: colorgrad.Viridis(),
	}

	// apply additional options
	for _, option := range options {
		option(&plt.contourOptions)
	}

	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	// add colormap and make a heatmap plotter
	pal := colorsGradient{
		ColorList: plt.contourOptions.gradient.Colors(uint(plt.contourOptions.nLevels)),
	}
	raster := plotter.NewHeatMap(m, pal)
	raster.Rasterized = true

	// add the plotters to the plot
	plt.plot.Add(raster)

	if plt.contourOptions.addLines {
		// add contour lines to contourf
		black, _ := colorgrad.NewGradient().Colors(color.Black).Build()
		options = append(options, WithGradient(black))
		plt.Contour(x, y, z, options...)
	}
}

// parameters to scatter plot
func (plt *PlotParameters) Scatter(x, y, z []float64, params ...interface{}) {
	// default values
	gradient := colorgrad.Viridis()
	for _, param := range params {
		switch v := param.(type) {
		case colorgrad.Gradient:
			gradient = v
		default:
			fmt.Printf("unknown options for scatter plot, type: %T\n", v)
		}
	}

	// prepare data to plot
	pts := make(plotter.XYZs, len(x))
	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
		pts[i].Z = z[i]
	}

	// add colormap and make a scatter plotter
	sc, err := plotter.NewScatter(pts)
	if err != nil {
		log.Panic(err)
	}

	// specify style and color for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		colors := gradient.Colors(uint(len(z)))
		return draw.GlyphStyle{Color: colors[i], Radius: vg.Points(3), Shape: draw.CircleGlyph{}}
	}
	plt.plot.Add(sc)
}

// parameters to image plot
func (plt *PlotParameters) ImShow(image image.Image) {
	// prepare data to plot
	b := image.Bounds()
	xmin := float64(b.Min.X)
	ymin := float64(b.Min.Y)
	xmax := float64(b.Max.X)
	ymax := float64(b.Max.Y)

	// add and make a image plotter
	img := plotter.NewImage(image, xmin, ymin, xmax, ymax)
	plt.plot.Add(img)
}

// save the plot to an image file
func (plt *PlotParameters) Save(name string) {
	// save the plot to a PNG file.
	xwdith := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywdith := font.Length(plt.figSize.ywidth) * vg.Centimeter
	err := plt.plot.Save(xwdith, ywdith, name)
	if err != nil {
		log.Panic(err)
	}
}

// size of the saved figure
func (plt *PlotParameters) FigSize(xwidth, ywidth int) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
}

// title for all plots
func (plt *PlotParameters) Title(title string) {
	plt.plot.Title.Text = title
}

// xlabel for all plots
func (plt *PlotParameters) XLabel(xlabel string) {
	plt.plot.X.Label.Text = xlabel
}

// ylabel for all plots
func (plt *PlotParameters) YLabel(ylabel string) {
	plt.plot.Y.Label.Text = ylabel
}

// legend mainly used in lines plots
func (plt *PlotParameters) Legend(str ...string) {
	// legend style
	for i, legend := range str {
		plt.plot.Legend.Add(legend, plt.legends[i]...)
		plt.plot.Legend.XOffs = -5. * vg.Millimeter
		plt.plot.Legend.Padding = vg.Millimeter
	}
}

// set the x-axis vies limits
func (plt *PlotParameters) XLim(xmin, xmax float64) {
	if xmin < xmax {
		plt.plot.X.Min = xmin
		plt.plot.X.Max = xmax
	}
}

// set the x-axis vies limits
func (plt *PlotParameters) YLim(ymin, ymax float64) {
	if ymin < ymax {
		plt.plot.Y.Min = ymin
		plt.plot.Y.Max = ymax
	}
}

// draw grid with both vertical and horizontal lines
func (plt *PlotParameters) Grid() {
	plt.plot.Add(plotter.NewGrid())
}
