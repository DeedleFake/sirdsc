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
	"reflect"
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

type ImageData struct {
	image.Image

	FileType int
	FileName string
}

func NewImageData(file string, dw, dh int) (img *ImageData, err os.Error) {
	img = new(ImageData)

	err = img.SetFileName(file)
	if err != nil {
		return img, err
	}

	_, err = os.Lstat(img.FileName)
	if (err == nil) && (file != "") {
		fl, err := os.Open(img.FileName)
		if err != nil {
			return img, err
		}
		defer fl.Close()

		var tmpImage image.Image
		switch img.FileType {
			case JPG:
				tmpImage, err = jpeg.Decode(fl)
				if err != nil {
					return img, err
				}
				img.Image = tmpImage
			case PNG:
				tmpImage, err = png.Decode(fl)
				if err != nil {
					return img, err
				}
				img.Image = tmpImage
		}
	} else {
		if (dw > 0) && (dh > 0) {
			img.Image = image.NewRGBA(dw, dh)
		} else {
			return img, err
		}
	}

	return img, nil
}

func (img *ImageData)SetFileName(file string) (err os.Error) {
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

func (img *ImageData)Save() (err os.Error) {
	fl, err := os.Create(img.FileName)
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

func (img *ImageData)MakeRandPat(x, y, w, h int) {
	sx := x
	sy := y

	for y = sy; y < h; y++ {
		for x = sx; x < w; x++ {
			cv := reflect.NewValue(img.Image.At(0, 0))
			if cv.Type().Kind() == reflect.Interface {
				cv = reflect.Indirect(cv.Elem())
			}

			c := randomColor(cv.Type())

			img.Set(x, y, c)
		}
	}
}

func (img *ImageData)Set(x, y int, c image.Color) {
	v := reflect.NewValue(img.Image)
	if v.Type().Kind() == reflect.Interface {
		v = reflect.Indirect(v.Elem())
	}

	cv := reflect.NewValue(c)
	if cv.Type().Kind() == reflect.Interface {
		cv = reflect.Indirect(cv.Elem())
	}

	args := []reflect.Value{
		v,
		reflect.NewValue(x),
		reflect.NewValue(y),
		cv,
	}

	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.Name == "Set" {
			m.Func.Call(args)
			return
		}
	}

	panic("no 'Set' method")
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

func randomColor(t reflect.Type) (image.Color) {
	c := reflect.Zero(t)
	c.FieldByName("R").Set(reflect.NewValue((uint8)(rand.Uint32())))
	c.FieldByName("G").Set(reflect.NewValue((uint8)(rand.Uint32())))
	c.FieldByName("B").Set(reflect.NewValue((uint8)(rand.Uint32())))
	c.FieldByName("A").Set(reflect.NewValue(255))

	return c.Interface().(image.Color)
}

func copyAndCheckPixel(in *ImageData, in2 *ImageData, inX, inY int, out *ImageData, outX, outY int) {
	depth := depthFromColor(in.At(inX, inY))

	cv := reflect.NewValue(out.Image.At(0, 0))
	if cv.Type().Kind() == reflect.Interface {
		cv = reflect.Indirect(cv.Elem())
	}

	out.Set(outX, outY, randomColor(cv.Type()))

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
		usageE("Could not convert third argument to int...")
		os.Exit(1)
	}

	in, err := NewImageData(os.Args[1], 0, 0)
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}

	out, err := NewImageData("", in.Bounds().Dx(), in.Bounds().Dy())
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}
	err = out.SetFileName(os.Args[2])
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}

	pat, err := NewImageData("", partSize, out.Bounds().Dy())
	if err != nil {
		usageE(err.String())
		os.Exit(1)
	}

	fmt.Printf("Generating SIRDS...\n")
	pat.MakeRandPat(0, 0, pat.Bounds().Dx(), pat.Bounds().Dy())
	for part := 0; part < (in.Bounds().Dx() / partSize); part++ {
		for y := 0; y < out.Bounds().Dy(); y++ {
			for outX := part * partSize; outX < (part + 1) * partSize; outX++ {
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
