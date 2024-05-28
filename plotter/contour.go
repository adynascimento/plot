package plotter

import (
	"github.com/mazznoer/colorgrad"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

type contourOptions struct {
	nLevels   int
	gradient  colorgrad.Gradient
	lineWidth font.Length
	lineStyle []vg.Length
	addLines  bool
}

func WithLevels(levels int) func(*contourOptions) {
	return func(co *contourOptions) {
		co.nLevels = levels
	}
}

func WithGradient(gradient colorgrad.Gradient) func(*contourOptions) {
	return func(co *contourOptions) {
		co.gradient = gradient
	}
}

func WithContourLineWidth(width float64) func(*contourOptions) {
	return func(co *contourOptions) {
		co.lineWidth = vg.Points(width)
	}
}

func WithContourLineStyle(style string) func(*contourOptions) {
	return func(co *contourOptions) {
		if l, exists := linestyle[style]; exists {
			co.lineStyle = l
		}
	}
}

func WithContourLines() func(*contourOptions) {
	return func(co *contourOptions) {
		co.addLines = true
	}
}
