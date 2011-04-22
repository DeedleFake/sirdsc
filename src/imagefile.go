package main

import(
	"os"
	"path"
	"reflect"
	"strings"
	"image"
	"image/png"
	"image/jpeg"
)

type ImageFile struct {
	image.Image

	FileType int
	fileName string
}

func NewImageFile(file string, w, h int) (img *ImageFile, err os.Error) {
	img = new(ImageFile)

	err = img.SetFileName(file)
	if err != nil {
		return
	}

	if (w <= 0) || (h <= 0) {
		err = os.NewError("Bad dimensions...")
		return
	}

	img.Image = image.NewRGBA(w, h)

	return
}

func LoadImageFile(file string) (img *ImageFile, err os.Error) {
	img = new(ImageFile)

	err = img.SetFileName(file)
	if err != nil {
		return
	}

	fl, err := os.Open(img.FileName())
	if err != nil {
		return
	}
	defer fl.Close()

	switch img.FileType {
		case JPG:
			tmpImage, err := jpeg.Decode(fl)
			if err != nil {
				return
			}
			img.Image = tmpImage
		case PNG:
			tmpImage, err := png.Decode(fl)
			if err != nil {
				return
			}
			img.Image = tmpImage
	}

	return
}

func CopyImageFile(src *ImageFile, slice image.Rectangle) (img *ImageFile, err os.Error) {
	img, err = NewImageFile(src.FileName(), slice.Dx(), slice.Dy())
	if err != nil {
		return
	}

	for inY := slice.Min.Y; inY < slice.Max.Y; inY++ {
		for inX := slice.Min.X; inX < slice.Max.X; inX++ {
			outX := inX - slice.Min.X
			outY := inY - slice.Min.Y

			img.Set(outX, outY, src.At(inX, inY))
		}
	}

	return
}

func NewRandPat(file string, w, h int) (img *ImageFile, err os.Error) {
	img, err = NewImageFile(file, w, h)
	if err != nil {
		return
	}

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := randomColor()

			img.Set(x, y, c)
		}
	}

	return
}

func (img *ImageFile)SetFileName(file string) (err os.Error) {
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

	img.fileName = file

	return nil
}

func (img *ImageFile)FileName() (string) {
	return img.fileName
}

func (img *ImageFile)Save() (err os.Error) {
	if img.FileName() == "" {
		return os.NewError("No associated file...")
	}

	fl, err := os.Create(img.FileName())
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

func (img *ImageFile)Set(x, y int, c image.Color) {
	v := reflect.NewValue(img.Image)

	c = img.Image.ColorModel().Convert(c)

	args := []reflect.Value{
		v,
		reflect.NewValue(x),
		reflect.NewValue(y),
	}

	t := v.Type()
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		if m.Name == "Set" {
			cv := reflect.Zero(m.Type.In(3))
			cv.Set(reflect.NewValue(c))
			args = append(args, cv)

			m.Func.Call(args)
			return
		}
	}

	panic("no 'Set' method")
}
