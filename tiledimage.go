package sirdsc

import (
	"image"
	"image/color"
)

// A TiledImage tiles a region of an image. Any pixel values outside
// of the rectangle specified are wrapped so that they come from
// inside that rectangle in the original image.
type TiledImage struct {
	image.Image
	image.Rectangle
}

func (img TiledImage) c(x, y int) (int, int) {
	x = (x-img.Min.X)%img.Dx() + img.Min.X
	y = (y-img.Min.Y)%img.Dy() + img.Min.Y
	return x, y
}

func (img TiledImage) Bounds() image.Rectangle {
	return image.Rect(-1e9, -1e9, 1e9, 1e9)
}

func (img TiledImage) At(x, y int) color.Color {
	x, y = img.c(x, y)
	return img.Image.At(x, y)
}
