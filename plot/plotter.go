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

func NewPlot() plotParameters {
	return plotParameters{}
}

// parameters to lines plots
func (plt *plotParameters) Plot(x []float64, y []float64) {
	plt.plotData.x = append(plt.plotData.x, x)
	plt.plotData.y = append(plt.plotData.y, y)
	plt.plotName = "lineplot"
}

// parameters to heatmap plot
func (plt *plotParameters) HeatMap(x, y, z *mat.Dense, n_levels int, gradient colorgrad.Gradient) {
	plt.contourData.x = x
	plt.contourData.y = y
	plt.contourData.z = z
	plt.n_levels = n_levels
	plt.gradient = gradient
	plt.plotName = "heatmap"
}

// parameters to contour plot
func (plt *plotParameters) Contour(x, y, z *mat.Dense, n_levels int, gradient colorgrad.Gradient) {
	plt.contourData.x = x
	plt.contourData.y = y
	plt.contourData.z = z
	plt.n_levels = n_levels
	plt.gradient = gradient
	plt.plotName = "contour"
}

// parameters to image plot
func (plt *plotParameters) ImShow(image image.Image) {
	plt.imageData = image
	plt.plotName = "image"
}

// parameters to scatter plot
func (plt *plotParameters) Scatter(x, y, z []float64, gradient colorgrad.Gradient) {
	plt.scatterData.x = x
	plt.scatterData.y = y
	plt.scatterData.z = z
	plt.gradient = gradient
	plt.plotName = "scatter"
}

// generate plot and save it to file
func (plt *plotParameters) Save(name string) {
	// create a new plot, set its title and axis labels
	p := plot.New()
	p.Title.Text = plt.title
	p.X.Label.Text = plt.axisLabel.xlabel
	p.Y.Label.Text = plt.axisLabel.ylabel

	// set the axis limits
	if plt.axisLimit.useXLim {
		p.X.Min = plt.axisLimit.xmin
		p.X.Max = plt.axisLimit.xmax
	}
	if plt.axisLimit.useYLim {
		p.Y.Min = plt.axisLimit.ymin
		p.Y.Max = plt.axisLimit.ymax
	}

	switch plt.plotName {
	case "lineplot":
		plt.linePlot(p) // make a line plotter
	case "heatmap":
		plt.heatMapPlot(p) // make a heatmap plotter
	case "contour":
		plt.contourPlot(p) // make a contour plotter
	case "image":
		plt.imagePlot(p) // make a image plotter
	case "scatter":
		plt.scatterPlot(p) // make a scatter plotter
	}

	// save the plot to a PNG file.
	xwdith := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywdith := font.Length(plt.figSize.ywidth) * vg.Centimeter
	err := p.Save(xwdith, ywdith, name)
	if err != nil {
		log.Panic(err)
	}
}

func (plt *plotParameters) linePlot(p *plot.Plot) {
	// draw a grid behind the data
	p.Add(plotter.NewGrid())

	// various plots to the figure
	lines := []*plotter.Line{}
	for nplot := 0; nplot < len(plt.plotData.x); nplot++ {
		pts := make(plotter.XYs, len(plt.plotData.x[nplot]))
		for j := range pts {
			pts[j].X = plt.plotData.x[nplot][j]
			pts[j].Y = plt.plotData.y[nplot][j]
		}

		// make a line plotter with points and set its style.
		line, _, _ := plotter.NewLinePoints(pts)
		line.Color = colors[nplot]
		line.LineStyle.Width = vg.Points(1.5)
		lines = append(lines, line)

		// add the plotters to the plot, with a legend
		p.Add(line)
	}

	// legend style
	for i, legend := range plt.legend {
		p.Legend.Add(legend, lines[i])
		p.Legend.XOffs = -5. * vg.Millimeter
		p.Legend.Padding = vg.Millimeter
	}
}

func (plt *plotParameters) heatMapPlot(p *plot.Plot) {
	// prepare data to plot
	m := unitGrid{x: plt.contourData.x, y: plt.contourData.y, Data: plt.contourData.z}

	// add colormap and make a heatmap plotter
	pal := colorsGradient{ColorList: plt.gradient.Colors(uint(plt.n_levels))}
	raster := plotter.NewHeatMap(m, pal)
	raster.Rasterized = true
	p.Add(raster)
}

func (plt *plotParameters) contourPlot(p *plot.Plot) {
	// prepare data to plot
	m := unitGrid{x: plt.contourData.x, y: plt.contourData.y, Data: plt.contourData.z}

	// add colormap and make a contour plotter
	pal := colorsGradient{ColorList: plt.gradient.Colors(uint(plt.n_levels))}
	levels := Linspace(mat.Min(plt.contourData.z), mat.Max(plt.contourData.z), plt.n_levels)
	c := plotter.NewContour(m, levels, pal)
	p.Add(c)
}

func (plt *plotParameters) imagePlot(p *plot.Plot) {
	// prepare data to plot
	b := plt.imageData.Bounds()
	xmin := float64(b.Min.X)
	ymin := float64(b.Min.Y)
	xmax := float64(b.Max.X)
	ymax := float64(b.Max.Y)

	// add and make a image plotter
	img := plotter.NewImage(plt.imageData, xmin, ymin, xmax, ymax)
	p.Add(img)
}

func (plt *plotParameters) scatterPlot(p *plot.Plot) {
	p.Add(plotter.NewGrid())

	// prepare data to plot
	pts := make(plotter.XYZs, len(plt.scatterData.x))
	for i := range pts {
		pts[i].X = plt.scatterData.x[i]
		pts[i].Y = plt.scatterData.y[i]
		pts[i].Z = plt.scatterData.z[i]
	}

	// add colormap and make a scatter plotter
	sc, err := plotter.NewScatter(pts)
	if err != nil {
		log.Panic(err)
	}

	// specify style and color for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		colors := plt.gradient.Colors(uint(len(plt.scatterData.z)))
		return draw.GlyphStyle{Color: colors[i], Radius: vg.Points(3), Shape: draw.CircleGlyph{}}
	}
	p.Add(sc)
}

// size of the saved figure
func (plt *plotParameters) FigSize(xwidth, ywidth int) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
}

// title for all plots
func (plt *plotParameters) Title(str string) {
	plt.title = str
}

// xlabel for all plots
func (plt *plotParameters) XLabel(xlabel string) {
	plt.axisLabel.xlabel = xlabel
}

// ylabel for all plots
func (plt *plotParameters) YLabel(ylabel string) {
	plt.axisLabel.ylabel = ylabel
}

// legend mainly used in lines plots
func (plt *plotParameters) Legend(str ...string) {
	plt.legend = append(plt.legend, str...)
}

// set the x-axis vies limits
func (plt *plotParameters) XLim(xmin, xmax float64) {
	if xmin < xmax {
		plt.axisLimit.xmin = xmin
		plt.axisLimit.xmax = xmax
		plt.axisLimit.useXLim = true
	}
}

// set the x-axis vies limits
func (plt *plotParameters) YLim(ymin, ymax float64) {
	if ymin < ymax {
		plt.axisLimit.ymin = ymin
		plt.axisLimit.ymax = ymax
		plt.axisLimit.useYLim = true
	}
}
