package sirdsc

import (
	"image"
	"image/color"
	"math/rand"
	"sync"
)

var (
	randpool = sync.Pool{
		New: func() interface{} {
			return rand.New(rand.NewSource(1))
		},
	}
)

// A RandImage deterministically generates random color values for
// each (x, y) coordinate, using itself as a seed. In other words,
// given two RandImages that are equal to each other, the color at
// the same (x, y) in each are also equal.
type RandImage int64

func (img RandImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (img RandImage) Bounds() image.Rectangle {
	return image.Rect(-1e9, -1e9, 1e9, 1e9)
}

func (img RandImage) At(x, y int) color.Color {
	r := randpool.Get().(*rand.Rand)
	defer randpool.Put(r)

	r.Seed(int64(img) + int64(y))
	r.Seed(r.Int63() + int64(x))

	return color.RGBA{
		R: uint8(r.Uint32()),
		G: uint8(r.Uint32()),
		B: uint8(r.Uint32()),
		A: 255,
	}
}
