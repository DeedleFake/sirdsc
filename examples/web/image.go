package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/DeedleFake/sirdsc"
	"golang.org/x/sync/errgroup"
)

type Image interface {
	image.Image
	Generate(ctx context.Context, w io.Writer, pat Image, q url.Values) error
}

func GetImage(ctx context.Context, url string) (Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer rsp.Body.Close()

	var buf bytes.Buffer
	r := io.TeeReader(rsp.Body, &buf)

	img, name, err := image.Decode(r)
	if name != "gif" {
		if err != nil {
			return nil, fmt.Errorf("decode image: %w", err)
		}
		return StillImage{img}, nil
	}

	g, err := gif.DecodeAll(io.MultiReader(&buf, rsp.Body))
	if err != nil {
		return nil, fmt.Errorf("decode GIF: %w", err)
	}
	return GIFImage{g}, nil
}

func GetPattern(ctx context.Context, q url.Values) (Image, error) {
	seed, _ := strconv.ParseUint(q.Get("seed"), 10, 0)

	patsrc := q.Get("pat")
	if patsrc == "" {
		if q.Get("sym") == "true" {
			return StillImage{sirdsc.SymmetricRandImage{Seed: seed}}, nil
		}
		return StillImage{sirdsc.RandImage{Seed: seed}}, nil
	}

	pat, err := GetImage(ctx, patsrc)
	if err != nil {
		return nil, fmt.Errorf("get image: %w", err)
	}
	return pat, nil
}

type StillImage struct {
	image.Image
}

func (img StillImage) Generate(ctx context.Context, w io.Writer, pat Image, q url.Values) error {
	partSize, _ := strconv.ParseInt(q.Get("partsize"), 10, 0)
	if partSize <= 0 {
		partSize = 100
	}

	max, _ := strconv.ParseInt(q.Get("depth"), 10, 0)
	if max <= 0 {
		max = 40
	}

	flat := q.Get("flat") == "true"
	inverse := q.Get("inverse") == "true"

	out := image.NewNRGBA(image.Rect(
		img.Bounds().Min.X,
		img.Bounds().Min.Y,
		img.Bounds().Max.X+int(partSize),
		img.Bounds().Max.Y,
	))

	sirdsc.Generate(
		out,
		sirdsc.ImageDepthMap{
			Image:   img,
			Max:     int(max),
			Flat:    flat,
			Inverse: inverse,
		},
		pat,
		int(partSize),
	)

	return png.Encode(w, out)
}

type GIFImage struct {
	*gif.GIF
}

func (img GIFImage) Generate(ctx context.Context, w io.Writer, pat Image, q url.Values) error {
	partSize, _ := strconv.ParseInt(q.Get("partsize"), 10, 0)
	if partSize <= 0 {
		partSize = 100
	}

	max, _ := strconv.ParseInt(q.Get("depth"), 10, 0)
	if max <= 0 {
		max = 40
	}

	flat := q.Get("flat") == "true"
	inverse := q.Get("inverse") == "true"

	eg, ctx := errgroup.WithContext(ctx)
	for i := range img.Image {
		eg.Go(func() error {
			out := image.NewPaletted(image.Rect(
				img.Image[i].Bounds().Min.X,
				img.Image[i].Bounds().Min.Y,
				img.Image[i].Bounds().Max.X+int(partSize),
				img.Image[i].Bounds().Max.Y,
			), palette.Plan9)

			sirdsc.Generate(
				out,
				sirdsc.ImageDepthMap{
					Image:   img.Image[i],
					Max:     int(max),
					Flat:    flat,
					Inverse: inverse,
				},
				pat,
				int(partSize),
			)

			img.Image[i] = out

			return nil
		})
	}

	img.Config.Width += int(partSize)

	err := eg.Wait()
	if err != nil {
		return err
	}

	return gif.EncodeAll(w, img.GIF)
}

func (img GIFImage) ColorModel() color.Model {
	return img.Image[0].ColorModel()
}

func (img GIFImage) Bounds() image.Rectangle {
	return img.Image[0].Bounds()
}

func (img GIFImage) At(x, y int) color.Color {
	return img.Image[0].At(x, y)
}
