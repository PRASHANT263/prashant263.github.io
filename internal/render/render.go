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
