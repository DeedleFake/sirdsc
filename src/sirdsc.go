package main

import(
	"os"
	"fmt"
	"path"
	"rand"
	"math"
	"strings"
	"strconv"
	"image"
	"image/jpeg"
	"image/png"
)

const(
	MaxDepth = 10
)

func usageE(err string) {
	fmt.Printf("\033[0;41mError:\033[m %s\n", err)
	fmt.Printf("---------------------------\n\n")

	usage()
}

func usage() {
	fmt.Printf("\033[0;1mUsage:\033[m\n")
	fmt.Printf("\t%s <in> <out> <part-size>\n\n", os.Args[0])
	fmt.Printf("\t\tin: A Jpeg or PNG file.\n")
	fmt.Printf("\t\tout: A PNG file.\n")
	fmt.Printf("\t\tpart-size: The width of each section of the SIRDS.\n")
}

func colorsAreEqual(c, c2 image.Color) (bool) {
	cR, cG, cB, cA := c.RGBA()
	c2R, c2G, c2B, c2A := c2.RGBA()

	if (cR == c2R) && (cG == c2G) && (cB == c2B) && (cA == c2A) {
		return true
	}

	return false
}/**/

func depthFromColor(c image.Color) (d int) {
	r, g, b, _ := c.RGBA()

	r = uint32(float(r * MaxDepth) / math.MaxUint32)
	g = uint32(float(g * MaxDepth) / math.MaxUint32)
	b = uint32(float(b * MaxDepth) / math.MaxUint32)

	d = int(float((r + g + b) * MaxDepth) / (MaxDepth * 3))

	return 5
}

func randomColor() (image.Color) {
	c := image.RGBAColor{
		R: (uint8)(rand.Int31n(255)),
		G: (uint8)(rand.Int31n(255)),
		B: (uint8)(rand.Int31n(255)),
		A: 255,
	}

	return c
}

func makeRandPat(img *image.RGBA, x, y, w, h int) {
	sx := x
	sy := y

	for y = sy; y < h; y++ {
		for x = sx; x < w; x++ {
			c := randomColor()

			img.Set(x, y, c)
		}
	}
}

func main() {
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}

	inN := os.Args[1]
	outN := os.Args[2]
	partSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		usageE("Could not convert third argument to int...")
		os.Exit(1)
	}

	switch strings.ToLower(path.Ext(outN)) {
		//case ".jpg", ".jpeg":
		case ".png":
		default:
			usageE("Output image format not supported...")
			os.Exit(1)
	}
	outF, err := os.Open(outN, os.O_RDWR | os.O_CREAT | os.O_TRUNC, 0666)
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}
	defer outF.Close()

	inF, err := os.Open(inN, os.O_RDONLY, 0666)
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}
	defer inF.Close()

	var in image.Image
	switch strings.ToLower(path.Ext(inN)) {
		case ".jpg", ".jpeg":
			fmt.Printf("Loading Jpeg...\n")
			in, err = jpeg.Decode(inF)
			if err != nil {
				usageE(err.String())
				os.Exit(1)
			}
		case ".png":
			fmt.Printf("Loading PNG...\n")
			in, err = png.Decode(inF)
			if err != nil {
				usageE(err.String())
				os.Exit(1)
			}
		default:
			usageE("Input image format either not supported or could not be detected...")
			os.Exit(1)
	}

	out := image.NewRGBA(in.Width() + partSize, in.Height())

	fmt.Printf("Generating SIRDS...\n")
	makeRandPat(out, 0, 0, partSize, out.Height())
	for part := 1; part < (in.Width() / partSize) + 1; part++ {
		for y := 0; y < out.Height(); y++ {
			for outX := part * partSize; outX < (part + 1) * partSize; outX++ {
				inX := outX - partSize

				if !colorsAreEqual(in.At(inX, y), image.Black) {
					depth := depthFromColor(in.At(inX, y))

					out.Set(outX - depth, y, out.At(inX, y))

					for i := 0; i < depth; i++ {
						out.Set(outX - i, y, randomColor())
					}
				} else {
					out.Set(outX, y, out.At(inX, y))
				}
			}
		}
	}

	fmt.Printf("Writing SIRDS...\n")
	png.Encode(outF, out)

	fmt.Printf("Done...\n")
}
