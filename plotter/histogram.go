package plotter

import (
	"image/color"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

type histogramOptions struct {
	fillColor    color.Color
	lineColor    color.Color
	lineStyle    []vg.Length
	lineWidth    font.Length
	densityCurve bool
	densityColor color.Color
}

func WithHistFillColor(color colorType) func(*histogramOptions) {
	return func(ho *histogramOptions) {
		ho.fillColor = color
	}
}

func WithHistLineColor(color colorType) func(*histogramOptions) {
	return func(ho *histogramOptions) {
		ho.lineColor = color
	}
}

func WithHistLineWidth(width float64) func(*histogramOptions) {
	return func(ho *histogramOptions) {
		ho.lineWidth = vg.Points(width)
	}
}

func WithHistLineStyle(style lineStyleType) func(*histogramOptions) {
	return func(ho *histogramOptions) {
		ho.lineStyle = style
	}
}

func WithHistNormalCurve(color colorType) func(*histogramOptions) {
	return func(ho *histogramOptions) {
		ho.densityCurve = true
		ho.densityColor = color
	}
}
