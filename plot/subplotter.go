package plot

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Subplot interface {
	Save(name string)
	FigSize(xwidth, ywidth int)
	Subplot(row, col int) PlotterInterface
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
	}
}

// initialize each subplot individually
func (plt *subplotParameters) Subplot(row, col int) PlotterInterface {
	p := plot.New()
	plt.subplots[row][col] = p
	return &plotParameters{
		plot: p,
	}
}

// save the plot to an image file
func (plt *subplotParameters) Save(name string) {
	// save the plot to a PNG file.
	xwdith := font.Length(plt.figSize.xwidth) * vg.Centimeter
	ywdith := font.Length(plt.figSize.ywidth) * vg.Centimeter

	format := strings.TrimPrefix(strings.ToLower(filepath.Ext(name)), ".")
	img, err := draw.NewFormattedCanvas(xwdith, ywdith, format)
	if err != nil {
		log.Panic(err)
	}

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

	w, err := os.Create(name)
	if err != nil {
		log.Panic(err)
	}
	defer w.Close()

	if _, err := img.WriteTo(w); err != nil {
		log.Panic(err)
	}
}

// size of the saved figure
func (plt *subplotParameters) FigSize(xwidth, ywidth int) {
	plt.figSize.xwidth = xwidth
	plt.figSize.ywidth = ywidth
}
