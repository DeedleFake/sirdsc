package main

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color/palette"
	"image/gif"
	_ "image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"

	"github.com/DeedleFake/sirdsc"
	"golang.org/x/sync/errgroup"
)

//go:generate bun run --bun build

//go:embed dist
var distFS embed.FS

func getImage(ctx context.Context, url string) (any, error) {
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
		return img, nil
	}

	g, err := gif.DecodeAll(io.MultiReader(&buf, rsp.Body))
	if err != nil {
		return nil, fmt.Errorf("decode GIF: %w", err)
	}
	return g, nil
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

func generateImage(w io.Writer, img image.Image, pat image.Image, q url.Values) error {
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

	imgC := make(chan any, 1)
	patC := make(chan image.Image, 1)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		seed, _ := strconv.ParseUint(q.Get("seed"), 10, 0)

		pat := image.Image(sirdsc.RandImage{Seed: seed})
		if q.Get("sym") == "true" {
			pat = sirdsc.SymmetricRandImage{Seed: seed}
		}
		if patsrc := q.Get("pat"); patsrc != "" {
			tmp, err := getImage(ctx, patsrc)
			if err != nil {
				return fmt.Errorf("Failed to get pattern from %q: %v", patsrc, err)
			}
			switch tmp := tmp.(type) {
			case image.Image:
				pat = tmp
			case *gif.GIF:
				return errors.New("pattern must not be a GIF")
			default:
				panic(reflect.TypeOf(tmp))
			}
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case patC <- pat:
			return nil
		}
	})

	eg.Go(func() error {
		img, err := getImage(ctx, src)
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
		var img any
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
			err = generateImage(rw, img, pat, q)
		default:
			panic(reflect.TypeOf(img))
		}
		if err != nil {
			return fmt.Errorf("Failed to encode image from %q: %v", src, err)
		}

		return nil
	})

	err := eg.Wait()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		slog.Error("generate failed", "err", err)
	}
}

func handleIndex(rw http.ResponseWriter, req *http.Request) {
	http.ServeFileFS(rw, req, distFS, "dist/index.html")
}

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		slog.Info("request", "remote", req.RemoteAddr, "method", req.Method, "url", req.URL)

		h.ServeHTTP(rw, req)
	})
}

func main() {
	addr := flag.String("addr", ":8080", "The address to listen on.")
	flag.Parse()

	http.Handle("GET /generate", logHandler(http.HandlerFunc(handleGenerate)))
	http.Handle("GET /dist/", logHandler(http.FileServerFS(distFS)))
	http.Handle("GET /", logHandler(http.HandlerFunc(handleIndex)))

	slog.Info("starting server", "addr", *addr)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		slog.Error("failed to start server", "err", err)
		os.Exit(1)
	}
}
