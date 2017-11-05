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
	// Bounds is the boundry of the depth map. This is analogous to
	// image.Image's Bounds() method.
	Bounds() image.Rectangle

	// At returns the depth at the given (x, y) coordinates. Zero is
	// considered to be the background plane, positive numbers are in
	// front of the background, and negative numbers are behind it.
	At(x, y int) int
}

// DefaultMaxImageDepth is the maximum depth used by ImageDepthMap if none is specified.
const DefaultMaxImageDepth = 40

// ImageDepthMap is a wrapper around an image.Image that allows it to
// be used as a DepthMap. It determines depth information from the
// values of the underlying pixels, considering higher value pixels to
// be closer and lower value pixels to be further away. Solid black
// pixels are considered to have a depth of zero. Alpha information is
// ignored.
type ImageDepthMap struct {
	// The image.Image to pull pixel data from.
	Image image.Image

	// Max is the maximum depth that can be calculated from the pixels.
	// In other words, a pixel with a value of 0 yields a depth of 0,
	// while a pixel with a value of 255 yields this depth.
	//
	// If Max is zero, DefaultMaxImageDepth is used instead.
	Max int

	// If flat is true, all non-black pixels are considered to be Max.
	Flat bool

	// If Inverse is true, pixel values are considered to be the inverse
	// of normal. In other words, lower value pixels are considered to
	// be closer, while higher value pixels are considered to be further
	// away.
	Inverse bool
}

// Bounds returns the same boundries as the underlying image.
func (dm ImageDepthMap) Bounds() image.Rectangle {
	return dm.Image.Bounds()
}

func (dm ImageDepthMap) At(x, y int) int { // nolint
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
