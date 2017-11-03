package sirdsc

import (
	"image"
	"image/color"
)

// A TiledImage tiles a region of an image. Any pixel values outside
// of the rectangle specified are wrapped so that they come from
// inside that rectangle in the original image.
type TiledImage struct {
	Image image.Image
	Rect  image.Rectangle
}

func (img TiledImage) c(x, y int) (int, int) {
	x = (x-img.Rect.Min.X)%img.Rect.Dx() + img.Rect.Min.X
	if x < 0 {
		x += img.Rect.Dx()
	}

	y = (y-img.Rect.Min.Y)%img.Rect.Dy() + img.Rect.Min.Y
	if y < 0 {
		y += img.Rect.Dy()
	}

	return x, y
}

func (img TiledImage) ColorModel() color.Model {
	return img.Image.ColorModel()
}

func (img TiledImage) Bounds() image.Rectangle {
	return image.Rect(-1e9, -1e9, 1e9, 1e9)
}

func (img TiledImage) At(x, y int) color.Color {
	x, y = img.c(x, y)
	return img.Image.At(x, y)
}
