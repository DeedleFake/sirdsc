package main

import (
	"bytes"
	"context"
	"embed"
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/DeedleFake/sirdsc"
	"golang.org/x/sync/errgroup"
)

//go:embed interface/public/build
var fsys embed.FS

func init() {
	// TODO: Replace with context usage.
	http.DefaultClient.Timeout = 5 * time.Second
}

func getImage(url string, decodeGIFs bool) (interface{}, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	var buf bytes.Buffer
	r := io.TeeReader(rsp.Body, &buf)

	img, name, err := image.Decode(r)
	if !decodeGIFs || (name != "gif") {
		return img, err
	}

	return gif.DecodeAll(io.MultiReader(&buf, rsp.Body))
}

func generateGIF(ctx context.Context, w io.Writer, img *gif.GIF, pat image.Image, q url.Values) error {
	partSize, _ := strconv.ParseInt(q.Get("partsize"), 10, 0)
	if partSize <= 0 {
		partSize = 100
	}

	max, _ := strconv.ParseInt(q.Get("depth"), 10, 0)
	if max <= 0 {
		max = 40
	}

	eg, ctx := errgroup.WithContext(ctx)
	for i := range img.Image {
		i := i
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
					Flat:    q.Get("flat") == "true",
					Inverse: q.Get("inverse") == "true",
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

	return gif.EncodeAll(w, img)
}

func generateImage(ctx context.Context, w io.Writer, img image.Image, pat image.Image, q url.Values) error {
	partSize, _ := strconv.ParseInt(q.Get("partsize"), 10, 0)
	if partSize <= 0 {
		partSize = 100
	}

	max, _ := strconv.ParseInt(q.Get("depth"), 10, 0)
	if max <= 0 {
		max = 40
	}

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
			Flat:    q.Get("flat") == "true",
			Inverse: q.Get("inverse") == "true",
		},
		pat,
		int(partSize),
	)

	return png.Encode(w, out)
}

func handleGenerate(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(req.Context())
	defer cancel()

	q := req.URL.Query()

	src := q.Get("src")
	if src == "" {
		http.Error(rw, "No source specified.", http.StatusBadRequest)
		return
	}

	imgC := make(chan interface{}, 1)
	patC := make(chan image.Image, 1)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		seed, _ := strconv.ParseUint(q.Get("seed"), 10, 0)

		pat := image.Image(sirdsc.RandImage{Seed: seed})
		if q.Get("sym") == "true" {
			pat = sirdsc.SymmetricRandImage{Seed: seed}
		}
		if patsrc := q.Get("pat"); patsrc != "" {
			tmp, err := getImage(patsrc, false)
			if err != nil {
				return fmt.Errorf("Failed to get pattern from %q: %v", patsrc, err)
			}
			pat = tmp.(image.Image)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case patC <- pat:
			return nil
		}
	})

	eg.Go(func() error {
		img, err := getImage(src, true)
		if err != nil {
			return fmt.Errorf("Failed to get depth map from %q: %v", src, err)
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case imgC <- img:
			return nil
		}
	})

	eg.Go(func() error {
		var img interface{}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case img = <-imgC:
		}

		var pat image.Image
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pat = <-patC:
		}

		var err error
		switch img := img.(type) {
		case *gif.GIF:
			err = generateGIF(ctx, rw, img, pat, q)
		case image.Image:
			err = generateImage(ctx, rw, img, pat, q)
		}
		if err != nil {
			return fmt.Errorf("Failed to encode image from %q: %v", src, err)
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func main() {
	addr := flag.String("addr", ":8080", "The address to listen on.")
	flag.Parse()

	sub, err := fs.Sub(fsys, "interface/public")
	if err != nil {
		log.Fatalf("Failed to create sub FS: %v", err)
	}

	http.HandleFunc("/generate", handleGenerate)
	http.Handle("/", http.FileServer(http.FS(sub)))

	log.Printf("Starting server on %q", *addr)
	err = http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
