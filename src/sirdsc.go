package main

import(
	"os"
	"fmt"
	"rand"
	"math"
	"image"
	"strconv"
)

const(
	MaxDepth = 10
)

const(
	JPG = iota + 1
	PNG
)

func usageE(err fmt.Stringer) {
	fmt.Printf("\033[0;41mError:\033[m %s\n", err)
	fmt.Printf("---------------------------\n\n")

	usage()
}

func usage() {
	fmt.Printf("\033[0;1mUsage:\033[m\n")
	fmt.Printf("\t%s <src> <dst> <part-size>\n\n", os.Args[0])
	fmt.Printf("\t\tsrc: A Jpeg or PNG file.\n")
	fmt.Printf("\t\tdst: A PNG file.\n")
	fmt.Printf("\t\tpart-size: The width of each section of the SIRDS.\n")
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

func depthFromColor(c image.Color) (d int) {
	r, g, b, _ := c.RGBA()

	v := math.Fmax(float64(g), math.Fmax(float64(b), float64(r)))
	//d = int(v * MaxDepth / float64(math.MaxUint32))
	if v != 0 {
		return 5
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

func copyAndCheckPixel(in *ImageFile, in2 *ImageFile, inX, inY int, out *ImageFile, outX, outY int) {
	depth := depthFromColor(in.At(inX, inY))

	out.Set(outX, outY, randomColor())

	if outX - depth >= 0 {
		out.Set(outX - depth, outY, in2.At(inX, inY))
	}
}

func main() {
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}

	partSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		usageE(os.NewError("Could not convert third argument to int."))
		os.Exit(1)
	}

	in, err := LoadImageFile(os.Args[1])
	if err != nil {
		usageE(err)
		os.Exit(1)
	}

	parts := in.Bounds().Dx() / partSize
	if (in.Bounds().Dx() % partSize) != 0 {
		parts++
	}

	out, err := NewImageFile(os.Args[2], in.Bounds().Dx(), in.Bounds().Dy())
	if err != nil {
		usageE(err)
		os.Exit(1)
	}

	pat, err := NewRandPat("", partSize, out.Bounds().Dy())
	if err != nil {
		usageE(err)
		os.Exit(1)
	}

	fmt.Printf("Generating SIRDS...\n")
	for part := 0; part < parts; part++ {
		for y := 0; y < out.Bounds().Dy(); y++ {
			for outX := part * partSize; outX < (part + 1) * partSize; outX++ {
				if outX > out.Bounds().Dx() {
					break
				}

				inX := outX - partSize

				if inX < 0 {
					copyAndCheckPixel(in, pat, outX, y, out, outX, y)
				} else {
					copyAndCheckPixel(in, out, inX, y, out, outX, y)
				}
			}
		}
	}

	fmt.Printf("Writing SIRDS...\n")
	out.Save()

	fmt.Printf("Done...\n")
}
