package sirdsc

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
)

// A Config specifies extra options for SIRDS generation.
type Config struct {
	MaxDepth int
	Flat     bool
	PartSize int
	Inverse  bool
}

var DefaultConfig = &Config{
	MaxDepth: 40,
	Flat:     false,
	PartSize: 100,
	Inverse:  false,
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

// Generate generates a new SIRDS from the depth map dm and draws it
// to out, using the pattern pat. If config is nil, DefaultConfig is
// used.
func Generate(out draw.Image, dm image.Image, pat image.Image, config *Config) {
	if config == nil {
		config = DefaultConfig
	}

	partSize := int(config.PartSize)
	if partSize <= 0 {
		partSize = pat.Bounds().Dx()
	}

	pat = TiledImage{
		Image: pat,
	}

	parts := dm.Bounds().Dx() / partSize
	if (dm.Bounds().Dx() % partSize) != 0 {
		parts++
	}

	var wg sync.WaitGroup
	wg.Add(out.Bounds().Dy())
	for y := 0; y < out.Bounds().Dy(); y++ {
		go func(y int) {
			defer wg.Done()

			for x := 0; x < out.Bounds().Dx(); x++ {
				depth := depthFromColor(dm.At(x-partSize, y), config.MaxDepth, config.Flat)
				if config.Inverse {
					depth = config.MaxDepth - depth
				}

				src := pat
				if x-partSize >= 0 {
					src = out
				}

				if x-depth >= 0 {
					c := src.At(x-partSize, y)

					out.Set(x, y, c)
					out.Set(x-depth, y, c)
				}
			}
		}(y)
	}
	wg.Wait()
}
