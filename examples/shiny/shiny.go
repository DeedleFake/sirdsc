package main

import (
	"context"
	"image"
	"log"
	"time"

	"github.com/DeedleFake/sirdsc"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
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

type UpdateEvent struct {
	Time time.Time
}

func main() {
	driver.Main(func(s screen.Screen) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		w, err := s.NewWindow(&screen.NewWindowOptions{
			Width:  740,
			Height: 480,
		})
		if err != nil {
			log.Fatalf("Failed to create window: %v", err)
		}
		defer w.Release()

		go func() {
			t := time.NewTicker(time.Second / 60)
			defer t.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case ts := <-t.C:
					w.Send(UpdateEvent{
						Time: ts,
					})
				}
			}
		}()

		keys := make(map[key.Code]bool)
		keyDown := func(c key.Code) bool {
			d, ok := keys[c]
			return d && ok
		}

		out, err := s.NewBuffer(image.Pt(740, 480))
		if err != nil {
			log.Fatalf("Failed to create buffer: %v", err)
		}
		defer out.Release()

		dm := DepthMap{
			Rect: image.Rect(100, 100, 200, 200),
		}

		for {
			switch ev := w.NextEvent().(type) {
			case lifecycle.Event:
				if ev.To == lifecycle.StageDead {
					return
				}

			case key.Event:
				keys[ev.Code] = ev.Direction == key.DirPress

			case UpdateEvent:
				if keyDown(key.CodeUpArrow) {
					dm.Rect.Min.Y--
					dm.Rect.Max.Y--
				}
				if keyDown(key.CodeDownArrow) {
					dm.Rect.Min.Y++
					dm.Rect.Max.Y++
				}
				if keyDown(key.CodeLeftArrow) {
					dm.Rect.Min.X--
					dm.Rect.Max.X--
				}
				if keyDown(key.CodeRightArrow) {
					dm.Rect.Min.X++
					dm.Rect.Max.X++
				}

				sirdsc.Generate(out.RGBA(), dm, sirdsc.RandImage(1), 100)
				w.Upload(image.ZP, out, out.Bounds())
				w.Publish()
			}
		}
	})
}
