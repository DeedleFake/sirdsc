package main

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/DeedleFake/sirdsc"
	"golang.org/x/sync/errgroup"
)

func init() {
	http.DefaultClient.Timeout = 5 * time.Second
}

func getImage(url string) (image.Image, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	img, _, err := image.Decode(rsp.Body)
	return img, err
}

func handleMain(rw http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	q := req.URL.Query()

	src := q.Get("src")
	if src == "" {
		http.Error(rw, "No source specified.", http.StatusBadRequest)
		return
	}

	imgC := make(chan image.Image, 1)
	patC := make(chan image.Image, 1)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		seed, _ := strconv.ParseInt(q.Get("seed"), 10, 0)

		pat := image.Image(sirdsc.RandImage(seed))
		if q.Get("sym") == "true" {
			pat = sirdsc.SymmetricRandImage(seed)
		}
		if patsrc := q.Get("pat"); patsrc != "" {
			tmp, err := getImage(patsrc)
			if err != nil {
				return fmt.Errorf("Failed to get pattern from %q: %v", patsrc, err)
			}
			pat = tmp
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case patC <- pat:
			return nil
		}
	})

	eg.Go(func() error {
		img, err := getImage(src)
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
		partSize, _ := strconv.ParseInt(q.Get("partsize"), 10, 0)
		if partSize <= 0 {
			partSize = 100
		}

		max, _ := strconv.ParseInt(q.Get("depth"), 10, 0)
		if max <= 0 {
			max = 40
		}

		var img image.Image
		select {
		case <-ctx.Done():
			return ctx.Err()
		case img = <-imgC:
		}

		out := image.NewNRGBA(image.Rect(
			img.Bounds().Min.X,
			img.Bounds().Min.Y,
			img.Bounds().Max.X+int(partSize),
			img.Bounds().Max.Y,
		))

		var pat image.Image
		select {
		case <-ctx.Done():
			return ctx.Err()
		case pat = <-patC:
		}

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

		err := png.Encode(rw, out)
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
	http.HandleFunc("/", handleMain)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Failed to start server: %v", err)
	}
}
