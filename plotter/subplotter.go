package plotter

import (
	"image/color"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

type Subplot interface {
	Subplot(row, col int) PlotterInterface
	FigSize(xwidth, ywidth int)
	Save(name string)
	Show()
}

func NewSubplot(rows, cols int) Subplot {
	subplots := make([][]*plot.Plot, rows)
	for j := range subplots {
		subplots[j] = make([]*plot.Plot, cols)
	}

	return &subplotParameters{
		rows:     rows,
		cols:     cols,
		subplots: subplots,
		figSize: figSize{
			xwidth: 15,
			ywidth: 10,
		},
		figure: vgimg.PngCanvas{
			Canvas: &vgimg.Canvas{},
		},
	}
}

// initialize each subplot individually
func (plt *subplotParameters) Subplot(row, col int) PlotterInterface {
	p := plot.New()
	plt.subplots[row][col] = p
	return &plotParameters{
		plot: p,
		lineOptions: lineOptions{
			usedColors: make(map[color.Color]bool),
		},
	}
}

// draw plot to a figure
func (plt *subplotParameters) DrawPlot() {
	xwidth := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywidth := font.Length(plt.figSize.ywidth) * vg.Centimeter

	// new image canvas
	img := vgimg.New(xwidth, ywidth)

	canvases := plot.Align(plt.subplots, draw.Tiles{
		Rows: plt.rows,
		Cols: plt.cols,
		PadX: vg.Centimeter,
		PadY: vg.Centimeter,
	}, draw.New(img))
	for j := 0; j < plt.rows; j++ {
		for i := 0; i < plt.cols; i++ {
			if plt.subplots[j][i] != nil {
				plt.subplots[j][i].Draw(canvases[j][i])
			}
		}
	}

	plt.figure = vgimg.PngCanvas{Canvas: img}
}

// show plot in graphical window
func (plt *subplotParameters) Show() {
	if plt.figure.Image() == nil {
		plt.DrawPlot()
	}

	// graphical window creation
	imgData := plt.figure.Image()
	window := app.NewWindow(
		app.Title("Plot Viewer"),
		app.Size(unit.Dp(float32(imgData.Bounds().Dx())),
			unit.Dp(float32(imgData.Bounds().Dy()))),
	)

	var ops op.Ops
	for e := range window.Events() {
		switch e := e.(type) {
		case system.DestroyEvent:
			return
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				imgOp := paint.NewImageOp(imgData)
				imgOp.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
			e.Frame(gtx.Ops)
		}
	}
	app.Main()
}

// save the plot to an image file
func (plt *subplotParameters) Save(file string) {
	// save the plot to a PNG file.
	if plt.figure.Image() == nil {
		plt.DrawPlot()
	}

	// save the image to a file
	w, err := os.Create(file)
	if err != nil {
		log.Panic(err)
	}
	defer w.Close()

	if _, err := plt.figure.WriteTo(w); err != nil {
		panic(err)
	}
}

// size of the saved figure
func (plt *subplotParameters) FigSize(xwidth, ywidth int) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
}
