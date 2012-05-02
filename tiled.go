package main

import(
	"image"
	"image/color"
)

type TiledImage struct {
	image.Image
}

func (img *TiledImage)cX(x int, b image.Rectangle) int {
	if x < b.Min.X {
		x += b.Dx()
	} else if x > b.Max.X - 1 {
		x -= b.Dx()
	} else {
		return x
	}

	return img.cX(x, b)
}

func (img *TiledImage)cY(y int, b image.Rectangle) int {
	if y < b.Min.Y {
		y += b.Dy()
	} else if y > b.Max.Y - 1 {
		y -= b.Dy()
	} else {
		return y
	}

	return img.cY(y, b)
}

func (img *TiledImage)c(x, y int) (int, int) {
	b := img.Bounds()
	return img.cX(x, b), img.cY(y, b)
}

func (img *TiledImage)At(x, y int) color.Color {
	x, y = img.c(x, y)

	return img.Image.At(x, y)
}
