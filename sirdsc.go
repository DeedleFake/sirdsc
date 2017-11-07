package sirdsc

import (
	"image"
	"image/draw"
	"sync"
)

// Generate generates a new SIRDS from the depth map dm and draws it
// to out, using the pattern pat. partSize specifies the width of a
// single section of the generated stereogram. If partSize is less
// than or equal to zero, the width of pat is used. 100 is recommended
// as a good default, but this heavily depends on the physical size of
// the screen that the stereogram will be displayed on.
func Generate(out draw.Image, dm DepthMap, pat image.Image, partSize int) {
	if partSize <= 0 {
		partSize = pat.Bounds().Dx()
	}

	pat = TiledImage{
		Image: pat,
	}

	b := out.Bounds()

	var wg sync.WaitGroup
	wg.Add(b.Dy())
	for y := b.Min.Y; y < b.Max.Y; y++ {
		go func(y int) {
			defer wg.Done()

			for x := b.Min.X; x < b.Max.X; x++ {
				depth := dm.At(x-partSize, y)

				src := pat
				if x-partSize >= 0 {
					src = out
				}

				c := src.At(x-partSize, y)
				out.Set(x, y, c)

				if (depth != 0) && (x-depth >= b.Min.X) && (x-depth <= b.Max.X) {
					out.Set(x-depth, y, c)
				}
			}
		}(y)
	}
	wg.Wait()
}
