package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"time"

	"github.com/DeedleFake/sirdsc"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const (
	// ScreenWidth is the width of the underlying depth map. The actual
	// display will be wider by PartSize pixels.
	ScreenWidth = 640

	// ScreenHeight is the height of both the depth map and the display.
	ScreenHeight = 480

	// PartSize is the size of a part of the stererogram.
	PartSize = 100

	// FPSDelay is the duration to wait in between each printing of the
	// FPS.
	FPSDelay = 5 * time.Second
)

// DepthMap is an implementation of sirdsc.DepthMap that 'draws'
// objects as depths, rather than colors.
type DepthMap struct {
	Depth int
	Rect  pixel.Rect

	Obstacle pixel.Rect
}

func (dm DepthMap) Bounds() image.Rectangle { // nolint
	return image.Rect(0, 0, ScreenWidth, ScreenHeight)
}

func (dm DepthMap) At(x, y int) int { // nolint
	v := pixel.V(float64(x), float64(y))

	if dm.Rect.Contains(v) && dm.Obstacle.Contains(v) {
		return int(math.Max(float64(dm.Depth), 10))
	}
	if dm.Rect.Contains(v) {
		return dm.Depth
	}
	if dm.Obstacle.Contains(v) {
		return 10
	}

	return 0
}

type PictureImage pixel.PictureData

func (img PictureImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (img PictureImage) Bounds() image.Rectangle {
	return image.Rect(
		int(img.Rect.Min.X),
		int(img.Rect.Min.Y),
		int(img.Rect.Max.X),
		int(img.Rect.Max.Y),
	)
}

func (img PictureImage) index(x, y int) int {
	y = int(img.Rect.Max.Y) - 1 - y
	return (y * img.Stride) + x
}

func (img PictureImage) At(x, y int) color.Color {
	return img.Pix[img.index(x, y)]
}

func (img PictureImage) Set(x, y int, c color.Color) {
	i := img.index(x, y)
	img.Pix[i] = img.ColorModel().Convert(c).(color.RGBA)
}

func (img *PictureImage) PictureData() *pixel.PictureData {
	return (*pixel.PictureData)(img)
}

func main() {
	pixelgl.Run(func() {
		win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
			Title:  "SIRDS",
			Bounds: pixel.R(0, 0, ScreenWidth+PartSize, ScreenHeight),
			VSync:  true,
		})
		if err != nil {
			log.Fatalf("Failed to create window: %v", err)
		}
		defer win.Destroy()
		win.SetMatrix(pixel.IM.Moved(win.Bounds().Center()))

		out := (*PictureImage)(pixel.MakePictureData(win.Bounds()))

		dm := DepthMap{
			Depth: 10,
			Rect:  pixel.R(100, 100, 200, 200),

			Obstacle: pixel.R(
				ScreenWidth/2-35,
				ScreenHeight/2-35,
				ScreenWidth/2+35,
				ScreenHeight/2+35,
			),
		}

		seed := time.Now().UnixNano()

		tick := time.NewTicker(time.Second / 60)
		defer tick.Stop()

		var frames uint
		last := struct {
			ts    time.Time
			frame uint
		}{
			ts: time.Now(),
		}

		for !win.Closed() {
			// TODO: Scale movement speed.
			start := time.Now()
			fps := float64(frames-last.frame) / start.Sub(last.ts).Seconds()
			if start.Sub(last.ts) > FPSDelay {
				fmt.Printf("FPS: %v\n", fps)

				last.ts = start
				last.frame = frames
			}

			if win.Pressed(pixelgl.KeyUp) {
				dm.Rect.Min.Y -= 10
				dm.Rect.Max.Y -= 10
			}
			if win.Pressed(pixelgl.KeyDown) {
				dm.Rect.Min.Y += 10
				dm.Rect.Max.Y += 10
			}
			if win.Pressed(pixelgl.KeyLeft) {
				dm.Rect.Min.X -= 10
				dm.Rect.Max.X -= 10
			}
			if win.Pressed(pixelgl.KeyRight) {
				dm.Rect.Min.X += 10
				dm.Rect.Max.X += 10
			}

			if win.Pressed(pixelgl.KeyW) {
				dm.Depth--
			}
			if win.Pressed(pixelgl.KeyS) {
				dm.Depth++
			}
			if dm.Depth < 5 {
				dm.Depth = 5
			}
			if dm.Depth > 20 {
				dm.Depth = 20
			}

			if s := time.Now().UnixNano(); s-seed > int64(time.Second/30) {
				seed = s
			}

			sirdsc.Generate(out, dm, sirdsc.RandImage{Seed: uint64(seed)}, PartSize)
			s := pixel.NewSprite(out.PictureData(), out.PictureData().Bounds())
			s.Draw(win, pixel.IM)

			win.Update()

			frames++
		}
	})
}
