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

func TestDescriptionFallbackToPitch(t *testing.T) {
	s := testSite()
	s.Profile.Pitch = "Expert engineer."
	// Post with no Summary: Description should fall back to Pitch.
	entry := content.BlogEntry{Slug: "test", Title: "Test", Body: "# Hello"}
	out, err := newTestRenderer(t).Post(s, entry)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), "<desc>Expert engineer.</desc>")
}

func TestNotFoundNoIndex(t *testing.T) {
	out, err := newTestRenderer(t).NotFound(testSite())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), `<meta name="robots" content="noindex">`)
}
