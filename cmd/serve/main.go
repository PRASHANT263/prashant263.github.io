// Command serve builds the site, serves dist/ locally, and rebuilds whenever
// content, templates or static files change. Refresh the browser to see edits.
package main

import (
	"errors"
	"flag"
	"fmt"
	"hash/fnv"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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

// handler serves dir the way GitHub Pages does: directories serve index.html
// and missing paths get 404.html with a 404 status.
func handler(dir string) http.Handler {
	files := http.FileServer(http.Dir(dir))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		p := filepath.Join(dir, filepath.FromSlash(path.Clean("/"+r.URL.Path)))
		if _, err := os.Stat(p); errors.Is(err, fs.ErrNotExist) {
			page, _ := os.ReadFile(filepath.Join(dir, "404.html"))
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			w.Write(page)
			return
		}
		files.ServeHTTP(w, r)
	})
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
