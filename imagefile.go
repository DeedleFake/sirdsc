package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"os"
	"path"
	"strings"
)

var (
	ETOLDNO = errors.New("Told not to save file.")
)

type ImageFile struct {
	draw.Image

	FileType ImageType
	fileName string

	jpegOpt *jpeg.Options
}

func NewImageFile(file string, w, h int) (img *ImageFile, err error) {
	img = new(ImageFile)

	err = img.SetFileName(file)
	if err != nil {
		return
	}

	if (w <= 0) || (h <= 0) {
		return nil, errors.New("Bad dimensions")
	}

	img.Image = image.NewRGBA(image.Rect(0, 0, w, h))

	return
}

func LoadImageFile(file string) (img *ImageFile, err error) {
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

	img.Image, err = img.FileType.Load(fl)
	if err != nil {
		return
	}

	return
}

func CopyImageFile(src *ImageFile, slice image.Rectangle) (img *ImageFile, err error) {
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

func NewRandPat(file string, w, h int) (img *ImageFile, err error) {
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

func (img *ImageFile) SetFileName(file string) (err error) {
	if file != "" {
		switch strings.ToLower(path.Ext(file)) {
		case ".jpg", ".jpeg":
			img.FileType = JPG
		case ".png":
			img.FileType = PNG
		case ".gif":
			img.FileType = GIF
		default:
			return errors.New("Image format not supported or could not be detected...")
		}
	}

	img.fileName = file

	return nil
}

func (img *ImageFile) FileName() string {
	return img.fileName
}

func (img *ImageFile) Save(askout io.Writer, askin io.Reader) (err error) {
	if img.FileName() == "" {
		return errors.New("No associated file...")
	}

	var rem bool
	_, err = os.Stat(img.FileName())
	if err != nil {
		err2, ok := err.(*os.PathError)
		if !ok {
			return
		}

		if os.IsNotExist(err2) {
			rem = true
		}
	} else {
		if (askin != nil) && (askout != nil) {
			_, err = fmt.Fprintf(askout, "Overwrite %v? [y/N] ", img.FileName())
			if err != nil {
				return
			}

			ans := make([]byte, 1)
			n, err := askin.Read(ans)
			if err != nil {
				return err
			}
			if (n == 0) || (bytes.ToLower(ans)[0] != 'y') {
				return ETOLDNO
			}
		}
	}

	fl, err := os.Create(img.FileName())
	if err != nil {
		return
	}
	defer fl.Close()

	err = img.FileType.Save(fl, img, img.jpegOpt)
	if err != nil {
		if (err == ENOSAVE) && rem {
			os.Remove(img.FileName())
		}

		return
	}

	return nil
}

func (img *ImageFile) SetJPEGOptions(opt *jpeg.Options) {
	img.jpegOpt = opt
}
