package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

//go:generate bun run --bun build

//go:embed dist
var distFS embed.FS

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
	patC := make(chan Image, 1)

	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		pat, err := GetPattern(ctx, q)
		if err != nil {
			return fmt.Errorf("get pattern: %w", err)
		}

		select {
		case <-ctx.Done():
			return context.Cause(ctx)
		case patC <- pat:
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
		var img, pat Image
		for img == nil || pat == nil {
			select {
			case <-ctx.Done():
				return context.Cause(ctx)
			case img = <-imgC:
			case pat = <-patC:
			}
		}

		err := img.Generate(ctx, rw, pat, q)
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
