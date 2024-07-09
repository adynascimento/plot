package main

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/adynascimento/plot/plotter"
	"gonum.org/v1/gonum/mat"
)

func main() {
	// image plot
	file, _ := os.Open("pixels.csv")
	lines, _ := csv.NewReader(file).ReadAll()

	rChannel, gChannel, bChannel := []float64{}, []float64{}, []float64{}
	for _, line := range lines {
		r, _ := strconv.ParseFloat(line[0], 64)
		rChannel = append(rChannel, r)

		g, _ := strconv.ParseFloat(line[1], 64)
		gChannel = append(gChannel, g)

		b, _ := strconv.ParseFloat(line[2], 64)
		bChannel = append(bChannel, b)
	}

	x := make([]*mat.Dense, 3)
	x[0] = mat.NewDense(280, 280, rChannel)
	x[1] = mat.NewDense(280, 280, gChannel)
	x[2] = mat.NewDense(280, 280, bChannel)

	plt := plotter.NewPlot()
	plt.FigSize(9, 9)

	plt.ImShow(x)
	plt.Title("image plot example")

	plt.Save("image.png")
}
