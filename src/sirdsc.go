package main

import(
	"os"
	"fmt"
	"flag"
	"rand"
	"math"
	"image"
)

const(
	JPG = iota + 1
	PNG
)

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %v [options] <src> <dest>\n", os.Args[0])
		fmt.Printf("\nOptions:\n")
		flag.PrintDefaults()
		fmt.Printf("\nParameters:\n")
		fmt.Printf("  src: Heightmap image\n")
		fmt.Printf("  dest: Output image file\n")
	}
}

func usage(err fmt.Stringer) {
	if err != nil {
		fmt.Printf("Error: %v\n", err) //\033[0;41m\033[m
		fmt.Printf("---------------------------\n\n")
	}

	flag.Usage()
}

//func colorsAreEqual(c, c2 image.Color) (bool) {
//	cR, cG, cB, cA := c.RGBA()
//	c2R, c2G, c2B, c2A := c2.RGBA()
//
//	if (cR == c2R) && (cG == c2G) && (cB == c2B) && (cA == c2A) {
//		return true
//	}
//
//	return false
//}

func depthFromColor(c image.Color, max int) (d int) {
	r, g, b, _ := c.RGBA()

	v := math.Fmax(float64(g), math.Fmax(float64(b), float64(r)))
	//d = int(v * max / float64(math.MaxUint32))
	if v != 0 {
		return max
	}

	return d
}

func randomColor() (image.Color) {
	c := image.RGBAColor{
		R: (uint8)(rand.Uint32()),
		G: (uint8)(rand.Uint32()),
		B: (uint8)(rand.Uint32()),
		A: 255,
	}

	return c
}

func main() {
	var(
		partSize int
		maxDepth int
	)
	flag.IntVar(&partSize, "partsize", 100, "Size of sections in the SIRDS")
	flag.IntVar(&maxDepth, "depth", 10, "Maximum depth")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		usage(nil)
		os.Exit(1)
	}
	inFile := args[0]
	outFile := args[1]

	fmt.Printf("Options:\n")
	fmt.Printf("  partsize: %v\n", partSize)
	fmt.Printf("  depth: %v\n", maxDepth)
	fmt.Printf("  src: %v\n", inFile)
	fmt.Printf("  dest: %v\n", outFile)
	fmt.Printf("\n")

	in, err := LoadImageFile(inFile)
	if err != nil {
		usage(err)
		os.Exit(1)
	}

	parts := in.Bounds().Dx() / partSize
	if (in.Bounds().Dx() % partSize) != 0 {
		parts++
	}

	out, err := NewImageFile(outFile, in.Bounds().Dx() + partSize, in.Bounds().Dy())
	if err != nil {
		usage(err)
		os.Exit(1)
	}

	pat, err := NewRandPat("", partSize, out.Bounds().Dy())
	if err != nil {
		usage(err)
		os.Exit(1)
	}

	fmt.Printf("Generating SIRDS...\n")
	for part := 0; part < parts + 1; part++ {
		for y := 0; y < out.Bounds().Dy(); y++ {
			for outX := part * partSize; outX < (part + 1) * partSize; outX++ {
				if outX > out.Bounds().Dx() {
					break
				}

				inX := outX - partSize
				depth := depthFromColor(in.At(inX, y), maxDepth)

				out.Set(outX, y, randomColor())

				if inX < 0 {
					if outX - depth >= 0 {
						out.Set(outX - depth, y, pat.At(outX, y))
					}
				} else {
					if outX - depth >= 0 {
						out.Set(outX - depth, y, out.At(inX, y))
					}
				}
			}
		}
	}

	fmt.Printf("Writing SIRDS...\n")
	out.Save()

	fmt.Printf("Done...\n")
}
