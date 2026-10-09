# Part 3 — Site build and commands

### Task 8: Link checker

**Files:** Create `internal/site/links.go`, `internal/site/links_test.go`

- [ ] **Step 1: Failing test**

```go
package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTestFile(t *testing.T, path, data string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestCheckLinks(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "index.html"), `
<a href="/blog/">blog</a> <a href="/blog">blog no slash</a> <a href="/#about">anchor</a>
<link href="/css/a.css?v=1"> <img src="/missing.png">
<a href="https://example.com/x">ext</a> <a href="mailto:a@b.c">mail</a> <a href="#top">frag</a>
<a href="//cdn.example.com/x.js">protocol-relative</a>`)
	writeTestFile(t, filepath.Join(dir, "blog", "index.html"), `<a href="/">home</a>`)
	writeTestFile(t, filepath.Join(dir, "css", "a.css"), ``)

	err := CheckLinks(dir)
	if err == nil {
		t.Fatal("expected broken link error")
	}
	msg := err.Error()
	if !strings.Contains(msg, "index.html: broken link /missing.png") {
		t.Errorf("error = %q", msg)
	}
	if strings.Count(msg, "broken link") != 1 {
		t.Errorf("expected exactly one broken link, got %q", msg)
	}
}

func TestCheckLinksClean(t *testing.T) {
	dir := t.TempDir()
	writeTestFile(t, filepath.Join(dir, "index.html"), `<a href="/">home</a>`)
	if err := CheckLinks(dir); err != nil {
		t.Fatal(err)
	}
}
```

- [ ] **Step 2: Run** `go test ./internal/site/` → FAIL (undefined: CheckLinks)

- [ ] **Step 3: `links.go`**

```go
package site

import (
	"errors"
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var linkAttr = regexp.MustCompile(`(?:href|src)="([^"]*)"`)

// CheckLinks reports every root-relative href or src in dir's HTML files that
// doesn't resolve to a file in dir. External, mailto and fragment-only links
// are skipped.
func CheckLinks(dir string) error {
	var broken []error
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".html" {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		for _, m := range linkAttr.FindAllSubmatch(data, -1) {
			link := html.UnescapeString(string(m[1]))
			if !strings.HasPrefix(link, "/") || strings.HasPrefix(link, "//") {
				continue
			}
			if !resolves(dir, link) {
				broken = append(broken, fmt.Errorf("%s: broken link %s", rel, link))
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return errors.Join(broken...)
}

func resolves(dir, link string) bool {
	if i := strings.IndexAny(link, "?#"); i >= 0 {
		link = link[:i]
	}
	p := filepath.Join(dir, filepath.FromSlash(link))
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	if info.IsDir() {
		_, err = os.Stat(filepath.Join(p, "index.html"))
		return err == nil
	}
	return true
}
```

- [ ] **Step 4: Run** `go test ./internal/site/` → PASS
- [ ] **Step 5: Commit** — `git add internal/site && git commit -m "feat(site): add internal link checker"`

---

### Task 9: Build

**Files:** Create `internal/site/config.go`, `internal/site/build.go`, `internal/site/build_test.go`, `internal/site/testdata/content/{profile,experience,skills,blog}.yaml`, `internal/site/testdata/content/posts/hello.md`, `internal/site/testdata/content/case-studies/sample.md`

- [ ] **Step 1: Fixture content**

`testdata/content/profile.yaml`:
```yaml
name: Test Person
title: Engineer
pitch: Testing the build.
email: test@example.com
photo: img/profile.svg
resume: resume.pdf
site_url: https://test.example.com
about: ["About text."]
socials:
  - { label: Email, icon: email, url: "mailto:test@example.com" }
```
`testdata/content/experience.yaml`:
```yaml
- name: Acme
  roles:
    - { title: Engineer, start: Jan 2020, bullets: [Did work.] }
```
`testdata/content/skills.yaml`:
```yaml
- { name: Languages, items: [Go] }
```
`testdata/content/blog.yaml`:
```yaml
- { title: External, date: 2026-02-01, url: "https://medium.com/@t/x", source: Medium, tags: [go] }
- { title: Hello, date: 2026-01-01, slug: hello, tags: [intro] }
```
`testdata/content/posts/hello.md`:
```markdown
Hello from a local post.
```
`testdata/content/case-studies/sample.md`:
```markdown
---
title: Sample
summary: A sample case study.
order: 1
---
## Problem

Something.
```

- [ ] **Step 2: Failing test `build_test.go`**

```go
package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func testConfig(t *testing.T, contentDir string) Config {
	return Config{ContentDir: contentDir, TemplatesDir: "../../templates", StaticDir: "../../static", OutDir: t.TempDir()}
}

func TestBuild(t *testing.T) {
	cfg := testConfig(t, "testdata/content")
	if err := Build(cfg); err != nil {
		t.Fatal(err)
	}
	for _, f := range []string{
		"index.html", "blog/index.html", "blog/hello/index.html", "case-studies/sample/index.html",
		"404.html", "sitemap.xml", "robots.txt", "css/style.css", "js/main.js", "resume.pdf",
	} {
		if _, err := os.Stat(filepath.Join(cfg.OutDir, f)); err != nil {
			t.Errorf("missing %s", f)
		}
	}
	sitemap, _ := os.ReadFile(filepath.Join(cfg.OutDir, "sitemap.xml"))
	for _, want := range []string{"https://test.example.com/case-studies/sample/", "https://test.example.com/blog/hello/"} {
		if !strings.Contains(string(sitemap), want) {
			t.Errorf("sitemap missing %s", want)
		}
	}
}

func TestBuildFailureKeepsPreviousOutput(t *testing.T) {
	cfg := testConfig(t, t.TempDir()) // empty content dir: profile.yaml missing
	marker := filepath.Join(cfg.OutDir, "keep.txt")
	writeTestFile(t, marker, "old")
	err := Build(cfg)
	if err == nil || !strings.Contains(err.Error(), "profile.yaml") {
		t.Fatalf("err = %v", err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Error("output was cleared even though the build failed")
	}
}

func TestResetDirRefusesDangerousPaths(t *testing.T) {
	for _, d := range []string{"", ".", "/"} {
		if err := resetDir(d); err == nil {
			t.Errorf("resetDir(%q) should refuse", d)
		}
	}
}
```

- [ ] **Step 3: Run** `go test ./internal/site/` → FAIL (undefined: Config, Build, resetDir)

- [ ] **Step 4: `config.go`**

```go
package site

import "flag"

// Config says where Build reads from and writes to.
type Config struct {
	ContentDir   string
	TemplatesDir string
	StaticDir    string
	OutDir       string
}

// RegisterFlags adds -content, -templates, -static and -out flags defaulting
// to the repo layout.
func (c *Config) RegisterFlags(flags *flag.FlagSet) {
	flags.StringVar(&c.ContentDir, "content", "content", "content directory")
	flags.StringVar(&c.TemplatesDir, "templates", "templates", "templates directory")
	flags.StringVar(&c.StaticDir, "static", "static", "static files directory")
	flags.StringVar(&c.OutDir, "out", "dist", "output directory (replaced on every build)")
}
```

- [ ] **Step 5: `build.go`**

```go
// Package site builds the static site: it is the only package that touches
// the filesystem.
package site

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"portfolio/internal/content"
	"portfolio/internal/render"
)

// Build loads and validates content, then replaces OutDir with the rendered
// site plus static files, and checks internal links. Content and template
// errors are reported before OutDir is touched.
func Build(cfg Config) error {
	s, err := content.Load(os.DirFS(cfg.ContentDir))
	if err != nil {
		return err
	}
	if err := content.Validate(s, os.DirFS(cfg.StaticDir)); err != nil {
		return err
	}
	r, err := render.New(os.DirFS(cfg.TemplatesDir))
	if err != nil {
		return err
	}
	if err := resetDir(cfg.OutDir); err != nil {
		return err
	}
	if err := copyDir(cfg.StaticDir, cfg.OutDir); err != nil {
		return err
	}
	if err := writePages(r, s, cfg.OutDir); err != nil {
		return err
	}
	return CheckLinks(cfg.OutDir)
}

func writePages(r *render.Renderer, s *content.Site, out string) error {
	type page struct {
		file   string
		render func() ([]byte, error)
	}
	pages := []page{
		{"index.html", func() ([]byte, error) { return r.Home(s) }},
		{"blog/index.html", func() ([]byte, error) { return r.BlogIndex(s) }},
		{"404.html", func() ([]byte, error) { return r.NotFound(s) }},
	}
	urls := []string{"/", "/blog/"}
	for _, cs := range s.CaseStudies {
		pages = append(pages, page{"case-studies/" + cs.Slug + "/index.html", func() ([]byte, error) { return r.CaseStudy(s, cs) }})
		urls = append(urls, "/case-studies/"+cs.Slug+"/")
	}
	for _, e := range s.Blog {
		if e.IsExternal() {
			continue
		}
		pages = append(pages, page{"blog/" + e.Slug + "/index.html", func() ([]byte, error) { return r.Post(s, e) }})
		urls = append(urls, e.Href())
	}
	for _, p := range pages {
		html, err := p.render()
		if err != nil {
			return err
		}
		if err := writeFile(filepath.Join(out, p.file), html); err != nil {
			return err
		}
	}
	if err := writeFile(filepath.Join(out, "sitemap.xml"), render.Sitemap(s.Profile.SiteURL, urls)); err != nil {
		return err
	}
	return writeFile(filepath.Join(out, "robots.txt"), render.Robots(s.Profile.SiteURL))
}

func writeFile(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// resetDir empties dir, refusing paths where that would be catastrophic.
func resetDir(dir string) error {
	clean := filepath.Clean(dir)
	if clean == "." || clean == "/" {
		return fmt.Errorf("refusing to clear output directory %q", dir)
	}
	if err := os.RemoveAll(clean); err != nil {
		return err
	}
	return os.MkdirAll(clean, 0o755)
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if d.Name() == ".DS_Store" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}
```

- [ ] **Step 6: Run** `go test ./internal/site/` → PASS
- [ ] **Step 7: Commit** — `git add internal/site && git commit -m "feat(site): build dist from content, templates and static"`

---

### Task 10: `cmd/build` and `cmd/serve`

**Files:** Create `cmd/build/main.go`, `cmd/serve/main.go`, `cmd/serve/main_test.go`

- [ ] **Step 1: `cmd/build/main.go`**

```go
// Command build renders the site into dist/.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"portfolio/internal/site"
)

func main() {
	var cfg site.Config
	cfg.RegisterFlags(flag.CommandLine)
	flag.Parse()

	start := time.Now()
	if err := site.Build(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "build failed:\n%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("built %s in %s\n", cfg.OutDir, time.Since(start).Round(time.Millisecond))
}
```

- [ ] **Step 2: Failing test `cmd/serve/main_test.go`**

```go
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestHandlerServesFilesAnd404Page(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "index.html"), []byte("home"), 0o644)
	os.WriteFile(filepath.Join(dir, "404.html"), []byte("not found page"), 0o644)
	srv := httptest.NewServer(handler(dir))
	defer srv.Close()

	get := func(path string) (int, string) {
		res, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		defer res.Body.Close()
		b, _ := io.ReadAll(res.Body)
		return res.StatusCode, string(b)
	}
	if code, body := get("/"); code != 200 || body != "home" {
		t.Errorf("/ = %d %q", code, body)
	}
	if code, body := get("/missing/"); code != 404 || body != "not found page" {
		t.Errorf("/missing/ = %d %q", code, body)
	}
}

func TestSnapshotChangesWhenFileChanges(t *testing.T) {
	dir := t.TempDir()
	f := filepath.Join(dir, "a.txt")
	os.WriteFile(f, []byte("one"), 0o644)
	before := snapshot([]string{dir})
	os.WriteFile(f, []byte("three"), 0o644)
	if snapshot([]string{dir}) == before {
		t.Error("snapshot did not change after edit")
	}
}
```

- [ ] **Step 3: Run** `go test ./cmd/serve/` → FAIL (undefined: handler, snapshot)

- [ ] **Step 4: `cmd/serve/main.go`**

```go
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
```

- [ ] **Step 5: Run** `go test ./cmd/serve/` → PASS
- [ ] **Step 6: Commit** — `git add cmd && git commit -m "feat(cmd): add build and local dev server"`
