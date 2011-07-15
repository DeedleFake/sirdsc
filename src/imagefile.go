package main

import(
	"os"
	"io"
	"fmt"
	"path"
	"bytes"
	"strings"
	"image"
	"image/jpeg"
	"image/draw"
)

var(
	ETOLDNO = os.NewError("Told not to save file.")
)

type ImageFile struct {
	draw.Image

	FileType ImageType
	fileName string

	jpegOpt *jpeg.Options
}

func NewImageFile(file string, w, h int) (img *ImageFile, err os.Error) {
	img = new(ImageFile)

	err = img.SetFileName(file)
	if err != nil {
		return
	}

	if (w <= 0) || (h <= 0) {
		return nil, os.NewError("Bad dimensions")
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

	img.Image, err = img.FileType.Load(fl)
	if err != nil {
		return
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
			case ".bmp":
				img.FileType = BMP
			case ".gif":
				img.FileType = GIF
			case ".tif", ".tiff":
				img.FileType = TIFF
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

func (img *ImageFile)Save(askout io.Writer, askin io.Reader) (err os.Error) {
	if img.FileName() == "" {
		return os.NewError("No associated file...")
	}

	var rem bool
	_, err = os.Stat(img.FileName())
	if err != nil {
		err2, ok := err.(*os.PathError)
		if !ok {
			return
		}

		if err2.Error == os.ENOENT {
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
				return
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

func (img *ImageFile)SetJPEGOptions(opt *jpeg.Options) {
	img.jpegOpt = opt
}
