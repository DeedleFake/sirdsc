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
// as a good default.
func Generate(out draw.Image, dm DepthMap, pat image.Image, partSize int) {
	if partSize <= 0 {
		partSize = pat.Bounds().Dx()
	}

	pat = TiledImage{
		Image: pat,
	}

	var wg sync.WaitGroup
	wg.Add(out.Bounds().Dy())
	for y := 0; y < out.Bounds().Dy(); y++ {
		go func(y int) {
			defer wg.Done()

			for x := 0; x < out.Bounds().Dx(); x++ {
				depth := dm.At(x-partSize, y)

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
