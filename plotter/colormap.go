package plotter

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
)

var (
	// sequential
	Viridis = colormap{gradient: colorgrad.Viridis()}
	Plasma  = colormap{gradient: colorgrad.Plasma()}
	Magma   = colormap{gradient: colorgrad.Magma()}
	Cool    = colormap{gradient: colorgrad.Cool()}
	Inferno = colormap{gradient: colorgrad.Inferno()}
	Cividis = colormap{gradient: colorgrad.Cividis()}
	Greys   = colormap{gradient: colorgrad.Greys()}
	Purples = colormap{gradient: colorgrad.Purples()}
	Blues   = colormap{gradient: colorgrad.Blues()}
	Greens  = colormap{gradient: colorgrad.Greens()}
	Oranges = colormap{gradient: colorgrad.Oranges()}
	Reds    = colormap{gradient: colorgrad.Reds()}
	YlOrBr  = colormap{gradient: colorgrad.YlOrBr()}
	YlOrRd  = colormap{gradient: colorgrad.YlOrRd()}
	YlGn    = colormap{gradient: colorgrad.YlGn()}
	YlGnBu  = colormap{gradient: colorgrad.YlGnBu()}
	GnBu    = colormap{gradient: colorgrad.GnBu()}
	BuGn    = colormap{gradient: colorgrad.BuGn()}
	BuPu    = colormap{gradient: colorgrad.BuPu()}
	PuBu    = colormap{gradient: colorgrad.PuBu()}
	PuBuGn  = colormap{gradient: colorgrad.PuBuGn()}
	PuRd    = colormap{gradient: colorgrad.PuRd()}
	RdPu    = colormap{gradient: colorgrad.RdPu()}
	OrRd    = colormap{gradient: colorgrad.OrRd()}

	// additional sequential
	Hot    = newColormap("#000000", "#ff0000", "#ffff00", "#ffffff")
	Spring = newColormap("#ff00ff", "#ffff00")
	Summer = newColormap("#008066", "#ffff66")
	Autumn = newColormap("#ff0000", "#ffff00")
	Winter = newColormap("#0000ff", "#00ffff")
	Copper = newColormap("#000000", "#ffcc99")

	// diverging
	BrBG     = colormap{gradient: colorgrad.BrBG()}
	PRGn     = colormap{gradient: colorgrad.PRGn()}
	PiYG     = colormap{gradient: colorgrad.PiYG()}
	PuOr     = colormap{gradient: colorgrad.PuOr()}
	RdBu     = colormap{gradient: colorgrad.RdBu()}
	RdGy     = colormap{gradient: colorgrad.RdGy()}
	RdYlBu   = colormap{gradient: colorgrad.RdYlBu()}
	RdYlGn   = colormap{gradient: colorgrad.RdYlGn()}
	Spectral = colormap{gradient: colorgrad.Spectral()}

	// additional diverging
	Coolwarm = newColormap("#3b4cc0", "#688aef", "#9abbff", "#dddcdc", "#f4987a", "#d32f27", "#b40426")
	Seismic  = newColormap("#00004c", "#0000ff", "#ffffff", "#ff0000", "#4c0000")
	Bwr      = newColormap("#0000ff", "#ffffff", "#ff0000")

	// qualitative
	Tab10   = newDiscreteColormap("#1f77b4", "#ff7f0e", "#2ca02c", "#d62728", "#9467bd", "#8c564b", "#e377c2", "#7f7f7f", "#bcbd22", "#17becf")
	Pastel1 = newDiscreteColormap("#fbb4ae", "#b3cde3", "#ccebc5", "#decbe4", "#fed9a6", "#ffffcc", "#e5d8bd", "#fddaec", "#f2f2f2")
	Pastel2 = newDiscreteColormap("#b3e2cd", "#fdcdac", "#cbd5e8", "#f4cae4", "#e6f5c9", "#fff2ae", "#f1e2cc", "#cccccc")
	Dark2   = newDiscreteColormap("#1b9e77", "#d95f02", "#7570b3", "#e7298a", "#66a61e", "#e6ab02", "#a6761d", "#666666")
	Set1    = newDiscreteColormap("#e41a1c", "#377eb8", "#4daf4a", "#984ea3", "#ff7f00", "#ffff33", "#a65628", "#f781bf", "#999999")
	Set2    = newDiscreteColormap("#66c2a5", "#fc8d62", "#8da0cb", "#e78ac3", "#a6d854", "#ffd92f", "#e5c494", "#b3b3b3")
	Set3    = newDiscreteColormap("#8dd3c7", "#ffffb3", "#bebada", "#fb8072", "#80b1d3", "#fdb462", "#b3de69", "#fccde5", "#d9d9d9", "#bc80bd", "#ccebc5", "#ffed6f")

	// miscellaneous
	Turbo   = colormap{gradient: colorgrad.Turbo()}
	Rainbow = colormap{gradient: colorgrad.Rainbow()}
	Warm    = colormap{gradient: colorgrad.Warm()}

	// additional miscellaneous
	Jet = newColormap("#00007f", "#0000ff", "#007fff", "#00ffff", "#7fff7f", "#ffff00", "#ff7f00", "#ff0000", "#7f0000")
)

type Colormap interface {
	At(v float64) color.Color
	Colors(n int) []color.Color
}

type colormap struct {
	gradient colorgrad.Gradient
}

func (c colormap) At(v float64) color.Color {
	return c.gradient.At(v)
}

func (c colormap) Colors(n int) []color.Color {
	return c.gradient.Colors(uint(n))
}

func newColormap(colors ...string) colormap {
	grad, _ := colorgrad.NewGradient().
		HtmlColors(colors...).
		Build()
	return colormap{gradient: grad}
}

func newDiscreteColormap(colors ...string) colormap {
	grad, _ := colorgrad.NewGradient().
		HtmlColors(colors...).
		Build()
	return colormap{gradient: grad.Sharp(uint(len(colors)), 0)}
}
