package render

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"portfolio/internal/content"
)

var update = flag.Bool("update", false, "rewrite golden files in testdata/golden")

func fixtureSite() *content.Site {
	d := func(s string) content.Date {
		t, _ := time.Parse("2006-01-02", s)
		return content.Date{Time: t}
	}
	return &content.Site{
		Profile: content.Profile{
			Name: "Ada", Title: "Backend Engineer", Pitch: "Builds reliable systems.",
			About: []content.AboutPoint{{Title: "What I do", Icon: "code", Text: "Hello."}}, Location: "Earth", Email: "ada@example.com",
			Photo: "img/ada.jpg", Resume: "resume.pdf", ResumePreview: "img/resume.png",
			ResumeUpdated: d("2026-01-02"), Web3FormsKey: "test-key", SiteURL: "https://ada.example.com",
			Status: "Currently building at Acme", Expertise: []string{"Go", "DDD"},
			Facts: []content.Fact{{Value: "5", Label: "Years"}},
			Socials: []content.Social{
				{Label: "GitHub", Icon: "github", URL: "https://github.com/ada"},
				{Label: "Email", Icon: "email", URL: "mailto:ada@example.com"},
			},
		},
		Experience: []content.Company{{Name: "Acme", Location: "Remote", Roles: []content.Role{
			{Title: "Engineer", Start: "Jan 2020", Highlights: []string{"+10% speed"}, Bullets: []string{"Built things."}},
		}}},
		Skills:   []content.SkillGroup{{Name: "Languages", Items: []string{"Go", "SQL"}}},
		Projects: []content.Project{{Title: "Tool", Description: "A tool.", Tech: []string{"Go"}, Repo: "https://github.com/ada/tool"}},
		Blog: []content.BlogEntry{
			{Title: "External", Date: d("2026-03-01"), URL: "https://medium.com/@ada/x", Source: "Medium", Tags: []string{"go"}, Summary: "An article."},
			{Title: "Local", Date: d("2026-02-01"), Slug: "local", Tags: []string{"ddd"}, Body: "Some *text*."},
		},
		Recognition: []content.Recognition{{Title: "Star Performer", When: "2024", Detail: "Twice.", Kind: "award"}, {Title: "B.Tech", When: "2018 – 2022", Kind: "education"}},
		CaseStudies: []content.CaseStudy{{Slug: "migration", Title: "Migration", Summary: "Moved to DDD.", Metrics: []string{"faster"}, Tags: []string{"ddd"}, Body: "## Problem\n\nIt was slow.\n"}},
	}
}

func realRenderer(t *testing.T) *Renderer {
	t.Helper()
	r, err := New(os.DirFS("../../templates"))
	if err != nil {
		t.Fatal(err)
	}
	r.now = func() time.Time { return time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }
	return r
}

func renderAll(t *testing.T) map[string][]byte {
	t.Helper()
	r, s := realRenderer(t), fixtureSite()
	out := map[string][]byte{}
	for name, fn := range map[string]func() ([]byte, error){
		"home.html":       func() ([]byte, error) { return r.Home(s) },
		"blog.html":       func() ([]byte, error) { return r.BlogIndex(s) },
		"case-study.html": func() ([]byte, error) { return r.CaseStudy(s, s.CaseStudies[0]) },
		"post.html":       func() ([]byte, error) { return r.Post(s, s.Blog[1]) },
		"404.html":        func() ([]byte, error) { return r.NotFound(s) },
	} {
		b, err := fn()
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		out[name] = b
	}
	return out
}

func TestRealTemplatesContent(t *testing.T) {
	pages := renderAll(t)
	mustContain(t, string(pages["home.html"]),
		`<title>Ada · Backend Engineer</title>`,
		`<link rel="canonical" href="https://ada.example.com/">`,
		`id="about"`, `id="experience"`, `id="work"`, `id="skills"`, `id="resume"`, `id="blog"`, `id="projects"`, `id="contact"`,
		`src="/img/ada.jpg"`,
		`href="/resume.pdf" download`,
		`src="/img/resume.png"`,
		`Last updated Jan 2, 2026`,
		`href="/case-studies/migration/"`,
		`href="https://medium.com/@ada/x" target="_blank" rel="noopener"`,
		`href="/blog/local/"`,
		`name="access_key" value="test-key"`,
		`© 2026 Ada`,
	)
	mustContain(t, string(pages["blog.html"]), `<title>Writing · Ada</title>`, `data-tag="ddd"`, `data-blog-list`)
	mustContain(t, string(pages["case-study.html"]), `<title>Migration · Ada</title>`, `<h2 id="problem">Problem</h2>`, `<li>faster</li>`)
	mustContain(t, string(pages["post.html"]), `<h1>Local</h1>`, `<em>text</em>`)
	mustContain(t, string(pages["404.html"]), `href="/"`)
}

func TestContactFormNoKey(t *testing.T) {
	r := realRenderer(t)
	s := fixtureSite()
	s.Profile.Web3FormsKey = ""
	out, err := r.Home(s)
	if err != nil {
		t.Fatal(err)
	}
	body := string(out)
	if strings.Contains(body, "api.web3forms.com") {
		t.Error("home with no Web3FormsKey must not render the web3forms form")
	}
	mustContain(t, body, "The contact form is coming soon.")
}

func TestHeroGreetings(t *testing.T) {
	r := realRenderer(t)
	s := fixtureSite()

	out, err := r.Home(s)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), `<p class="eyebrow mono">Hi, I'm</p>`)

	s.Profile.Greetings = []string{"Hi", "नमस्ते", "ನಮಸ್ಕಾರ"}
	out, err = r.Home(s)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out),
		`<span class="sr-only">Hi, I'm</span>`,
		`data-greetings="Hi|नमस्ते|ನಮಸ್ಕಾರ"`,
		`<span class="greet" data-greetings="Hi|नमस्ते|ನಮಸ್ಕಾರ">Hi</span><span class="cursor"></span>, I'm`,
	)
}

func TestRealTemplatesGolden(t *testing.T) {
	for name, got := range renderAll(t) {
		path := filepath.Join("testdata", "golden", name)
		if *update {
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				t.Fatal(err)
			}
			if err := os.WriteFile(path, got, 0o644); err != nil {
				t.Fatal(err)
			}
			continue
		}
		want, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("%v (run: go test ./internal/render -update)", err)
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s differs from golden file; review the diff, then run: go test ./internal/render -update", name)
		}
	}
}

func TestHomeSectionsFromSashwatIdeas(t *testing.T) {
	out, err := realRenderer(t).Home(fixtureSite())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out),
		// status pill above the name
		`<p class="status-pill"><span class="status-dot" aria-hidden="true"></span>Currently building at Acme</p>`,
		// eyebrow labels with short headings
		`<p class="section-eyebrow mono">Career</p>`,
		`<h2 class="section-title">Where I've built.</h2>`,
		// recognition section
		`id="recognition"`, `<h3>Star Performer</h3>`, `<p class="when mono">2018 – 2022</p>`,
		// closing call-to-action band
		`class="cta-band`, `href="mailto:ada@example.com"`,
		// richer footer
		`<h4>Quick links</h4>`, `<h4>Expertise</h4>`, `<li>DDD</li>`, `<h4>Resources</h4>`,
	)
}

func TestAboutPillars(t *testing.T) {
	out, err := realRenderer(t).Home(fixtureSite())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out),
		`<article class="glass card about-card">`,
		`<h3>What I do</h3>`,
		`<p>Hello.</p>`,
		`<polyline points="16 18 22 12 16 6"/>`, // the "code" icon
	)
}

func TestResumeDownloadNotRepeated(t *testing.T) {
	out, err := realRenderer(t).Home(fixtureSite())
	if err != nil {
		t.Fatal(err)
	}
	// Only the Resume section downloads; the hero's "Download CV" scrolls there.
	if n := strings.Count(string(out), `href="/resume.pdf" download`); n != 1 {
		t.Errorf("resume download links = %d, want 1", n)
	}
	mustContain(t, string(out), `<a class="btn" href="#resume">Download CV`, `href="/resume.pdf" download>Download CV`)
	if strings.Contains(string(out), "Resume (PDF)") {
		t.Error("footer should not repeat the resume link")
	}
}

func TestFactsRenderAsHeroBadge(t *testing.T) {
	s := fixtureSite()
	s.Profile.Facts = []content.Fact{{Value: "4+ years", Label: "of backend experience"}}
	out, err := realRenderer(t).Home(s)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), `<span class="hero-badge"><strong>4&#43; years</strong> of backend experience</span>`)
	if strings.Contains(string(out), `<ul class="facts">`) {
		t.Error("stat cards should no longer render")
	}
}

func TestHeroPhotoSitsAfterText(t *testing.T) {
	out, err := realRenderer(t).Home(fixtureSite())
	if err != nil {
		t.Fatal(err)
	}
	body := string(out)
	text, photo := strings.Index(body, `<div class="hero-text">`), strings.Index(body, `<div class="hero-photo">`)
	if text < 0 || photo < 0 || photo < text {
		t.Errorf("hero photo should come after the hero text (right column); text=%d photo=%d", text, photo)
	}
	mustContain(t, body, `<img src="/img/ada.jpg" alt="Photo of Ada" width="400" height="500"`)
}

func TestHeroAvatarCircleOnLeft(t *testing.T) {
	s := fixtureSite()
	s.Profile.Avatar = "img/avatar.svg"
	out, err := realRenderer(t).Home(s)
	if err != nil {
		t.Fatal(err)
	}
	body := string(out)
	avatar, text, photo := strings.Index(body, `<div class="hero-avatar">`), strings.Index(body, `<div class="hero-text">`), strings.Index(body, `<div class="hero-photo">`)
	if avatar < 0 || !(avatar < text && text < photo) {
		t.Errorf("want avatar, text, photo in that order; got %d, %d, %d", avatar, text, photo)
	}
	mustContain(t, body, `<img src="/img/avatar.svg" alt="" width="220" height="220">`)
}

func TestEmailAddressShownOnce(t *testing.T) {
	out, err := realRenderer(t).Home(fixtureSite())
	if err != nil {
		t.Fatal(err)
	}
	if n := strings.Count(string(out), ">ada@example.com<") + strings.Count(string(out), "</svg> ada@example.com<"); n != 1 {
		t.Errorf("email address shown as text %d times, want 1 (Contact only)", n)
	}
}

func TestFooterLocationHasPinIcon(t *testing.T) {
	out, err := realRenderer(t).Home(fixtureSite())
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out), `<span class="location">`, `<circle cx="12" cy="10" r="3"/></svg> Earth</span>`)
}

func TestMediumSocialHasItsOwnIcon(t *testing.T) {
	s := fixtureSite()
	s.Profile.Socials = append(s.Profile.Socials, content.Social{Label: "Medium", Icon: "medium", URL: "https://medium.com/@ada"})
	out, err := realRenderer(t).Home(s)
	if err != nil {
		t.Fatal(err)
	}
	mustContain(t, string(out),
		`<a href="https://medium.com/@ada" aria-label="Medium" target="_blank" rel="noopener"><svg class="icon medium"`,
		`<li><a href="https://medium.com/@ada" target="_blank" rel="noopener">Medium</a></li>`,
	)
}
