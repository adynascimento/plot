package plotter

import (
	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

type contourOptions struct {
	nLevels      int
	gradient     colorgrad.Gradient
	lineSettings lineSettings
	colorBar     colorBar
}

type lineSettings struct {
	show  bool
	width font.Length
	style []vg.Length
}

type colorBar struct {
	show     bool
	gradient colorgrad.Gradient
	min, max float64
	position positionType
}

func WithLevels(levels int) func(*contourOptions) {
	return func(co *contourOptions) {
		co.nLevels = levels
	}
}

func WithGradient(gradient colorgrad.Gradient) func(*contourOptions) {
	return func(co *contourOptions) {
		co.gradient = gradient
		co.colorBar.gradient = gradient
	}
}

func WithContourLines() func(*contourOptions) {
	return func(co *contourOptions) {
		co.lineSettings.show = true
	}
}

func WithContourLineWidth(width float64) func(*contourOptions) {
	return func(co *contourOptions) {
		co.lineSettings.width = vg.Points(width)
	}
}

func WithContourLineStyle(style lineStyleType) func(*contourOptions) {
	return func(co *contourOptions) {
		co.lineSettings.style = style
	}
}

func WithColorbar(position positionType) func(*contourOptions) {
	return func(co *contourOptions) {
		if co.gradient == (colorgrad.Gradient{}) {
			co.colorBar.show = false
		} else {
			co.colorBar.show = true
			co.colorBar.position = position
		}
	}
}
