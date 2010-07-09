package main

import(
	"os"
	"fmt"
	"path"
	"rand"
	"strings"
	"strconv"
	"image"
	"image/jpeg"
	"image/png"
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

func makeRandPat(img *image.RGBA, x, y, w, h int) {
	for y < h {
		for x < w {
			c := image.RGBAColor{
				R: (uint8)rand.Int31n(255),
				G: (uint8)rand.Int31n(255),
				B: (uint8)rand.Int31n(255),
				A: 255,
			}

			img.Set(x, y, c)

			x++
		}

		y++
	}
}

func main() {
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}

	inN := os.Args[1]
	outN := os.Args[2]
	partSize, _ := strconv.Atoi(os.Args[3])

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

	makeRandPat(out)
	for part := 1; part < (in.Width() / partSize) + 1; part++ {
		for y := 0; y < out.Height(); y++ {
			for outX := part * partSize; outX < (part + 1) * partSize; outX++ {
				inX := outX - partSize

				out.Set(outX, y, in.At(inX, y))
			}
		}
	}

	fmt.Printf("Writing SIRDS...\n")
	png.Encode(outF, out)

	fmt.Printf("Done...\n")
}
