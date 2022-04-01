package plot

import (
	"log"

	ngo "plot/numeric"
	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

type plotContour struct {
	x, y, z_values *mat.Dense
	gradient       colorgrad.Gradient
	n_levels       int
	title          string
	xlabel         string
	ylabel         string
	xwidth         int
	ywidth         int
}

func NewContour() plotContour {
	return plotContour{}
}

func (plt *plotContour) Save(name string) {
	// create a new plot, set its title and axis labels
	p := plot.New()
	p.Title.Text = plt.title
	p.X.Label.Text = plt.xlabel
	p.Y.Label.Text = plt.ylabel

	// prepare data to plot
	m := unitGrid{x: plt.x, y: plt.x, Data: plt.z_values}

	// add colormap and make a contour plotter
	pal := colorsGradient{ColorList: plt.gradient.Colors(uint(plt.n_levels))}
	levels := ngo.Linspace(mat.Min(plt.z_values), mat.Max(plt.z_values), plt.n_levels)
	c := plotter.NewContour(m, levels, pal)
	p.Add(c)
	
	// save the plot to a PNG file.
	err := p.Save(font.Length(plt.xwidth)*vg.Centimeter, font.Length(plt.ywidth)*vg.Centimeter, name)
	if err != nil {
		log.Panic(err)
	}
}

func (plt *plotContour) Contour(x, y, z_values *mat.Dense, n_levels int, gradient colorgrad.Gradient) {
	plt.x = x
	plt.y = y
	plt.z_values = z_values
	plt.n_levels = n_levels
	plt.gradient = gradient
}

func (plt *plotContour) FigSize(xwidth, ywidth int) {
	plt.xwidth = xwidth
	plt.ywidth = ywidth
}

func (plt *plotContour) Title(str string) {
	plt.title = str
}

func (plt *plotContour) XLabel(strx string) {
	plt.xlabel = strx
}

func (plt *plotContour) YLabel(stry string) {
	plt.ylabel = stry
}

