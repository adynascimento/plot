package plot

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
)

// default colors to lines plot
var colors = []color.Color{
	color.RGBA{000, 000, 000, 255},
	color.RGBA{255, 000, 000, 255},
	color.RGBA{122, 195, 106, 255},
	color.RGBA{90, 155, 212, 255},
	color.RGBA{250, 167, 91, 255},
	color.RGBA{158, 103, 171, 255},
	color.RGBA{206, 112, 88, 255},
	color.RGBA{215, 127, 180, 255},
}

type plotParameters struct {
	x_dense   *mat.Dense         // used in heatmap and contour
	y_dense   *mat.Dense         // used in heatmap and contour
	z_dense   *mat.Dense         // used in heatmap and contour
	x_values  []float64          // used in scatter
	y_values  []float64          // used in scatter
	z_values  []float64          // used in scatter
	x, y      [][]float64        // used in lines plot
	gradient  colorgrad.Gradient // colormap
	n_levels  int                // colormap levels
	title     string             // title for all
	xlabel    string             // xlabel for all
	ylabel    string             // ylabel for all
	legend    []string           // mainly used in lines plots
	xwidth    int                // xwidth of the saved figure
	ywidth    int                // ywidth of the saved figure
	plot_name string             // flag to call plotter
}

// struct that defines methods to match the GridXYZ interface defined in gonum plot library
// used in heatmap and contour plots
type unitGrid struct {
	x, y, Data *mat.Dense
}

// methods to match the GridXYZ interface defined in gonum plot library
func (g unitGrid) Dims() (c, r int)   { r, c = g.Data.Dims(); return c, r }
func (g unitGrid) Z(c, r int) float64 { return g.Data.At(r, c) }
func (g unitGrid) X(c int) float64    { return g.x.At(0, c) }
func (g unitGrid) Y(r int) float64    { return g.y.At(0, r) }

// struct that defines methods to match the Palette interface defined in gonum plot library
// used in heatmap and contour plots
type colorsGradient struct {
	ColorList []color.Color
}

// methods to match the Palette interface defined in gonum plot library
func (g colorsGradient) Colors() []color.Color { return g.ColorList }

// generate linearly spaced slice of float64
func Linspace(start, stop float64, num int) []float64 {
	var step float64
	if num == 1 {
		return []float64{start}
	}
	step = (stop - start) / float64(num-1)

	r := make([]float64, num)
	for i := 0; i < num; i++ {
		r[i] = start + float64(i)*step
	}
	return r
}
