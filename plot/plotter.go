package plot

import (
	"image"
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
	Plot(x []float64, y []float64)
	HeatMap(x, y, z *mat.Dense, n_levels int, gradient colorgrad.Gradient)
	Contour(x, y, z *mat.Dense, n_levels int, gradient colorgrad.Gradient)
	ImShow(image image.Image)
	Scatter(x, y, z []float64, gradient colorgrad.Gradient)
	Title(str string)
	XLabel(xlabel string)
	YLabel(ylabel string)
	Legend(str ...string)
	XLim(xmin, xmax float64)
	YLim(ymin, ymax float64)
}

type Plot interface {
	PlotterInterface
	Save(name string)
	FigSize(xwidth, ywidth int)
}

func NewPlot() Plot {
	return &plotParameters{
		plot: plot.New(),
	}
}

// parameters to lines plots
func (plt *plotParameters) Plot(x []float64, y []float64) {
	// draw a grid behind the data
	plt.plot.Add(plotter.NewGrid())

	// various plots to the figure
	pts := make(plotter.XYs, len(x))
	for j := range pts {
		pts[j].X = x[j]
		pts[j].Y = y[j]
	}

	// make a line plotter with points and set its style.
	line, _, _ := plotter.NewLinePoints(pts)
	line.LineStyle.Width = vg.Points(1.5)
	line.Color = colors[plt.nplot]

	plt.lines = append(plt.lines, line)
	plt.nplot++

	// add the plotters to the plot, with a legend
	plt.plot.Add(line)
}

// parameters to heatmap plot
func (plt *plotParameters) HeatMap(x, y, z *mat.Dense, nLevels int, gradient colorgrad.Gradient) {
	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	// add colormap and make a heatmap plotter
	pal := colorsGradient{ColorList: gradient.Colors(uint(nLevels))}
	raster := plotter.NewHeatMap(m, pal)
	raster.Rasterized = true
	plt.plot.Add(raster)
}

// parameters to contour plot
func (plt *plotParameters) Contour(x, y, z *mat.Dense, nLevels int, gradient colorgrad.Gradient) {
	// prepare data to plot
	m := unitGrid{x: x, y: y, Data: z}

	// add colormap and make a contour plotter
	pal := colorsGradient{ColorList: gradient.Colors(uint(nLevels))}
	levels := Linspace(mat.Min(z), mat.Max(z), nLevels)
	c := plotter.NewContour(m, levels, pal)
	plt.plot.Add(c)
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

// parameters to scatter plot
func (plt *plotParameters) Scatter(x, y, z []float64, gradient colorgrad.Gradient) {
	plt.plot.Add(plotter.NewGrid())

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

// save the plot to an image file
func (plt *plotParameters) Save(name string) {
	// save the plot to a PNG file.
	xwdith := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywdith := font.Length(plt.figSize.ywidth) * vg.Centimeter
	err := plt.plot.Save(xwdith, ywdith, name)
	if err != nil {
		log.Panic(err)
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
		plt.plot.Legend.Add(legend, plt.lines[i])
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
