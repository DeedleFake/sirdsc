package main

import(
	"os"
	"fmt"
	"flag"
	"rand"
	"math"
	"image"
	"image/jpeg"
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

func depthFromColor(c image.Color, max int, flat bool) (int) {
	c = image.RGBAColorModel.Convert(c)
	tr, tg, tb, _ := c.RGBA()
	r := uint8(tr)
	g := uint8(tg)
	b := uint8(tb)

	v := math.Fmax(float64(g), math.Fmax(float64(b), float64(r)))
	d := v * float64(max) / math.MaxUint8

	if (flat) && (d != 0) {
		return max
	}

	return int(d)
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

type Config struct {
	MaxDepth int
	Flat bool
}

func GenerateSIRDS(out *ImageFile, in *ImageFile, pat *ImageFile, config Config) {
	partSize := pat.Bounds().Dx()

	parts := in.Bounds().Dx() / partSize
	if (in.Bounds().Dx() % partSize) != 0 {
		parts++
	}

	for part := 0; part < parts + 1; part++ {
		for y := 0; y < out.Bounds().Dy(); y++ {
			for outX := part * partSize; outX < (part + 1) * partSize; outX++ {
				if outX > out.Bounds().Dx() {
					break
				}

				inX := outX - partSize
				depth := depthFromColor(in.At(inX, y), config.MaxDepth, config.Flat)

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
}

func main() {
	var(
		partSize int
		jpegOpt jpeg.Options
		config Config
	)
	flag.IntVar(&partSize, "partsize", 100, "Size of sections in the SIRDS")
	flag.IntVar(&config.MaxDepth, "depth", 40, "Maximum depth")
	flag.BoolVar(&config.Flat, "flat", false, "Generate a flat image")
	flag.IntVar(&jpegOpt.Quality, "jpeg:quality", 95, "Quality of output JPEG image")
	flag.Parse()
	args := flag.Args()
	if len(args) != 2 {
		usage(nil)
		os.Exit(1)
	}
	inFile := args[0]
	outFile := args[1]

	fmt.Printf("Options:\n")
	fmt.Printf("  depth: %v\n", config.MaxDepth)
	fmt.Printf("  flat: %v\n", config.Flat)
	fmt.Printf("  partsize: %v\n", partSize)
	fmt.Printf("  jpeg:quality: %v\n", jpegOpt.Quality)
	fmt.Printf("  src: %v\n", inFile)
	fmt.Printf("  dest: %v\n", outFile)
	fmt.Printf("\n")

	in, err := LoadImageFile(inFile)
	if err != nil {
		usage(err)
		os.Exit(1)
	}

	out, err := NewImageFile(outFile, in.Bounds().Dx() + partSize, in.Bounds().Dy())
	if err != nil {
		usage(err)
		os.Exit(1)
	}
	out.SetJPEGOptions(&jpegOpt)

	fmt.Printf("Generating Random Pattern...\n")
	pat, err := NewRandPat("", partSize, out.Bounds().Dy())
	if err != nil {
		usage(err)
		os.Exit(1)
	}

	fmt.Printf("Generating SIRDS...\n")
	GenerateSIRDS(out, in, pat, config)

	fmt.Printf("Writing SIRDS...\n")
	out.Save()

	fmt.Printf("Done.\n")
}
