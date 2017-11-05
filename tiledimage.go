package sirdsc

import (
	"image"
	"image/color"
)

// A TiledImage extends another image by tiling it infinitely in every
// direction.
type TiledImage struct {
	image.Image
}

func (img TiledImage) c(x, y int) (int, int) {
	b := img.Image.Bounds()

	x = (x-b.Min.X)%b.Dx() + b.Min.X
	if x < 0 {
		x += b.Dx()
	}

	y = (y-b.Min.Y)%b.Dy() + b.Min.Y
	if y < 0 {
		y += b.Dy()
	}

	return x, y
}

func (img TiledImage) Bounds() image.Rectangle { // nolint
	return image.Rect(-1e9, -1e9, 1e9, 1e9)
}

func (img TiledImage) At(x, y int) color.Color { // nolint
	x, y = img.c(x, y)
	return img.Image.At(x, y)
}
