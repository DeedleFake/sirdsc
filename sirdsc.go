package sirdsc

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(int64(time.Since(time.Time{})))
}

func depthFromColor(c color.Color, max int, flat bool) int {
	c = color.RGBAModel.Convert(c)
	tr, tg, tb, _ := c.RGBA()
	r := uint8(tr)
	g := uint8(tg)
	b := uint8(tb)

	v := math.Max(float64(g), math.Max(float64(b), float64(r)))
	d := v * float64(max) / math.MaxUint8

	if (flat) && (d != 0) {
		return max
	}

	return int(d)
}

// A Config specifies extra options for SIRDS generation.
type Config struct {
	MaxDepth int
	Flat     bool
	PartSize int
}

var DefaultConfig = &Config{
	MaxDepth: 40,
	Flat:     false,
	PartSize: 100,
}

// Generate generates a new SIRDS from the depth map dm and draws it
// to out, using the pattern pat. If config is nil, DefaultConfig is
// used.
func Generate(out draw.Image, dm image.Image, pat image.Image, config *Config) {
	if config == nil {
		config = DefaultConfig
	}

	partSize := int(config.PartSize)
	if partSize == 0 {
		partSize = pat.Bounds().Dx()
	}

	patTile := TiledImage{
		Image:     pat,
		Rectangle: image.Rect(0, 0, out.Bounds().Dy(), config.PartSize),
	}

	parts := dm.Bounds().Dx() / partSize
	if (dm.Bounds().Dx() % partSize) != 0 {
		parts++
	}

	for part := 0; part < parts+1; part++ {
		for y := 0; y < out.Bounds().Dy(); y++ {
			for outX := part * partSize; outX < (part+1)*partSize; outX++ {
				if outX > out.Bounds().Dx() {
					break
				}

				inX := outX - partSize
				depth := depthFromColor(dm.At(inX, y), config.MaxDepth, config.Flat)

				if inX < 0 {
					if outX-depth >= 0 {
						out.Set(outX, y, patTile.At(outX, y))
						out.Set(outX-depth, y, patTile.At(outX, y))
					}
				} else {
					if outX-depth >= 0 {
						out.Set(outX, y, out.At(inX, y))
						out.Set(outX-depth, y, out.At(inX, y))
					}
				}
			}
		}
	}
}