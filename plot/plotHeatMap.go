package plot

import (
	"image/color"
	"log"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type plotHeatMap struct {
	x, y, z_values *mat.Dense
	gradient       colorgrad.Gradient
	n_levels       int
	title          string
	xlabel         string
	ylabel         string
	xwidth         int
	ywidth         int
}

type unitGrid struct {
	x, y, Data *mat.Dense
}

func (g unitGrid) Dims() (c, r int)   { r, c = g.Data.Dims(); return c, r }
func (g unitGrid) Z(c, r int) float64 { return g.Data.At(r, c) }
func (g unitGrid) X(c int) float64    { return g.x.At(0, c) }
func (g unitGrid) Y(r int) float64    { return g.y.At(0, r) }

type colorsGradient struct {
	ColorList []color.Color
}

func (g colorsGradient) Colors() []color.Color { return g.ColorList }

func NewHeatMap() plotHeatMap {
	return plotHeatMap{}
}

func (plt *plotHeatMap) Save(name string) {
	// create a new plot, set its title and axis labels
	p := plot.New()
	p.Title.Text = plt.title
	p.X.Label.Text = plt.xlabel
	p.Y.Label.Text = plt.ylabel

	// prepare data to plot
	m := unitGrid{x: plt.x, y: plt.x, Data: plt.z_values}

	// add colormap and make a heatmap plotter
	pal := colorsGradient{ColorList: plt.gradient.Colors(uint(plt.n_levels))}
	raster := plotter.NewHeatMap(m, pal)
	raster.Rasterized = true
	p.Add(raster)

	// save the plot to a PNG file.
	err := p.Save(font.Length(plt.xwidth)*vg.Centimeter, font.Length(plt.ywidth)*vg.Centimeter, name)
	if err != nil {
		log.Panic(err)
	}
}

func (plt *plotHeatMap) HeatMap(x, y, z_values *mat.Dense, n_levels int, gradient colorgrad.Gradient) {
	plt.x = x
	plt.y = y
	plt.z_values = z_values
	plt.n_levels = n_levels
	plt.gradient = gradient
}

func (plt *plotHeatMap) FigSize(xwidth, ywidth int) {
	plt.xwidth = xwidth
	plt.ywidth = ywidth
}

func (plt *plotHeatMap) Title(str string) {
	plt.title = str
}

func (plt *plotHeatMap) XLabel(strx string) {
	plt.xlabel = strx
}

func (plt *plotHeatMap) YLabel(stry string) {
	plt.ylabel = stry
}
