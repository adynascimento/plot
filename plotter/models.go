package plotter

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
)

type plotParameters struct {
	plot           *plot.Plot           // initialize new plot
	lineOptions    lineOptions          // line plotter options
	contourOptions contourOptions       // contour plotter options
	scatterOptions scatterOptions       // scatter plotter options
	legends        [][]plot.Thumbnailer // legend plotter config
	figSize        figSize              // xwidth and ywidth of the saved figure
	colorBar       colorBar             // show colorbar with gradient
}

type subplotParameters struct {
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
	colorList []color.Color
	gradient  colorgrad.Gradient
	min, max  float64
}

// methods to match the Palette interface defined in gonum plot library
func (g *colorsGradient) Colors() []color.Color { return g.colorList }
func (g *colorsGradient) Alpha() float64        { return 1.0 }
func (g *colorsGradient) At(v float64) (color.Color, error) {
	return g.gradient.At(v), nil
}
func (g *colorsGradient) Max() float64 { return g.max }
func (g *colorsGradient) Min() float64 { return g.min }
func (g *colorsGradient) Palette(n int) palette.Palette {
	return &colorsGradient{
		colorList: g.gradient.Colors(uint(n)),
	}
}
func (g *colorsGradient) SetMax(v float64)   { g.max = v }
func (g *colorsGradient) SetMin(v float64)   { g.min = v }
func (g *colorsGradient) SetAlpha(a float64) {}

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
