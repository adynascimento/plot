package plotter

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type positionType string

var (
	Vertical   positionType = "vertical"
	Horizontal positionType = "horizontal"
)

type colorbar struct {
	show     bool
	colormap Colormap
	min, max float64
	position positionType
}

func (plt *plotParameters) drawVerticalColorbar(xwidth, ywidth font.Length) *vgimg.Canvas {
	// create a new plot for vertical colorbar
	c := plot.New()
	c.HideX()
	c.Y.Padding = 0
	l := &plotter.ColorBar{
		ColorMap: &colormapPalette{
			colormap: plt.colorbar.colormap,
		},
	}
	l.ColorMap.SetMin(plt.colorbar.min)
	l.ColorMap.SetMax(plt.colorbar.max)
	l.Vertical = true
	c.Add(l)

	// spacing between plot and colorbar
	spacing := font.Length(1.5) * vg.Centimeter

	img := vgimg.New(xwidth+1.35*spacing, ywidth)
	dCanvas := draw.New(img)

	// draw the principal plot in the first column
	plt.plot.Draw(draw.Canvas{
		Canvas: dCanvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: xwidth, Y: ywidth},
		},
	})

	// draw the colorbar in the second column
	c.Draw(draw.Canvas{
		Canvas: dCanvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: (xwidth + 0.3*spacing), Y: vg.Centimeter},
			Max: vg.Point{X: (xwidth + 1.3*spacing), Y: 0.93 * ywidth},
		},
	})

	return img
}

func (plt *plotParameters) drawHorizontalColorbar(xwidth, ywidth font.Length) *vgimg.Canvas {
	// create a new plot for horizontal colorbar
	c := plot.New()
	c.HideY()
	c.X.Padding = 0
	l := &plotter.ColorBar{
		ColorMap: &colormapPalette{
			colormap: plt.colorbar.colormap,
		},
	}
	l.ColorMap.SetMin(plt.colorbar.min)
	l.ColorMap.SetMax(plt.colorbar.max)
	c.Add(l)

	// spacing between plot and colorbar
	spacing := font.Length(1.5) * vg.Centimeter

	img := vgimg.New(xwidth, ywidth+spacing)
	dCanvas := draw.New(img)

	// draw the principal plot in the first row
	plt.plot.Draw(draw.Canvas{
		Canvas: dCanvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: spacing},
			Max: vg.Point{X: xwidth, Y: ywidth + spacing},
		},
	})

	// Draw the color bar in the second column
	c.Draw(draw.Canvas{
		Canvas: dCanvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: vg.Centimeter, Y: 0},
			Max: vg.Point{X: 0.99 * xwidth, Y: 1.2 * vg.Centimeter},
		},
	})

	return img
}
