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

const(
	JPG = iota
	PNG = iota
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

type imageData struct {
	*image.RGBA

	FileType int
	FileName string
}

func NewImageData(file string, dw, dh int) (img imageData, err os.Error) {
	err = img.SetFileName(file)
	if err != nil {
		return img, err
	}

	_, err = os.Lstat(img.FileName)
	if (err == nil) && (file != "") {
		fl, err := os.Open(img.FileName, os.O_RDONLY, 0666)
		if err != nil {
			return img, err
		}
		defer fl.Close()

		switch img.FileType {
			case JPG:
				tmpImage, err := jpeg.Decode(fl)
				if err != nil {
					return img, err
				}
				img.RGBA = tmpImage.(*image.RGBA)
			case PNG:
				tmpImage, err := png.Decode(fl)
				if err != nil {
					return img, err
				}
				img.RGBA = tmpImage.(*image.RGBA)
		}
	} else {
		if (dw > 0) && (dh > 0) {
			img.RGBA = image.NewRGBA(dw, dh)
		} else {
			return img, err
		}
	}

	return img, nil
}

func (img *imageData)SetFileName(file string) (err os.Error) {
	if file != "" {
		switch strings.ToLower(path.Ext(file)) {
			case ".jpg", ".jpeg":
				img.FileType = JPG
			case ".png":
				img.FileType = PNG
			default:
				return os.NewError("Image format not supported or could not be detected...")
		}
	}

	img.FileName = file

	return nil
}

func (img *imageData)Save() (err os.Error) {
	fl, err := os.Open(img.FileName, os.O_RDWR | os.O_CREAT | os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer fl.Close()

	switch img.FileType {
		case PNG:
			png.Encode(fl, img)
		default:
			return os.NewError("Image format can't be saved to...")
	}

	return nil
}

func (img *imageData)MakeRandPat(x, y, w, h int) {
	sx := x
	sy := y

	for y = sy; y < h; y++ {
		for x = sx; x < w; x++ {
			c := randomColor()

			img.Set(x, y, c)
		}
	}
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

func copyAndCheckPixel(in *image.RGBA, in2 *image.RGBA, inX, inY int, out *image.RGBA, outX, outY int) {
	if !colorsAreEqual(in.At(inX, inY), image.Black) {
		depth := depthFromColor(in.At(inX, inY))

		out.Set(outX - depth, outY, in2.At(inX, inY))

		for i := 0; i < depth; i++ {
			out.Set(outX - i, outY, randomColor())
		}
	} else {
		out.Set(outX, outY, in2.At(inX, inY))
	}
}

func main() {
	if len(os.Args) != 4 {
		usage()
		os.Exit(1)
	}

	partSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		usageE("Could not convert third argument to int...")
		os.Exit(1)
	}

	in, err := NewImageData(os.Args[1], 0, 0)
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}

	out, err := NewImageData("", in.Width(), in.Height())
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}
	err = out.SetFileName(os.Args[2])
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}

	pat, err := NewImageData("", partSize, out.Height())
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}

	fmt.Printf("Generating SIRDS...\n")
	pat.MakeRandPat(0, 0, pat.Width(), pat.Height())
	for part := 1; part < (in.Width() / partSize); part++ {
		for y := 0; y < out.Height(); y++ {
			for outX := part * partSize; outX < (part + 1) * partSize; outX++ {
				inX := outX - partSize

				if inX < 0 {
					copyAndCheckPixel(in.RGBA, pat.RGBA, outX, y, out.RGBA, outX, y)
				} else {
					copyAndCheckPixel(in.RGBA, out.RGBA, inX, y, out.RGBA, outX, y)
				}
			}
		}
	}

	fmt.Printf("Writing SIRDS...\n")
	out.Save()

	fmt.Printf("Done...\n")
}
