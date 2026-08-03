package plotter

import (
	"image"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"os"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
)

type animationOptions struct {
	nFrames  int
	update   func(frame int)
	interval time.Duration
	frame    int
	loop     bool
}

func (plt *plotParameters) Animation(nFrames int, update func(frame int), options ...func(*animationOptions)) {
	// default options
	plt.animation = &animationOptions{
		nFrames:  nFrames,
		update:   update,
		interval: 20 * time.Millisecond,
	}

	// apply additional options
	for _, option := range options {
		option(plt.animation)
	}
}

func WithAnimationInterval(interval time.Duration) func(*animationOptions) {
	return func(opts *animationOptions) {
		opts.interval = interval
	}
}

func WithAnimationLoop(loop bool) func(*animationOptions) {
	return func(opts *animationOptions) {
		opts.loop = loop
	}
}

func (plt *plotParameters) drawAnimationFrame(gtx layout.Context, e app.FrameEvent) {
	// restart animation if looping is enabled
	if plt.animation.frame >= plt.animation.nFrames && plt.animation.loop {
		plt.animation.frame = 0
		plt.Clear()
	}

	// update and draw the next animation frame
	if plt.animation.frame < plt.animation.nFrames {
		plt.animation.frame++
		plt.animation.update(plt.animation.frame)
		plt.DrawPlot()
	}

	// schedule the next animation frame
	if plt.animation.frame < plt.animation.nFrames || plt.animation.loop {
		gtx.Execute(op.InvalidateCmd{At: e.Now.Add(plt.animation.interval)})
	}
}

func (plt *plotParameters) saveAnimation(file string) {
	animation := gif.GIF{}

	// render and collect all animation frames
	for frame := 0; frame < plt.animation.nFrames; frame++ {
		plt.animation.update(frame + 1)
		plt.DrawPlot()

		// convert image to paletted image for GIF encoding
		img := plt.figure.Image()
		paletted := image.NewPaletted(img.Bounds(), palette.Plan9)
		draw.Draw(paletted, img.Bounds(), img, img.Bounds().Min, draw.Src)

		animation.Image = append(animation.Image, paletted)
		animation.Delay = append(animation.Delay, int(plt.animation.interval.Milliseconds()/10))
	}

	// save the animation to a GIF file
	w, err := os.Create(file)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	if err := gif.EncodeAll(w, &animation); err != nil {
		panic(err)
	}
}
