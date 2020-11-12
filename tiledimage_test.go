package sirdsc_test

import (
	"image"
	"image/color"
	"testing"

	"github.com/DeedleFake/sirdsc"
)

type subImage struct {
	img  image.Image
	rect image.Rectangle
}

func (img subImage) ColorModel() color.Model {
	return img.img.ColorModel()
}

func (img subImage) Bounds() image.Rectangle {
	return img.rect
}

func (img subImage) At(x, y int) color.Color {
	return img.img.At(x, y)
}

func TestTiledImage(t *testing.T) {
	img := sirdsc.TiledImage{
		Image: subImage{
			img:  &sirdsc.RandImage{Seed: 1},
			rect: image.Rect(0, 0, 5, 5),
		},
	}

	c1 := img.At(0, 0)
	c2 := img.At(5, 5)
	if c1 != c2 {
		t.Fatalf("c1 == %#v\nc2 == %#v", c1, c2)
	}

	c1 = img.At(1, 1)
	c2 = img.At(11, -19)
	if c1 != c2 {
		t.Fatalf("c1 == %#v\nc2 == %#v", c1, c2)
	}
}
