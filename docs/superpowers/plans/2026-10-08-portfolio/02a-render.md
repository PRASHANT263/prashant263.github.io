# Part 2a — Renderer core

### Task 5: Renderer, funcs, sitemap

**Files:** Create `internal/render/{render.go,funcs.go,sitemap.go,render_test.go}`, `internal/render/testdata/templates/{base.html,partials/footer.html,pages/home.html,pages/blog.html,pages/case-study.html,pages/post.html,pages/404.html}`

- [ ] **Step 1: Fixture templates** (tiny, exist only to test Renderer mechanics)

`testdata/templates/base.html`:
```
{{define "base"}}<title>{{.Title}}</title>{{template "main" .}}{{template "footer" .}}{{end}}
```
`testdata/templates/partials/footer.html`:
```
{{define "footer"}}<footer>{{.Year}}</footer>{{end}}
```
`testdata/templates/pages/home.html`:
```
{{define "main"}}<h1>{{.Site.Profile.Name}}</h1>{{range take 1 .Site.Blog}}<a href="{{.Href}}">{{.Title}}</a>{{end}}{{end}}
```
`testdata/templates/pages/blog.html`:
```
{{define "main"}}{{range .Tags}}<b>{{.}}</b>{{end}}{{end}}
```
`testdata/templates/pages/case-study.html`:
```
{{define "main"}}{{.CaseStudy.Title}}{{.Body}}{{end}}
```
`testdata/templates/pages/post.html`:
```
{{define "main"}}{{.Post.Title}}{{.Body}}{{end}}
```
`testdata/templates/pages/404.html`:
```
{{define "main"}}missing{{end}}
```

- [ ] **Step 2: Failing tests `render_test.go`**

```go
package render

import (
	"os"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"portfolio/internal/content"
)

func testSite() *content.Site {
	return &content.Site{
		Profile: content.Profile{Name: "Ada", SiteURL: "https://ada.example.com/"},
		Blog: []content.BlogEntry{
			{Title: "Newest", URL: "https://x.example/new", Tags: []string{"go", "ddd"}},
			{Title: "Older", Slug: "older", Tags: []string{"go"}, Body: "Hello **world**"},
		},
		CaseStudies: []content.CaseStudy{
			{Slug: "migration", Title: "Migration", Body: "## Problem\n\n<svg class=\"diagram\"></svg>\n"},
		},
	}
}

func newTestRenderer(t *testing.T) *Renderer {
	t.Helper()
	r, err := New(os.DirFS("testdata/templates"))
	if err != nil {
		t.Fatal(err)
	}
	r.now = func() time.Time { return time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC) }
	return r
}

func mustContain(t *testing.T, got string, wants ...string) {
	t.Helper()
	for _, w := range wants {
		if !strings.Contains(got, w) {
			t.Errorf("output missing %q\n---\n%s", w, got)
		}
	}
}

func TestHome(t *testing.T) {
	out, err := newTestRenderer(t).Home(testSite())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "<h1>Ada</h1>", `<a href="https://x.example/new">Newest</a>`, "<footer>2030</footer>")
	if strings.Contains(string(out), "Older") {
		t.Error("take 1 should limit entries")
	}
}

func TestBlogIndexTagsSortedUnique(t *testing.T) {
	out, err := newTestRenderer(t).BlogIndex(testSite())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "<b>ddd</b><b>go</b>")
}

func TestCaseStudyRendersMarkdownAndRawHTML(t *testing.T) {
	s := testSite()
	out, err := newTestRenderer(t).CaseStudy(s, s.CaseStudies[0])
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "<title>Migration</title>", `<h2 id="problem">Problem</h2>`, `<svg class="diagram"></svg>`)
}

func TestPost(t *testing.T) {
	s := testSite()
	out, err := newTestRenderer(t).Post(s, s.Blog[1])
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "Older", "<strong>world</strong>")
}

func TestNotFound(t *testing.T) {
	out, err := newTestRenderer(t).NotFound(testSite())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "missing")
}

func TestPageCanonical(t *testing.T) {
	p := Page{Site: testSite(), Path: "/blog/"}
	if got := p.Canonical(); got != "https://ada.example.com/blog/" {
		t.Errorf("Canonical() = %q", got)
	}
}

func TestNewFailsWithoutTemplates(t *testing.T) {
	if _, err := New(fstest.MapFS{}); err == nil {
		t.Fatal("expected error for empty templates fs")
	}
}

func TestSitemapAndRobots(t *testing.T) {
	sm := string(Sitemap("https://ada.example.com/", []string{"/", "/blog/"}))
	mustContain(t, sm, "<loc>https://ada.example.com/</loc>", "<loc>https://ada.example.com/blog/</loc>")
	mustContain(t, string(Robots("https://ada.example.com/")), "Sitemap: https://ada.example.com/sitemap.xml")
}
```

- [ ] **Step 3: Run** `go test ./internal/render/` → FAIL (undefined: New, Page, ...)

- [ ] **Step 4: `render.go`**

```go
// Package render turns content into HTML pages. It reads templates from an
// fs.FS and returns bytes; it never writes to disk.
package render

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"portfolio/internal/content"
)

// pageNames are the files in templates/pages/. Each is parsed together with
// base.html and every partial, so pages can't clash on block names.
var pageNames = []string{"home", "blog", "case-study", "post", "404"}

type Renderer struct {
	pages map[string]*template.Template
	md    goldmark.Markdown
	now   func() time.Time
}

// Page is the data every template receives.
type Page struct {
	Site        *content.Site
	Path        string // site-relative URL path, e.g. "/blog/"
	Title       string // empty means "<name> · <title>"
	Description string
	Body        template.HTML // rendered Markdown for case studies and posts
	CaseStudy   *content.CaseStudy
	Post        *content.BlogEntry
	Tags        []string // every blog tag, for the blog filter
	Year        int
}

// SiteRoot is the site URL without a trailing slash.
func (p Page) SiteRoot() string { return strings.TrimSuffix(p.Site.Profile.SiteURL, "/") }

func (p Page) Canonical() string { return p.SiteRoot() + p.Path }

func New(templates fs.FS) (*Renderer, error) {
	r := &Renderer{
		pages: map[string]*template.Template{},
		md: goldmark.New(
			goldmark.WithExtensions(extension.GFM),
			goldmark.WithParserOptions(parser.WithAutoHeadingID()),
			// Case studies embed inline SVG diagrams; content is owner-written.
			goldmark.WithRendererOptions(html.WithUnsafe()),
		),
		now: time.Now,
	}
	for _, name := range pageNames {
		t, err := template.New(name).Funcs(funcs).ParseFS(templates, "base.html", "partials/*.html", "pages/"+name+".html")
		if err != nil {
			return nil, fmt.Errorf("templates: %w", err)
		}
		r.pages[name] = t
	}
	return r, nil
}

func (r *Renderer) Home(s *content.Site) ([]byte, error) {
	return r.execute("home", Page{Site: s, Path: "/", Description: s.Profile.Pitch})
}

func (r *Renderer) BlogIndex(s *content.Site) ([]byte, error) {
	return r.execute("blog", Page{
		Site: s, Path: "/blog/", Title: "Writing",
		Description: "Articles by " + s.Profile.Name, Tags: blogTags(s.Blog),
	})
}

func (r *Renderer) CaseStudy(s *content.Site, cs content.CaseStudy) ([]byte, error) {
	body, err := r.markdown(cs.Body)
	if err != nil {
		return nil, fmt.Errorf("case study %s: %w", cs.Slug, err)
	}
	return r.execute("case-study", Page{
		Site: s, Path: "/case-studies/" + cs.Slug + "/", Title: cs.Title,
		Description: cs.Summary, Body: body, CaseStudy: &cs,
	})
}

func (r *Renderer) Post(s *content.Site, e content.BlogEntry) ([]byte, error) {
	body, err := r.markdown(e.Body)
	if err != nil {
		return nil, fmt.Errorf("post %s: %w", e.Slug, err)
	}
	return r.execute("post", Page{
		Site: s, Path: e.Href(), Title: e.Title, Description: e.Summary, Body: body, Post: &e,
	})
}

func (r *Renderer) NotFound(s *content.Site) ([]byte, error) {
	return r.execute("404", Page{Site: s, Path: "/404.html", Title: "Page not found"})
}

func (r *Renderer) execute(name string, p Page) ([]byte, error) {
	p.Year = r.now().Year()
	var buf bytes.Buffer
	if err := r.pages[name].ExecuteTemplate(&buf, "base", p); err != nil {
		return nil, fmt.Errorf("render %s: %w", name, err)
	}
	return buf.Bytes(), nil
}

func (r *Renderer) markdown(src string) (template.HTML, error) {
	var buf bytes.Buffer
	if err := r.md.Convert([]byte(src), &buf); err != nil {
		return "", err
	}
	return template.HTML(buf.String()), nil
}
```

- [ ] **Step 5: `funcs.go`**

```go
package render

import (
	"html/template"
	"sort"
	"strings"

	"portfolio/internal/content"
)

var funcs = template.FuncMap{
	"take": func(n int, entries []content.BlogEntry) []content.BlogEntry {
		if len(entries) > n {
			return entries[:n]
		}
		return entries
	},
	"join":      strings.Join,
	"hasPrefix": strings.HasPrefix,
	// asset turns a static/-relative path from content into a root URL.
	"asset": func(p string) string { return "/" + strings.TrimPrefix(p, "/") },
}

// blogTags returns every tag used by entries, sorted and de-duplicated.
func blogTags(entries []content.BlogEntry) []string {
	seen := map[string]bool{}
	var tags []string
	for _, e := range entries {
		for _, t := range e.Tags {
			if !seen[t] {
				seen[t] = true
				tags = append(tags, t)
			}
		}
	}
	sort.Strings(tags)
	return tags
}
```

- [ ] **Step 6: `sitemap.go`**

```go
package render

import (
	"encoding/xml"
	"strings"
)

// Sitemap returns sitemap.xml content listing paths under siteURL.
func Sitemap(siteURL string, paths []string) []byte {
	base := strings.TrimSuffix(siteURL, "/")
	var b strings.Builder
	b.WriteString(xml.Header)
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, p := range paths {
		b.WriteString("  <url><loc>")
		_ = xml.EscapeText(&b, []byte(base+p))
		b.WriteString("</loc></url>\n")
	}
	b.WriteString("</urlset>\n")
	return []byte(b.String())
}

func Robots(siteURL string) []byte {
	return []byte("User-agent: *\nAllow: /\nSitemap: " + strings.TrimSuffix(siteURL, "/") + "/sitemap.xml\n")
}
```

- [ ] **Step 7: Run** `go test ./internal/render/` → PASS
- [ ] **Step 8: Commit** — `git add internal/render && git commit -m "feat(render): add template renderer, markdown and sitemap"`
