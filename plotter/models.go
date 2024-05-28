package plotter

import (
	"image/color"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
)

type PlotParameters struct {
	plot           *plot.Plot           // initialize new plot
	lineOptions    lineOptions          // line plotter options
	contourOptions contourOptions       // contour plotter options
	legends        [][]plot.Thumbnailer // legend plotter config
	figSize        figSize              // xwidth and ywidth of the saved figure
}

type SubplotParameters struct {
	rows     int
	cols     int
	subplots [][]*plot.Plot // plots for subplot
	figSize  figSize        // xwidth and ywidth of the saved figure
}

type figSize struct{ xwidth, ywidth int }

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

// applies the function fn to each of the elements of a. The function fn takes a row/column
// index and element value and returns some function of that tuple
func Apply(fn func(i, j int, v float64) float64, a mat.Matrix) *mat.Dense {
	m := new(mat.Dense)
	m.Apply(fn, a)

	return m
}
