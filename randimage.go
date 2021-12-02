package sirdsc

import (
	"image"
	"image/color"

	"github.com/DeedleFake/sirdsc/spcg"
)

// A RandImage deterministically generates random color values for
// each (x, y) coordinate, using itself as a seed. In other words,
// given two RandImages that are equal to each other, the color at
// the same (x, y) in each are also equal.
//
// A RandImage has infinite size.
type RandImage struct {
	Seed uint64
}

func (img RandImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (img RandImage) Bounds() image.Rectangle {
	return image.Rect(-1e9, -1e9, 1e9, 1e9)
}

func (img RandImage) At(x, y int) color.Color {
	c, _, _ := spcg.Next(uint64(x)^img.Seed, uint64(y)^img.Seed)

	return color.RGBA{
		R: uint8(c),
		G: uint8(c >> 8),
		B: uint8(c >> 16),
		A: 255,
	}
}

// A SymmetricRandImage is a variant of RandImage that is symmetric,
// such that a pixel at (a, b) is equal to a pixel at (b, a).
type SymmetricRandImage struct {
	Seed uint64
}

func (img SymmetricRandImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (img SymmetricRandImage) Bounds() image.Rectangle {
	return image.Rect(-1e9, -1e9, 1e9, 1e9)
}

func (img SymmetricRandImage) At(x, y int) color.Color {
	c, _, _ := spcg.Next(uint64(x^y)^img.Seed, uint64(x^y)^img.Seed)

	return color.RGBA{
		R: uint8(c),
		G: uint8(c >> 8),
		B: uint8(c >> 16),
		A: 255,
	}
}
