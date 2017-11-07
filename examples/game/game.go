package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"time"

	"github.com/DeedleFake/sirdsc"
	"github.com/DeedleFake/sirdsc/examples/game/internal/sdl"
)

type DepthMap struct {
	Rect image.Rectangle
}

func (dm DepthMap) Bounds() image.Rectangle {
	return image.Rect(0, 0, 640, 480)
}

func (dm DepthMap) At(x, y int) int {
	if image.Pt(x, y).In(dm.Rect) {
		return 10
	}
	return 0
}

func main() {
	err := sdl.Init()
	if err != nil {
		log.Fatalf("Failed to initialize SDL: %v", err)
	}
	defer sdl.Quit()

	win, ren, err := sdl.CreateWindowAndRenderer(740, 480, 0)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	defer ren.Destroy()
	defer win.Destroy()

	out, err := ren.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		740, 480,
	)
	if err != nil {
		log.Fatalf("Failed to create texture: %v", err)
	}
	defer out.Destroy()

	keys := make(map[uint32]bool)
	keyDown := func(c uint32) bool {
		d, _ := keys[c]
		return d
	}

	dm := DepthMap{
		Rect: image.Rect(100, 100, 200, 200),
	}

	tick := time.NewTicker(time.Second / 60)
	defer tick.Stop()

	var frames uint
	last := struct {
		ts    time.Time
		frame uint
	}{
		ts: time.Now(),
	}

	for start := range tick.C {
		if sec := start.Sub(last.ts).Seconds(); sec > 5 {
			fmt.Printf("FPS: %v\n", 1/(sec/float64(frames-last.frame)))

			last.ts = start
			last.frame = frames
		}

		for {
			ev := sdl.PollEvent()
			if ev == nil {
				break
			}

			switch ev := ev.(type) {
			case sdl.QuitEvent:
				return

			case sdl.KeyUpEvent:
				keys[ev.Keysym().Sym()] = false
			case sdl.KeyDownEvent:
				keys[ev.Keysym().Sym()] = true
			}
		}

		if keyDown(sdl.K_UP) {
			dm.Rect.Min.Y -= 10
			dm.Rect.Max.Y -= 10
		}
		if keyDown(sdl.K_DOWN) {
			dm.Rect.Min.Y += 10
			dm.Rect.Max.Y += 10
		}
		if keyDown(sdl.K_LEFT) {
			dm.Rect.Min.X -= 10
			dm.Rect.Max.X -= 10
		}
		if keyDown(sdl.K_RIGHT) {
			dm.Rect.Min.X += 10
			dm.Rect.Max.X += 10
		}

		img := out.Image()
		sirdsc.Generate(img, dm, sirdsc.RandImage(rand.Uint64()), 100)
		img.Close()

		ren.Copy(out, image.ZR, image.ZR)
		ren.Present()

		frames++
	}
}
