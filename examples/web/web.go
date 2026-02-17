package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"golang.org/x/sync/errgroup"
)

//go:generate bun run --bun build

//go:embed dist
var distFS embed.FS

func configFromQuery(ctx context.Context, q url.Values) (*GenerateConfig, error) {
	seed, _ := strconv.ParseUint(q.Get("seed"), 10, 0)
	pat, err := GetPattern(ctx, seed, q.Get("sym") == "true", q.Get("pat"))
	if err != nil {
		return nil, fmt.Errorf("get pattern: %w", err)
	}

	partSize, _ := strconv.ParseInt(q.Get("partsize"), 10, 0)
	if partSize <= 0 {
		partSize = 100
	}

	maxDepth, _ := strconv.ParseInt(q.Get("depth"), 10, 0)
	if maxDepth <= 0 {
		maxDepth = 40
	}

	return &GenerateConfig{
		Pattern:  pat,
		PartSize: int(partSize),
		MaxDepth: int(maxDepth),
		Flat:     q.Get("flat") == "true",
		Inverse:  q.Get("inverse") == "true",
	}, nil
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
	slog := slog.With("src", src)

	imgC := make(chan Image, 1)
	configC := make(chan *GenerateConfig, 1)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		config, err := configFromQuery(ctx, q)
		if err != nil {
			return fmt.Errorf("parse query: %w", err)
		}

		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case configC <- config:
			return nil
		}
	})

	eg.Go(func() error {
		img, err := GetImage(ctx, src)
		if err != nil {
			return fmt.Errorf("get depth map: %w", err)
		}

		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case imgC <- img:
			return nil
		}
	})

	eg.Go(func() error {
		var img Image
		var config *GenerateConfig
		for img == nil || config == nil {
			select {
			case <-ctx.Done():
				return context.Cause(ctx)
			case img = <-imgC:
			case config = <-configC:
			}
		}

		err := img.Generate(ctx, rw, config)
		if err != nil {
			return fmt.Errorf("encode image: %w", err)
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
