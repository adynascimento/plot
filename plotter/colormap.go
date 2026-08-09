package plotter

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
)

var (
	Viridis          = colormap{gradient: colorgrad.Viridis()}
	Turbo            = colormap{gradient: colorgrad.Turbo()}
	Plasma           = colormap{gradient: colorgrad.Plasma()}
	Magma            = colormap{gradient: colorgrad.Magma()}
	Inferno          = colormap{gradient: colorgrad.Inferno()}
	Cividis          = colormap{gradient: colorgrad.Cividis()}
	Rainbow          = colormap{gradient: colorgrad.Rainbow()}
	Cool             = colormap{gradient: colorgrad.Cool()}
	BrBG             = colormap{gradient: colorgrad.BrBG()}
	PRGn             = colormap{gradient: colorgrad.PRGn()}
	PiYG             = colormap{gradient: colorgrad.PiYG()}
	PuOr             = colormap{gradient: colorgrad.PuOr()}
	RdBu             = colormap{gradient: colorgrad.RdBu()}
	RdGy             = colormap{gradient: colorgrad.RdGy()}
	RdYlBu           = colormap{gradient: colorgrad.RdYlBu()}
	RdYlGn           = colormap{gradient: colorgrad.RdYlGn()}
	Spectral         = colormap{gradient: colorgrad.Spectral()}
	Blues            = colormap{gradient: colorgrad.Blues()}
	Greens           = colormap{gradient: colorgrad.Greens()}
	Greys            = colormap{gradient: colorgrad.Greys()}
	Oranges          = colormap{gradient: colorgrad.Oranges()}
	Purples          = colormap{gradient: colorgrad.Purples()}
	Reds             = colormap{gradient: colorgrad.Reds()}
	BuGn             = colormap{gradient: colorgrad.BuGn()}
	BuPu             = colormap{gradient: colorgrad.BuPu()}
	GnBu             = colormap{gradient: colorgrad.GnBu()}
	OrRd             = colormap{gradient: colorgrad.OrRd()}
	PuBuGn           = colormap{gradient: colorgrad.PuBuGn()}
	PuBu             = colormap{gradient: colorgrad.PuBu()}
	PuRd             = colormap{gradient: colorgrad.PuRd()}
	RdPu             = colormap{gradient: colorgrad.RdPu()}
	YlGnBu           = colormap{gradient: colorgrad.YlGnBu()}
	YlGn             = colormap{gradient: colorgrad.YlGn()}
	YlOrBr           = colormap{gradient: colorgrad.YlOrBr()}
	YlOrRd           = colormap{gradient: colorgrad.YlOrRd()}
	Sinebow          = colormap{gradient: colorgrad.Sinebow()}
	Warm             = colormap{gradient: colorgrad.Warm()}
	CubehelixDefault = colormap{gradient: colorgrad.CubehelixDefault()}
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
