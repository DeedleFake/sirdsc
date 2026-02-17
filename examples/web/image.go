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
	"sync"
	"time"

	"github.com/DeedleFake/sirdsc"
	"golang.org/x/sync/errgroup"
)

type Image interface {
	image.Image
	Generate(ctx context.Context, w io.Writer, config *GenerateConfig) error
}

type GenerateConfig struct {
	Pattern  image.Image
	PartSize int
	MaxDepth int
	Flat     bool
	Inverse  bool
}

var cache sync.Map

type cacheEntry struct {
	url    string
	img    Image
	cancel func()
}

func loadCachedImage(url string) (Image, bool) {
	v, ok := cache.Load(url)
	if !ok {
		return nil, false
	}

	entry := v.(*cacheEntry)
	entry.reset()
	return entry.img, true
}

func storeImage(url string, img Image) Image {
	entry := cacheEntry{url: url, img: img}
	old, ok := entry.reset()
	if !ok {
		return img
	}
	old.cancel()
	return img
}

func (entry cacheEntry) reset() (*cacheEntry, bool) {
	if entry.cancel != nil {
		entry.cancel()
	}

	t := time.AfterFunc(10*time.Second, func() {
		cache.CompareAndDelete(entry.url, &entry)
	})
	entry.cancel = func() { t.Stop() }

	old, ok := cache.Swap(entry.url, &entry)
	if !ok {
		return nil, false
	}
	return old.(*cacheEntry), true
}

func GetImage(ctx context.Context, url string) (Image, error) {
	old, ok := loadCachedImage(url)
	if ok {
		return old, nil
	}

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
		return storeImage(url, StillImage{img}), nil
	}

	g, err := gif.DecodeAll(io.MultiReader(&buf, rsp.Body))
	if err != nil {
		return nil, fmt.Errorf("decode GIF: %w", err)
	}
	return storeImage(url, GIFImage{g}), nil
}

func GetPattern(ctx context.Context, seed uint64, sym bool, patsrc string) (Image, error) {
	if patsrc == "" {
		if sym {
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

func (img StillImage) Generate(ctx context.Context, w io.Writer, config *GenerateConfig) error {
	out := image.NewNRGBA(image.Rect(
		img.Bounds().Min.X,
		img.Bounds().Min.Y,
		img.Bounds().Max.X+config.PartSize,
		img.Bounds().Max.Y,
	))

	sirdsc.Generate(
		out,
		sirdsc.ImageDepthMap{
			Image:   img,
			Max:     config.MaxDepth,
			Flat:    config.Flat,
			Inverse: config.Inverse,
		},
		config.Pattern,
		config.PartSize,
	)

	return png.Encode(w, out)
}

type GIFImage struct {
	*gif.GIF
}

func (img GIFImage) copy() *gif.GIF {
	newGIF := *img.GIF
	newGIF.Image = make([]*image.Paletted, len(newGIF.Image))
	return &newGIF
}

func (img GIFImage) Generate(ctx context.Context, w io.Writer, config *GenerateConfig) error {
	newGIF := img.copy()
	newGIF.Config.Width += config.PartSize

	eg, ctx := errgroup.WithContext(ctx)
	for i := range img.Image {
		eg.Go(func() error {
			out := image.NewPaletted(image.Rect(
				img.Image[i].Bounds().Min.X,
				img.Image[i].Bounds().Min.Y,
				img.Image[i].Bounds().Max.X+config.PartSize,
				img.Image[i].Bounds().Max.Y,
			), palette.Plan9)

			sirdsc.Generate(
				out,
				sirdsc.ImageDepthMap{
					Image:   img.Image[i],
					Max:     config.MaxDepth,
					Flat:    config.Flat,
					Inverse: config.Inverse,
				},
				config.Pattern,
				config.PartSize,
			)

			newGIF.Image[i] = out

			return nil
		})
	}

	err := eg.Wait()
	if err != nil {
		return err
	}

	return gif.EncodeAll(w, newGIF)
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
