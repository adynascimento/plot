package plotter

import (
	"image/color"

	"gioui.org/app"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/vg/vgimg"
)

type plotParameters struct {
	plot      *plot.Plot           // initialize new plot
	line      lineOptions          // line plotter options
	contour   contourOptions       // contour plotter options
	scatter   scatterOptions       // scatter plotter options
	histogram histogramOptions     // histogram plotter options
	legends   [][]plot.Thumbnailer // legend plotter config
	figSize   figSize              // xwidth and ywidth of the saved figure
	figure    vgimg.PngCanvas      // figure to plot and save
	window    *app.Window          // window to show the figure
	animation *animationOptions    // animation options
	colorbar  colorbar             // show colorbar with gradient
}

type subplotParameters struct {
	rows     int
	cols     int
	subplots [][]*plotParameters //  plots for subplot
	figSize  figSize             // xwidth and ywidth of the saved figure
	figure   vgimg.PngCanvas     // figure to plot and save
}

type figSize struct {
	xwidth, ywidth font.Length
}

// set figure size using Gonum units
func (plt *plotParameters) setFigSize(xwidth, ywidth font.Length) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
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
type colormapPalette struct {
	colorList []color.Color
	colormap  Colormap
	min, max  float64
}

// methods to match the Palette interface defined in gonum plot library
func (p *colormapPalette) Colors() []color.Color {
	if len(p.colorList) == 0 {
		p.colorList = p.colormap.Colors(256)
	}
	return p.colorList
}
func (p *colormapPalette) Alpha() float64 { return 1.0 }
func (p *colormapPalette) At(v float64) (color.Color, error) {
	if p.max == p.min {
		return p.colormap.At(0.5), nil
	}
	return p.colormap.At((v - p.min) / (p.max - p.min)), nil
}
func (p *colormapPalette) Max() float64 { return p.max }
func (p *colormapPalette) Min() float64 { return p.min }
func (p *colormapPalette) Palette(n int) palette.Palette {
	return &colormapPalette{
		colorList: p.colormap.Colors(n),
		colormap:  p.colormap,
		min:       p.min,
		max:       p.max,
	}
}
func (p *colormapPalette) SetMax(v float64)   { p.max = v }
func (p *colormapPalette) SetMin(v float64)   { p.min = v }
func (p *colormapPalette) SetAlpha(a float64) {}

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
