package plotter

import (
	"image/color"

	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
)

type histogramOptions struct {
	fillColor        color.Color
	lineColor        color.Color
	lineStyle        []vg.Length
	lineWidth        font.Length
	kdeCurve         bool
	kdeCurveColor    color.Color
	normalCurve      bool
	normalCurveColor color.Color
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

func WithHistKDECurve(color colorType) func(*histogramOptions) {
	return func(ho *histogramOptions) {
		ho.kdeCurve = true
		ho.kdeCurveColor = color
	}
}

func WithHistNormalCurve(color colorType) func(*histogramOptions) {
	return func(ho *histogramOptions) {
		ho.normalCurve = true
		ho.normalCurveColor = color
	}
}
