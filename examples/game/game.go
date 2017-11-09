package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/DeedleFake/sirdsc"
	"github.com/DeedleFake/sirdsc/examples/game/internal/sdl"
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
	Rect  image.Rectangle

	Obstacle image.Rectangle
}

func (dm DepthMap) Bounds() image.Rectangle { // nolint
	return image.Rect(0, 0, ScreenWidth, ScreenHeight)
}

func (dm DepthMap) At(x, y int) int { // nolint
	p := image.Pt(x, y)

	if p.In(dm.Rect) && p.In(dm.Obstacle) {
		return int(math.Max(float64(dm.Depth), 10))
	}
	if p.In(dm.Rect) {
		return dm.Depth
	}
	if p.In(dm.Obstacle) {
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

	win, ren, err := sdl.CreateWindowAndRenderer(ScreenWidth+PartSize, ScreenHeight, 0)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	defer ren.Destroy()
	defer win.Destroy()

	out, err := ren.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		ScreenWidth+PartSize, ScreenHeight,
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
		Depth: 10,
		Rect:  image.Rect(100, 100, 200, 200),

		Obstacle: image.Rect(
			ScreenWidth/2-35,
			ScreenHeight/2-35,
			ScreenWidth/2+35,
			ScreenHeight/2+35,
		),
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
		// TODO: Scale movement speed.
		fps := float64(frames-last.frame) / start.Sub(last.ts).Seconds()
		if start.Sub(last.ts) > FPSDelay {
			fmt.Printf("FPS: %v\n", fps)

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

		if keyDown(sdl.K_w) {
			dm.Depth--
		}
		if keyDown(sdl.K_s) {
			dm.Depth++
		}
		if dm.Depth < 5 {
			dm.Depth = 5
		}
		if dm.Depth > 20 {
			dm.Depth = 20
		}

		img := out.Image()
		sirdsc.Generate(img, dm, sirdsc.RandImage(rand.Uint64()), PartSize)
		img.Close()

		ren.Copy(out, image.ZR, image.ZR)
		ren.Present()

		frames++
	}
}
