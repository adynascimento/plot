package plot

import (
	"log"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type plotScatter struct {
	x, y, z_values []float64
	gradient       colorgrad.Gradient
	title          string
	xlabel         string
	ylabel         string
	xwidth         int
	ywidth         int
}

func NewScatter() plotScatter {
	return plotScatter{}
}

func (plt *plotScatter) Save(name string) {
	// create a new plot, set its title and axis labels
	p := plot.New()
	p.Title.Text = plt.title
	p.X.Label.Text = plt.xlabel
	p.Y.Label.Text = plt.ylabel
	p.Add(plotter.NewGrid())

	// prepare data to plot
	pts := make(plotter.XYZs, len(plt.x))
	for i := range pts {
		pts[i].X = plt.x[i]
		pts[i].Y = plt.y[i]
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

	// save the plot to a PNG file.
	err = p.Save(font.Length(plt.xwidth)*vg.Centimeter, font.Length(plt.ywidth)*vg.Centimeter, name)
	if err != nil {
		log.Panic(err)
	}
}

func (plt *plotScatter) Scatter(x, y, z_values []float64, gradient colorgrad.Gradient) {
	plt.x = x
	plt.y = y
	plt.z_values = z_values
	plt.gradient = gradient
}

func (plt *plotScatter) FigSize(xwidth, ywidth int) {
	plt.xwidth = xwidth
	plt.ywidth = ywidth
}

func (plt *plotScatter) Title(str string) {
	plt.title = str
}

func (plt *plotScatter) XLabel(strx string) {
	plt.xlabel = strx
}

func (plt *plotScatter) YLabel(stry string) {
	plt.ylabel = stry
}
