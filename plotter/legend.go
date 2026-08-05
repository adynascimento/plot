package plotter

import "gonum.org/v1/plot/vg"

type location string

const (
	UpperRight location = "upper_right"
	UpperLeft  location = "upper_left"
	LowerRight location = "lower_right"
	LowerLeft  location = "lower_left"
)

type LegendLocation interface {
	Location(loc location)
}

// legend mainly used in line plots
func (plt *plotParameters) Legend(str ...string) LegendLocation {
	// legend style
	if len(plt.legends) > 0 {
		str = str[:min(len(str), len(plt.legends))]
		for i, legend := range str {
			plt.plot.Legend.Add(legend, plt.legends[i]...)
		}
		plt.plot.Legend.XOffs = -5. * vg.Millimeter
		plt.plot.Legend.Padding = vg.Millimeter
	}

	return plt
}

// legend location
func (plt *plotParameters) Location(loc location) {
	switch loc {
	case UpperRight:
		plt.plot.Legend.Top = true
		plt.plot.Legend.Left = false
		plt.plot.Legend.XOffs = -5. * vg.Millimeter
		plt.plot.Legend.Padding = vg.Millimeter

	case UpperLeft:
		plt.plot.Legend.Top = true
		plt.plot.Legend.Left = true
		plt.plot.Legend.XOffs = 3. * vg.Millimeter
		plt.plot.Legend.Padding = vg.Millimeter

	case LowerRight:
		plt.plot.Legend.Top = false
		plt.plot.Legend.Left = false
		plt.plot.Legend.XOffs = -5. * vg.Millimeter
		plt.plot.Legend.Padding = vg.Millimeter

	case LowerLeft:
		plt.plot.Legend.Top = false
		plt.plot.Legend.Left = true
		plt.plot.Legend.XOffs = 3. * vg.Millimeter
		plt.plot.Legend.Padding = vg.Millimeter
	}
}
