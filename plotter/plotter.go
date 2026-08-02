package plotter

import (
	"image"
	"image/color"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type PlotterInterface interface {
	Plot(x, y []float64, options ...func(*lineOptions))
	Contour(x, y, z *mat.Dense, options ...func(*contourOptions))
	ContourF(x, y, z *mat.Dense, options ...func(*contourOptions))
	Scatter(x, y, z []float64, options ...func(*scatterOptions))
	Hist(x []float64, n int, options ...func(*histogramOptions))
	ImShow(x []*mat.Dense)
	Title(str string)
	XLabel(xlabel string)
	YLabel(ylabel string)
	Legend(str ...string) LegendLocation
	XLim(xmin, xmax float64)
	YLim(ymin, ymax float64)
	Grid()
}

type Plot interface {
	PlotterInterface
	FigSize(xwidth, ywidth int)
	Save(name string)
	Show()
}

func NewPlot() Plot {
	return &plotParameters{
		plot: plot.New(),
		lineOptions: lineOptions{
			usedColors: make(map[color.Color]bool),
		},
		figSize: figSize{
			xwidth: 10,
			ywidth: 10,
		},
		figure: vgimg.PngCanvas{
			Canvas: &vgimg.Canvas{},
		},
	}
}

// parameters to image plot
func (plt *plotParameters) ImShow(x []*mat.Dense) {
	// prepare data to plot
	var img image.Image
	rows, cols := x[0].Dims()

	xmin := 0.0
	ymin := 0.0
	xmax := float64(cols)
	ymax := float64(rows)

	if len(x) == 1 {
		// grayscale image
		grayImg := image.NewGray(image.Rect(int(xmin), int(ymin), int(xmax), int(ymax)))
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				v := uint8(x[0].At(i, j))
				grayImg.SetGray(j, i, color.Gray{Y: v})
			}
		}
		img = grayImg
	} else {
		// RGB image
		rgbImg := image.NewRGBA(image.Rect(int(xmin), int(ymin), int(xmax), int(ymax)))
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				r := uint8(x[0].At(i, j))
				g := uint8(x[1].At(i, j))
				b := uint8(x[2].At(i, j))
				rgbImg.SetRGBA(j, i, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}
		img = rgbImg
	}

	// add and make a image plotter
	plt.plot.Add(plotter.NewImage(img, xmin, ymin, xmax, ymax))
	plt.plot.X.Max = float64(rows) + 0.02*float64(rows)
}

// draw plot to a figure
func (plt *plotParameters) DrawPlot() {
	xwidth := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywidth := font.Length(plt.figSize.ywidth) * vg.Centimeter

	// new image canvas
	img := vgimg.New(xwidth, ywidth)

	// draw the plot
	plt.plot.Draw(draw.Canvas{
		Canvas: draw.New(img),
		Rectangle: vg.Rectangle{
			Min: vg.Point{X: 0, Y: 0},
			Max: vg.Point{X: xwidth, Y: ywidth},
		},
	})

	// add colorbar to plot
	if plt.colorBar.show {
		switch plt.colorBar.position {
		case Vertical:
			img = plt.drawVerticalColorBar(xwidth, ywidth)
		case Horizontal:
			img = plt.drawHorizontalColorBar(xwidth, ywidth)
		}
	}

	plt.figure = vgimg.PngCanvas{Canvas: img}
}

// show plot in graphical window
func (plt *plotParameters) Show() {
	if plt.figure.Image() == nil {
		plt.DrawPlot()
	}

	// graphical window creation
	imgData := plt.figure.Image()
	go func() {
		window := new(app.Window)
		window.Option(
			app.Title("Plot Viewer"),
			app.Size(
				unit.Dp(float32(imgData.Bounds().Dx())),
				unit.Dp(float32(imgData.Bounds().Dy())),
			),
		)

		img := widget.Image{
			Src:      paint.NewImageOp(imgData),
			Fit:      widget.Contain,
			Position: layout.Center,
			Scale:    1,
		}

		var ops op.Ops
		for {
			switch e := window.Event().(type) {
			case app.DestroyEvent:
				os.Exit(0)
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				img.Layout(gtx)
				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}

// save the plot to an image file
func (plt *plotParameters) Save(file string) {
	// save the plot to a PNG file.
	if plt.figure.Image() == nil {
		plt.DrawPlot()
	}

	// save the image to a file
	w, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if _, err := plt.figure.WriteTo(w); err != nil {
		panic(err)
	}
}

// size of the saved figure
func (plt *plotParameters) FigSize(xwidth, ywidth int) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
}

// title for all plots
func (plt *plotParameters) Title(title string) {
	plt.plot.Title.Text = title
}

// xlabel for all plots
func (plt *plotParameters) XLabel(xlabel string) {
	plt.plot.X.Label.Text = xlabel
}

// ylabel for all plots
func (plt *plotParameters) YLabel(ylabel string) {
	plt.plot.Y.Label.Text = ylabel
}

// set the x-axis vies limits
func (plt *plotParameters) XLim(xmin, xmax float64) {
	if xmin < xmax {
		plt.plot.X.Min = xmin
		plt.plot.X.Max = xmax
	}
}

// set the x-axis vies limits
func (plt *plotParameters) YLim(ymin, ymax float64) {
	if ymin < ymax {
		plt.plot.Y.Min = ymin
		plt.plot.Y.Max = ymax
	}
}

// draw grid with both vertical and horizontal lines
func (plt *plotParameters) Grid() {
	plt.plot.Add(plotter.NewGrid())
}
