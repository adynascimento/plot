package main

import (
	"image"
	"log"
	"os"

	"github.com/adynascimento/plot/plot"
)

func main() {
	// image plot
	file, err := os.Open("image.png")
	if err != nil {
		log.Println("error openning image:", err.Error())
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Println("error decodinng image:", err.Error())
	}

	plt := plot.NewPlot()
	plt.FigSize(18, 10)

	plt.ImShow(img)
	plt.Title("image plot example")

	plt.Save("image.png")
}
