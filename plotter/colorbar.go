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
	p := plot.New()
	p.HideX()
	p.HideY()
	p.Y.Padding = 0
	c := &plotter.ColorBar{
		ColorMap: &colormapPalette{
			colormap: plt.colorbar.colormap,
		},
	}
	c.ColorMap.SetMin(plt.colorbar.min)
	c.ColorMap.SetMax(plt.colorbar.max)
	c.Vertical = true
	p.Add(c)

	// spacing between plot and colorbar
	spacing := font.Length(1.5) * vg.Centimeter

	img := vgimg.New(xwidth+1.5*spacing, ywidth)
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
	p.Draw(draw.Canvas{
		Canvas: dCanvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: (xwidth + 0.3*spacing), Y: 1.25 * vg.Centimeter},
			Max: vg.Point{X: (xwidth + 0.7*spacing), Y: ywidth - 0.65*vg.Centimeter},
		},
	})

	// create a separate plot for the right y-axis
	axisPlot := plot.New()
	axisPlot.HideX()
	axisPlot.Y.Min = plt.colorbar.min
	axisPlot.Y.Max = plt.colorbar.max
	axisPlot.Y.Padding = 0

	// draw the y-axis on the right side of the colorbar
	drawRightYAxis(
		draw.Canvas{
			Canvas: dCanvas,
			Rectangle: vg.Rectangle{
				Min: vg.Point{
					X: xwidth + 0.7*spacing,
					Y: 1.44 * vg.Centimeter,
				},
				Max: vg.Point{
					X: xwidth + 1.0*spacing,
					Y: ywidth - 0.67*vg.Centimeter,
				},
			},
		},
		axisPlot.Y,
	)

	return img
}

func (plt *plotParameters) drawHorizontalColorbar(xwidth, ywidth font.Length) *vgimg.Canvas {
	// create a new plot for horizontal colorbar
	p := plot.New()
	p.HideY()
	p.X.Padding = 0
	c := &plotter.ColorBar{
		ColorMap: &colormapPalette{
			colormap: plt.colorbar.colormap,
		},
	}
	c.ColorMap.SetMin(plt.colorbar.min)
	c.ColorMap.SetMax(plt.colorbar.max)
	p.Add(c)

	// spacing between plot and colorbar
	spacing := font.Length(1.5) * vg.Centimeter

	img := vgimg.New(1.05*xwidth, ywidth+spacing)
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
	p.Draw(draw.Canvas{
		Canvas: dCanvas,
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 1.22 * vg.Centimeter, Y: 0},
			Max: vg.Point{X: xwidth, Y: 1.2 * vg.Centimeter},
		},
	})

	return img
}

func drawRightYAxis(c draw.Canvas, a plot.Axis) {
	marks := a.Tick.Marker.Ticks(a.Min, a.Max)

	// draw axis line
	x := c.Min.X
	c.StrokeLine2(a.LineStyle, x, c.Min.Y, x, c.Max.Y)

	// draw ticks
	major := false
	length := a.Tick.Length
	if a.Tick.Width > 0 && a.Tick.Length > 0 {
		for _, t := range marks {
			y := c.Y(a.Norm(t.Value))
			if !c.ContainsY(y) {
				continue
			}
			end := length
			if t.IsMinor() {
				end = length / 2
			}
			c.StrokeLine2(a.Tick.LineStyle, x, y, x+end, y)
			if !t.IsMinor() {
				major = true
			}
		}
	}

	// leave the same spacing used by Gonum between ticks and labels
	labelX := x + length
	if major {
		labelX += a.Tick.Label.Width(" ")
	}

	// draw labels to the right of the ticks
	labelStyle := a.Tick.Label
	labelStyle.XAlign = draw.XLeft
	for _, t := range marks {
		if t.IsMinor() {
			continue
		}
		y := c.Y(a.Norm(t.Value))
		if !c.ContainsY(y) {
			continue
		}
		c.FillText(
			labelStyle,
			vg.Point{
				X: labelX,
				Y: y + a.Tick.Label.FontExtents().Descent,
			},
			t.Label,
		)
	}
}
