// Command serve builds the site, serves dist/ locally, and rebuilds whenever
// content, templates or static files change. Refresh the browser to see edits.
package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"portfolio/internal/site"
)

func main() {
	var cfg site.Config
	cfg.RegisterFlags(flag.CommandLine)
	addr := flag.String("addr", "localhost:8080", "listen address")
	flag.Parse()

	rebuild := func() {
		start := time.Now()
		if err := site.Build(cfg); err != nil {
			log.Printf("build failed:\n%v", err)
			return
		}
		log.Printf("built in %s", time.Since(start).Round(time.Millisecond))
	}
	rebuild()
	go watch([]string{cfg.ContentDir, cfg.TemplatesDir, cfg.StaticDir}, rebuild)

	log.Printf("serving %s at http://%s", cfg.OutDir, *addr)
	log.Fatal(http.ListenAndServe(*addr, handler(cfg.OutDir)))
}

// handler serves dir mimicking GitHub Pages: directories are served via
// index.html, directories without an index.html get a 404, and missing paths
// get 404.html with a 404 status. Traversal and backslash paths are rejected
// by os.DirFS.
func handler(dir string) http.Handler {
	fsys := os.DirFS(dir)
	files := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")

		// Derive a valid fs.FS name from the URL path.
		// path.Clean rejects traversal; os.DirFS rejects any ".." component.
		name := strings.TrimPrefix(path.Clean("/"+r.URL.Path), "/")
		if name == "" {
			name = "."
		}

		info, err := fs.Stat(fsys, name)
		if err != nil {
			// Any stat error (not found, traversal attempt, etc.) → 404.
			serve404(w, dir)
			return
		}
		if info.IsDir() {
			// Only serve a directory if it has an index.html; no listings.
			if _, err := fs.Stat(fsys, path.Join(name, "index.html")); err != nil {
				serve404(w, dir)
				return
			}
		}
		files.ServeHTTP(w, r)
	})
}

func serve404(w http.ResponseWriter, dir string) {
	const fallback = "404 page not found (no 404.html; check the build log)"
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	page, err := os.ReadFile(filepath.Join(dir, "404.html"))
	if err != nil {
		w.Write([]byte(fallback))
		return
	}
	w.Write(page)
}

// watch polls dirs and calls onChange when any file is added, removed or
// modified. Polling keeps the tool dependency-free.
func watch(dirs []string, onChange func()) {
	last := snapshot(dirs)
	for range time.Tick(500 * time.Millisecond) {
		if cur := snapshot(dirs); cur != last {
			last = cur
			onChange()
		}
	}
}

func snapshot(dirs []string) uint64 {
	h := fnv.New64a()
	for _, d := range dirs {
		filepath.WalkDir(d, func(p string, e fs.DirEntry, err error) error {
			if err != nil {
				return nil
			}
			if info, err := e.Info(); err == nil {
				fmt.Fprintf(h, "%s|%d|%d\n", p, info.Size(), info.ModTime().UnixNano())
			}
			return nil
		})
	}
	return h.Sum64()
}
