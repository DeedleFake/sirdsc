package sirdsc

import (
	"image"
	"image/color"
	"math"
)

// A DepthMap maps pixel coordinates to depths. If you would like to
// use a traditional grayscale image as a depth map, see the
// ImageDepthMap type.
type DepthMap interface {
	Bounds() image.Rectangle
	At(x, y int) int
}

const DefaultMaxImageDepth = 40

type ImageDepthMap struct {
	Image image.Image

	Max     int
	Flat    bool
	Inverse bool
}

func (dm ImageDepthMap) Bounds() image.Rectangle {
	return dm.Image.Bounds()
}

func (dm ImageDepthMap) At(x, y int) int {
	max := dm.Max
	if max <= 0 {
		max = DefaultMaxImageDepth
	}

	c := color.RGBAModel.Convert(dm.Image.At(x, y))
	tr, tg, tb, _ := c.RGBA()
	r := uint8(tr)
	g := uint8(tg)
	b := uint8(tb)

	v := math.Max(float64(g), math.Max(float64(b), float64(r)))
	d := v * float64(max) / math.MaxUint8

	if (dm.Flat) && (d != 0) {
		return max
	}

	if dm.Inverse {
		d = float64(max) - d
	}

	return int(d)
}
