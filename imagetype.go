package main

import (
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

var (
	ENOLOAD = errors.New("Image format can't be loaded...")
	ENOSAVE = errors.New("Image format can't be saved to...")
)

type LoadFunc func(io.Reader) (image.Image, error)
type SaveFunc func(io.Writer, image.Image, interface{}) error

type ImageType struct {
	LoadFunc LoadFunc
	SaveFunc SaveFunc
}

func (it ImageType) Load(r io.Reader) (img draw.Image, err error) {
	if it.LoadFunc == nil {
		return nil, ENOLOAD
	}

	tmp, err := it.LoadFunc(r)
	if err != nil {
		return
	}

	if img, ok := tmp.(draw.Image); ok {
		return img, nil
	}
	return nil, errors.New("Couldn't load image: Does not satisfy draw.Image")
}

func (it ImageType) Save(w io.Writer, img draw.Image, data interface{}) (err error) {
	if it.SaveFunc == nil {
		return ENOSAVE
	}

	err = it.SaveFunc(w, img, data)

	return
}

var JPG = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err error) {
		img, err = jpeg.Decode(r)

		return
	},
	SaveFunc: func(w io.Writer, img image.Image, options interface{}) (err error) {
		err = jpeg.Encode(w, img, options.(*jpeg.Options))

		return
	},
}

var PNG = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err error) {
		img, err = png.Decode(r)

		return
	},
	SaveFunc: func(w io.Writer, img image.Image, ignored interface{}) (err error) {
		err = png.Encode(w, img)

		return
	},
}

var GIF = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err error) {
		img, err = gif.Decode(r)

		return
	},
}
