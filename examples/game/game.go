package main

import (
	"image"
	"log"

	"github.com/DeedleFake/sirdsc"
	"github.com/DeedleFake/sirdsc/examples/game/sdl"
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

	win, err := sdl.CreateWindow(
		"SIRDS",
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		740,
		480,
		0,
	)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	defer win.Destroy()

	screen, err := win.GetSurface()
	if err != nil {
		log.Fatalf("Failed to create surface: %v", err)
	}

	keys := make(map[uint32]bool)
	keyDown := func(c uint32) bool {
		d, _ := keys[c]
		return d
	}

	dm := DepthMap{
		Rect: image.Rect(100, 100, 200, 200),
	}

	for {
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
			dm.Rect.Min.Y--
			dm.Rect.Max.Y--
		}
		if keyDown(sdl.K_DOWN) {
			dm.Rect.Min.Y++
			dm.Rect.Max.Y++
		}
		if keyDown(sdl.K_LEFT) {
			dm.Rect.Min.X--
			dm.Rect.Max.X--
		}
		if keyDown(sdl.K_RIGHT) {
			dm.Rect.Min.X++
			dm.Rect.Max.X++
		}

		sirdsc.Generate(screen, dm, sirdsc.RandImage(1), 100)
		win.UpdateSurface()
	}
}
