package plot

import (
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
	plt.x = append(plt.x, x)
	plt.y = append(plt.y, y)
	plt.plot_name = "lineplot"
}

// parameters to heatmap plot
func (plt *plotParameters) HeatMap(x_dense, y_dense, z_dense *mat.Dense, n_levels int, gradient colorgrad.Gradient) {
	plt.x_dense = x_dense
	plt.y_dense = y_dense
	plt.z_dense = z_dense
	plt.n_levels = n_levels
	plt.gradient = gradient
	plt.plot_name = "heatmap"
}

// parameters to contour plot
func (plt *plotParameters) Contour(x_dense, y_dense, z_dense *mat.Dense, n_levels int, gradient colorgrad.Gradient) {
	plt.x_dense = x_dense
	plt.y_dense = y_dense
	plt.z_dense = z_dense
	plt.n_levels = n_levels
	plt.gradient = gradient
	plt.plot_name = "contour"
}

// parameters to scatter plot
func (plt *plotParameters) Scatter(x_values, y_values, z_values []float64, gradient colorgrad.Gradient) {
	plt.x_values = x_values
	plt.y_values = y_values
	plt.z_values = z_values
	plt.gradient = gradient
	plt.plot_name = "scatter"
}

// generate plot and save it to file
func (plt *plotParameters) Save(name string) {
	// create a new plot, set its title and axis labels
	p := plot.New()
	p.Title.Text = plt.title
	p.X.Label.Text = plt.xlabel
	p.Y.Label.Text = plt.ylabel

	switch plt.plot_name {
	case "lineplot":
		plt.linePlot(p) // make a line plotter
	case "heatmap":
		plt.heatMapPlot(p) // make a heatmap plotter
	case "contour":
		plt.contourPlot(p) // make a contour plotter
	case "scatter":
		plt.scatterPlot(p) // make a scatter plotter
	}

	// save the plot to a PNG file.
	err := p.Save(font.Length(plt.xwidth)*vg.Centimeter, font.Length(plt.ywidth)*vg.Centimeter, name)
	if err != nil {
		log.Panic(err)
	}
}

func (plt *plotParameters) linePlot(p *plot.Plot) {
	// draw a grid behind the data
	p.Add(plotter.NewGrid())

	// various plots to the figure
	lines := []*plotter.Line{}
	for nplot := 0; nplot < len(plt.x); nplot++ {
		pts := make(plotter.XYs, len(plt.x[nplot]))
		for j := range pts {
			pts[j].X = plt.x[nplot][j]
			pts[j].Y = plt.y[nplot][j]
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
	m := unitGrid{x: plt.x_dense, y: plt.y_dense, Data: plt.z_dense}

	// add colormap and make a heatmap plotter
	pal := colorsGradient{ColorList: plt.gradient.Colors(uint(plt.n_levels))}
	raster := plotter.NewHeatMap(m, pal)
	raster.Rasterized = true
	p.Add(raster)
}

func (plt *plotParameters) contourPlot(p *plot.Plot) {
	// prepare data to plot
	m := unitGrid{x: plt.x_dense, y: plt.y_dense, Data: plt.z_dense}

	// add colormap and make a contour plotter
	pal := colorsGradient{ColorList: plt.gradient.Colors(uint(plt.n_levels))}
	levels := Linspace(mat.Min(plt.z_dense), mat.Max(plt.z_dense), plt.n_levels)
	c := plotter.NewContour(m, levels, pal)
	p.Add(c)
}

func (plt *plotParameters) scatterPlot(p *plot.Plot) {
	p.Add(plotter.NewGrid())

	// prepare data to plot
	pts := make(plotter.XYZs, len(plt.x_values))
	for i := range pts {
		pts[i].X = plt.x_values[i]
		pts[i].Y = plt.y_values[i]
		pts[i].Z = plt.z_values[i]
	}

	// add colormap and make a scatter plotter
	sc, err := plotter.NewScatter(pts)
	if err != nil {
		log.Panic(err)
	}

	// specify style and color for individual points.
	sc.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		colors := plt.gradient.Colors(uint(len(plt.z_values)))
		return draw.GlyphStyle{Color: colors[i], Radius: vg.Points(3), Shape: draw.CircleGlyph{}}
	}
	p.Add(sc)
}

// size of the saved figure
func (plt *plotParameters) FigSize(xwidth, ywidth int) {
	plt.xwidth = xwidth
	plt.ywidth = ywidth
}

// title for all plots
func (plt *plotParameters) Title(str string) {
	plt.title = str
}

// xlabel for all plots
func (plt *plotParameters) XLabel(strx string) {
	plt.xlabel = strx
}

// ylabel for all plots
func (plt *plotParameters) YLabel(stry string) {
	plt.ylabel = stry
}

// legend mainly used in lines plots
func (plt *plotParameters) Legend(str ...string) {
	plt.legend = append(plt.legend, str...)
}
