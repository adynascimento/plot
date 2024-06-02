package plotter

import (
	"image"
	"image/color"
	"log"
	"os"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type PlotterInterface interface {
	Plot(x, y []float64, options ...func(*lineOptions))
	Contour(x, y, z *mat.Dense, options ...func(*contourOptions))
	ContourF(x, y, z *mat.Dense, options ...func(*contourOptions))
	Scatter(x, y, z []float64, options ...func(*scatterOptions))
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
	return &plotParameters{
		plot: plot.New(),
		lineOptions: lineOptions{
			usedColors: make(map[color.Color]bool),
		},
	}
}

// parameters to lines plots
func (plt *plotParameters) Plot(x, y []float64, options ...func(*lineOptions)) {
	var thumbs []plot.Thumbnailer
	var plotters []plot.Plotter

	// default options
	plt.lineOptions.params = params{
		lineStyle:     Solid,
		lineWidth:     vg.Points(1.5),
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

// parameters to scatter plot
func (plt *plotParameters) Scatter(x, y, z []float64, options ...func(*scatterOptions)) {
	// default options
	plt.scatterOptions = scatterOptions{
		color:      Blue,
		marker:     Circle,
		markerSize: vg.Points(3),
	}

	// apply additional options
	for _, option := range options {
		option(&plt.scatterOptions)
	}
	plt.colorBar = plt.scatterOptions.colorBar

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
		Color:  plt.scatterOptions.color,
		Radius: plt.scatterOptions.markerSize,
		Shape:  plt.scatterOptions.marker,
	}

	if plt.scatterOptions.gradient != (colorgrad.Gradient{}) {
		// specify style and color for individual points.
		sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
			colors := plt.scatterOptions.gradient.Colors(uint(len(z)))
			return draw.GlyphStyle{Color: colors[i], Radius: plt.scatterOptions.markerSize,
				Shape: plt.scatterOptions.marker}
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

// parameters to image plot
func (plt *plotParameters) ImShow(image image.Image) {
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
func (plt *plotParameters) Save(file string) {
	// save the plot to a PNG file.
	xwidth := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywidth := font.Length(plt.figSize.ywidth) * vg.Centimeter

	// new image canvas
	img := vgimg.New(xwidth, ywidth)

	// draw the plot
	plt.plot.Draw(draw.Canvas{
		Canvas: draw.New(img),
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: xwidth, Y: ywidth},
		},
	})

	// add colorbar to plot
	if plt.colorBar.show {
		switch plt.colorBar.position {
		case Vertical:
			img = plt.drawVerticalColorBar(xwidth, ywidth)
		case Horizontal:
			img = plt.drawHorizontalColorBar(xwidth, ywidth)
		}
	}

	// save the image to a file
	w, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		panic(err)
	}
}

// size of the saved figure
func (plt *plotParameters) FigSize(xwidth, ywidth int) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
}

// title for all plots
func (plt *plotParameters) Title(title string) {
	plt.plot.Title.Text = title
}

// xlabel for all plots
func (plt *plotParameters) XLabel(xlabel string) {
	plt.plot.X.Label.Text = xlabel
}

// ylabel for all plots
func (plt *plotParameters) YLabel(ylabel string) {
	plt.plot.Y.Label.Text = ylabel
}

// legend mainly used in lines plots
func (plt *plotParameters) Legend(str ...string) {
	// legend style
	for i, legend := range str {
		plt.plot.Legend.Add(legend, plt.legends[i]...)
		plt.plot.Legend.XOffs = -5. * vg.Millimeter
		plt.plot.Legend.Padding = vg.Millimeter
	}
}

// set the x-axis vies limits
func (plt *plotParameters) XLim(xmin, xmax float64) {
	if xmin < xmax {
		plt.plot.X.Min = xmin
		plt.plot.X.Max = xmax
	}
}

// set the x-axis vies limits
func (plt *plotParameters) YLim(ymin, ymax float64) {
	if ymin < ymax {
		plt.plot.Y.Min = ymin
		plt.plot.Y.Max = ymax
	}
}

// draw grid with both vertical and horizontal lines
func (plt *plotParameters) Grid() {
	plt.plot.Add(plotter.NewGrid())
}
