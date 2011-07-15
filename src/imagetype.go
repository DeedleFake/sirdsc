package main

import(
	"os"
	"io"
	"image"
	"image/bmp"
	"image/gif"
	"image/png"
	"image/jpeg"
	"image/tiff"
	"image/draw"
)

var(
	ENOLOAD = os.NewError("Image format can't be loaded...")
	ENOSAVE = os.NewError("Image format can't be saved to...")
)

type LoadFunc func(io.Reader) (image.Image, os.Error)
type SaveFunc func(io.Writer, image.Image, interface{}) (os.Error)

type ImageType struct {
	LoadFunc LoadFunc
	SaveFunc SaveFunc
}

func (it ImageType)Load(r io.Reader) (img draw.Image, err os.Error) {
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
	return nil, os.NewError("Couldn't load image: Does not satisfy draw.Image")
}

func (it ImageType)Save(w io.Writer, img draw.Image, data interface{}) (err os.Error) {
	if it.SaveFunc == nil {
		return ENOSAVE
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

var BMP = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err os.Error) {
		img, err = bmp.Decode(r)

		return
	},
}

var GIF = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err os.Error) {
		img, err = gif.Decode(r)

		return
	},
}

var TIFF = ImageType{
	LoadFunc: func(r io.Reader) (img image.Image, err os.Error) {
		img, err = tiff.Decode(r)

		return
	},
}
