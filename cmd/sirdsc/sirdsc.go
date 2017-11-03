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
	var (
		config  sirdsc.Config
		seed    int64
		patFile string
		outFile string
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options] <src>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}
	flag.IntVar(&config.PartSize, "partsize", sirdsc.DefaultConfig.PartSize, "Size of sections in the SIRDS")
	flag.IntVar(&config.MaxDepth, "depth", sirdsc.DefaultConfig.MaxDepth, "Maximum depth")
	flag.BoolVar(&config.Flat, "flat", sirdsc.DefaultConfig.Flat, "Generate a flat image")
	flag.BoolVar(&config.Inverse, "inverse", sirdsc.DefaultConfig.Inverse, "Inverse depth math")
	flag.Int64Var(&seed, "seed", int64(time.Since(time.Time{})), "Color generation seed")
	flag.StringVar(&patFile, "pat", "", "Custom pattern")
	flag.StringVar(&outFile, "o", "sirds.png", "Output file")
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	inFile := flag.Arg(0)

	in, err := loadImage(inFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open %q: %v\n", inFile, err)
		os.Exit(1)
	}

	pat := image.Image(sirdsc.RandImage(seed))
	if patFile != "" {
		pat, err = loadImage(patFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open %q: %v\n", patFile, err)
			os.Exit(1)
		}
	}

	out := image.NewNRGBA(in.Bounds())
	sirdsc.Generate(out, in, pat, &config)

	err = saveImage(outFile, out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to %q: %v", outFile, err)
		os.Exit(1)
	}
}
