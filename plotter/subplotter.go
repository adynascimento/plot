package plotter

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/vgimg"
)

type Subplot interface {
	Subplot(row, col int) PlotterInterface
	FigSize(xwidth, ywidth int)
	Save(name string)
	Show()
}

func NewSubplot(rows, cols int) Subplot {
	subplots := make([][]*plotParameters, rows)
	for j := range subplots {
		subplots[j] = make([]*plotParameters, cols)
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

// initialize each subplot individually: row and column indices are 1-based
func (plt *subplotParameters) Subplot(row, col int) PlotterInterface {
	if row < 1 || row > plt.rows {
		panic(fmt.Sprintf("subplot row out of range: %d", row))
	}
	if col < 1 || col > plt.cols {
		panic(fmt.Sprintf("subplot column out of range: %d", col))
	}

	// convert to 0-based slice indices
	row--
	col--

	pParameters := &plotParameters{
		plot: plot.New(),
		line: lineOptions{
			usedColors: make(map[color.Color]bool),
		},
	}
	plt.subplots[row][col] = pParameters

	return pParameters
}

// draw plot to a figure
func (plt *subplotParameters) DrawPlot() {
	xwidth := plt.figSize.xwidth
	ywidth := plt.figSize.ywidth

	// padding between adjacent subplots
	padX := vg.Centimeter
	padY := vg.Centimeter

	// calculate the available size for each subplot.
	sXWidth := (xwidth - vg.Length(plt.cols-1)*padX) / vg.Length(plt.cols)
	sYWidth := (ywidth - vg.Length(plt.rows-1)*padY) / vg.Length(plt.rows)

	// spacing between plot and colorbar
	// calculate the extra space required by vertical and horizontal
	spacing := font.Length(1.5) * vg.Centimeter
	spaceVCB := make([]float64, plt.cols)
	spaceHCB := make([]float64, plt.rows)

	// draw each subplot independently
	for j := 0; j < plt.rows; j++ {
		for i := 0; i < plt.cols; i++ {
			p := plt.subplots[j][i]

			if p.colorbar.show {
				switch p.colorbar.position {
				case Vertical:
					spaceVCB[i] = float64(1.5 * spacing)

				case Horizontal:
					spaceHCB[j] = float64(spacing)
				}
			}

			plt.subplots[j][i].setFigSize(sXWidth, sYWidth)
			plt.subplots[j][i].DrawPlot()
		}
	}

	vColorbar := floats.Sum(spaceVCB)
	hColorbar := floats.Sum(spaceHCB)

	// new image canvas including the space required by the colorbars
	img := vgimg.New(1.02*xwidth+font.Length(vColorbar), ywidth+font.Length(hColorbar))

	// compose the individually rendered subplots into the final figure
	plt.Align(padX, padY, img)
	plt.figure = vgimg.PngCanvas{Canvas: img}
}

// compose subplot images into the final figure
func (plt *subplotParameters) Align(padX, padY vg.Length, img *vgimg.Canvas) {
	dst := img.Image()

	totalWidth := dst.Bounds().Dx()
	totalHeight := dst.Bounds().Dy()

	xwidth, ywidth := img.Size()

	// convert the subplot padding from Gonum units to pixels
	padXpx := int(float64(padX) / float64(xwidth) * float64(totalWidth))
	padYpx := int(float64(padY) / float64(ywidth) * float64(totalHeight))

	// determine the required width of each column and height of each row
	// colorbars are already included in the rendered subplot images
	colWidths := make([]int, plt.cols)
	rowHeights := make([]int, plt.rows)
	for j := 0; j < plt.rows; j++ {
		for i := 0; i < plt.cols; i++ {
			p := plt.subplots[j][i]
			if p == nil || p.figure.Image() == nil {
				continue
			}
			bounds := p.figure.Image().Bounds()
			width := bounds.Dx()
			height := bounds.Dy()
			if width > colWidths[i] {
				colWidths[i] = width
			}
			if height > rowHeights[j] {
				rowHeights[j] = height
			}
		}
	}

	// calculate the horizontal position of each column
	colX := make([]int, plt.cols)
	for i := 1; i < plt.cols; i++ {
		colX[i] = colX[i-1] + colWidths[i-1] + padXpx
	}
	// calculate the vertical position of each row
	rowY := make([]int, plt.rows)
	for j := 1; j < plt.rows; j++ {
		rowY[j] = rowY[j-1] + rowHeights[j-1] + padYpx
	}

	// copy each rendered subplot image into its corresponding position
	for j := 0; j < plt.rows; j++ {
		for i := 0; i < plt.cols; i++ {
			p := plt.subplots[j][i]
			if p == nil {
				continue
			}
			src := p.figure.Image()
			if src == nil {
				continue
			}
			x := colX[i]
			y := rowY[j]
			draw.Draw(
				dst,
				image.Rect(
					x,
					y,
					x+src.Bounds().Dx(),
					y+src.Bounds().Dy(),
				),
				src,
				src.Bounds().Min,
				draw.Over,
			)
		}
	}
}

// show plot in graphical window
func (plt *subplotParameters) Show() {
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
	plt.figSize.xwidth = font.Length(xwidth) * vg.Centimeter
	plt.figSize.ywidth = font.Length(ywidth) * vg.Centimeter
}
