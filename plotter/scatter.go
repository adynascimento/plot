package plotter

import (
	"image/color"

	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type scatterOptions struct {
	color      color.Color
	gradient   colorgrad.Gradient
	marker     draw.GlyphDrawer
	markerSize font.Length
	colorBar   colorBar
}

func WithMarkerColor(color colorType) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.color = color
	}
}

func WithScatterGradient(gradient colorgrad.Gradient) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.gradient = gradient
		so.colorBar.gradient = gradient
	}
}

func WithScatterMarker(marker markerType) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.marker = marker
	}
}

func WithScatterMarkerSize(size float64) func(*scatterOptions) {
	return func(so *scatterOptions) {
		so.markerSize = vg.Points(size)
	}
}

func WithScatterColorbar(position positionType) func(*scatterOptions) {
	return func(so *scatterOptions) {
		if so.gradient == (colorgrad.Gradient{}) {
			so.colorBar.show = false
		} else {
			so.colorBar.show = true
			so.colorBar.position = position
		}
	}
}
