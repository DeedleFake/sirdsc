package main

import(
	"os"
	"io"
	"image"
	"image/png"
	"image/jpeg"
)

type LoadFunc func(io.Reader) (image.Image, os.Error)
type SaveFunc func(io.Writer, image.Image, interface{}) (os.Error)

type ImageType struct {
	LoadFunc LoadFunc
	SaveFunc SaveFunc
}

func (it ImageType)Load(r io.Reader) (img image.Image, err os.Error) {
	if it.LoadFunc == nil {
		return nil, os.NewError("Image format can't be loaded...")
	}

	img, err = it.LoadFunc(r)

	return
}

func (it ImageType)Save(w io.Writer, img image.Image, data interface{}) (err os.Error) {
	if it.SaveFunc == nil {
		return os.NewError("Image format can't be saved to...")
	}

	err = it.SaveFunc(w, img, data)

	return
}

var JPG = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err os.Error) {
		img, err = jpeg.Decode(r)

		return
	},
	SaveFunc: func(w io.Writer, img image.Image, options interface{}) (err os.Error) {
		err = jpeg.Encode(w, img, options.(*jpeg.Options))

		return
	},
}

var PNG = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err os.Error) {
		img, err = png.Decode(r)

		return
	},
	SaveFunc: func(w io.Writer, img image.Image, ignored interface{}) (err os.Error) {
		err = png.Encode(w, img)

		return
	},
}
