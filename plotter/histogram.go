package plotter

import (
	"image/color"
	"math"

	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/gonum/stat/distuv"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
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

// parameters to histogram plot
func (plt *plotParameters) Hist(x []float64, n int, options ...func(*histogramOptions)) {
	// default options
	plt.histogramOptions = histogramOptions{
		fillColor: Gray,
		lineStyle: Solid,
		lineWidth: vg.Points(1.5),
	}

	// apply additional options
	for _, option := range options {
		option(&plt.histogramOptions)
	}

	// prepare data to plot
	vs := make(plotter.Values, len(x))
	copy(vs, x)

	// make a histogram plotter
	h, _ := plotter.NewHist(vs, n)
	h.FillColor = plt.histogramOptions.fillColor
	h.LineStyle.Color = plt.histogramOptions.lineColor
	h.LineStyle.Width = plt.histogramOptions.lineWidth
	h.LineStyle.Dashes = plt.histogramOptions.lineStyle
	h.Normalize(1)

	// add the plotters to the plot
	plt.plot.Add(h)

	var std float64
	var xline []float64
	if plt.histogramOptions.kdeCurve || plt.histogramOptions.normalCurve {
		xmin, xmax, _, _ := h.DataRange()
		xline = Linspace(xmin, xmax, 1000)

		// standard deviation
		std = math.Sqrt(stat.Variance(x, nil))
	}

	// calculate density curve (kde) mathematically
	if plt.histogramOptions.kdeCurve {
		n := float64(len(x))
		yline := make([]float64, len(xline))

		// scott's rule
		kernel := distuv.Normal{
			Sigma: std * math.Pow(n, -0.2),
		}
		for i, xi := range xline {
			sum := 0.0
			for _, sample := range x {
				sum += kernel.Prob(xi - sample)
			}
			yline[i] = sum / n
		}

		// add the curve to the plot
		plt.Plot(xline, yline,
			WithLineColor(plt.histogramOptions.kdeCurveColor),
			WithLineWidth(2.0),
		)
	}

	// calculate normal theoretical density curve mathematically
	if plt.histogramOptions.normalCurve {
		yline := make([]float64, len(xline))

		distNormal := distuv.Normal{
			Mu:    stat.Mean(x, nil),
			Sigma: std,
		}
		for i, xi := range xline {
			yline[i] = distNormal.Prob(xi)
		}

		// add the curve to the plot
		plt.Plot(xline, yline,
			WithLineColor(plt.histogramOptions.normalCurveColor),
			WithLineWidth(2.0),
		)
	}
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
