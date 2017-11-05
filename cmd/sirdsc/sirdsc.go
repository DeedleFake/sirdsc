package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"os"
	"time"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"

	"github.com/DeedleFake/sirdsc"
)

func loadImage(file string) (image.Image, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	return img, err
}

// TODO: Encode different types based on the extension.
func saveImage(file string, img image.Image) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <src>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	partSize := flag.Int("partsize", 100, "Size of sections in the SIRDS")
	maxDepth := flag.Int("depth", sirdsc.DefaultMaxImageDepth, "Maximum depth")
	flat := flag.Bool("flat", false, "Generate an image with only two planes")
	inverse := flag.Bool("inverse", false, "Treat darker pixels as closer in the depth map")
	seed := flag.Int64("seed", int64(time.Since(time.Time{})), "Color generation seed")
	patFile := flag.String("pat", "", "If not empty, use the specified file as the pattern instead of randomizing")
	outFile := flag.String("o", "sirds.png", "Output file")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	inFile := flag.Arg(0)

	inImg, err := loadImage(inFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %q: %v\n", inFile, err)
		os.Exit(1)
	}
	in := sirdsc.ImageDepthMap{
		Image: inImg,

		Max:     *maxDepth,
		Flat:    *flat,
		Inverse: *inverse,
	}

	pat := image.Image(sirdsc.RandImage(*seed))
	if *patFile != "" {
		pat, err = loadImage(*patFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open %q: %v\n", *patFile, err)
			os.Exit(1)
		}
	}

	inb := in.Bounds()
	out := image.NewNRGBA(image.Rect(
		inb.Min.X,
		inb.Min.Y,
		inb.Max.X+*partSize,
		inb.Max.Y,
	))
	sirdsc.Generate(out, in, pat, *partSize)

	err = saveImage(*outFile, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to %q: %v", *outFile, err)
		os.Exit(1)
	}
}
