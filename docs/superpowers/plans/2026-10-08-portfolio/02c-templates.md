# Part 2c — Real templates

### Task 7: Templates + golden tests

**Files:** Create `templates/base.html`, `templates/partials/{icons,nav,footer,socials,blog-card}.html`, `templates/pages/{home,blog,case-study,post,404}.html`, `internal/render/golden_test.go`, `internal/render/testdata/golden/*.html` (generated)

- [ ] **Step 1: Failing test `internal/render/golden_test.go`**

```go
package render

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
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
			About: []string{"Hello."}, Location: "Earth", Email: "ada@example.com",
			Photo: "img/ada.jpg", Resume: "resume.pdf", ResumePreview: "img/resume.png",
			ResumeUpdated: d("2026-01-02"), Web3FormsKey: "test-key", SiteURL: "https://ada.example.com",
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
		`data-configured="true"`,
		`© 2026 Ada`,
	)
	mustContain(t, string(pages["blog.html"]), `<title>Writing · Ada</title>`, `data-tag="ddd"`, `data-blog-list`)
	mustContain(t, string(pages["case-study.html"]), `<title>Migration · Ada</title>`, `<h2 id="problem">Problem</h2>`, `<li>faster</li>`)
	mustContain(t, string(pages["post.html"]), `<h1>Local</h1>`, `<em>text</em>`)
	mustContain(t, string(pages["404.html"]), `href="/"`)
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
```

- [ ] **Step 2: Run** `go test ./internal/render/ -run RealTemplates` → FAIL (templates: pattern matches no files)

- [ ] **Step 3: `templates/base.html`**

```html
{{define "base"}}<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>{{if .Title}}{{.Title}} · {{.Site.Profile.Name}}{{else}}{{.Site.Profile.Name}} · {{.Site.Profile.Title}}{{end}}</title>
<meta name="description" content="{{.Description}}">
<link rel="canonical" href="{{.Canonical}}">
<meta property="og:type" content="website">
<meta property="og:title" content="{{if .Title}}{{.Title}} · {{.Site.Profile.Name}}{{else}}{{.Site.Profile.Name}} · {{.Site.Profile.Title}}{{end}}">
<meta property="og:description" content="{{.Description}}">
<meta property="og:url" content="{{.Canonical}}">
<meta property="og:image" content="{{.SiteRoot}}/img/og.png">
<meta name="twitter:card" content="summary_large_image">
<meta name="theme-color" content="#0b1020">
<link rel="icon" href="/img/favicon.svg" type="image/svg+xml">
<link rel="preload" href="/fonts/inter-var.woff2" as="font" type="font/woff2" crossorigin>
<link rel="stylesheet" href="/css/style.css">
<script>document.documentElement.classList.add("js");try{var t=localStorage.getItem("theme");if(t)document.documentElement.dataset.theme=t}catch(e){}</script>
<script src="/js/main.js" defer></script>
</head>
<body>
<div class="mesh" aria-hidden="true"><span class="blob b1"></span><span class="blob b2"></span><span class="blob b3"></span><span class="blob b4"></span></div>
{{template "nav" .}}
<main id="main">
{{template "main" .}}
</main>
{{template "footer" .}}
</body>
</html>
{{end}}
```

- [ ] **Step 4: `templates/partials/icons.html`** (Feather icons, MIT)

```html
{{define "icon-open"}}<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">{{end}}
{{define "icon-github"}}{{template "icon-open"}}<path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/></svg>{{end}}
{{define "icon-linkedin"}}{{template "icon-open"}}<path d="M16 8a6 6 0 0 1 6 6v7h-4v-7a2 2 0 0 0-4 0v7h-4v-7a6 6 0 0 1 6-6z"/><rect x="2" y="9" width="4" height="12"/><circle cx="4" cy="4" r="2"/></svg>{{end}}
{{define "icon-mail"}}{{template "icon-open"}}<path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/><polyline points="22,6 12,13 2,6"/></svg>{{end}}
{{define "icon-download"}}{{template "icon-open"}}<path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"/><polyline points="7 10 12 15 17 10"/><line x1="12" y1="15" x2="12" y2="3"/></svg>{{end}}
{{define "icon-external"}}{{template "icon-open"}}<path d="M18 13v6a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2V8a2 2 0 0 1 2-2h6"/><polyline points="15 3 21 3 21 9"/><line x1="10" y1="14" x2="21" y2="3"/></svg>{{end}}
{{define "icon-arrow"}}{{template "icon-open"}}<line x1="5" y1="12" x2="19" y2="12"/><polyline points="12 5 19 12 12 19"/></svg>{{end}}
{{define "icon-menu"}}{{template "icon-open"}}<line x1="3" y1="12" x2="21" y2="12"/><line x1="3" y1="6" x2="21" y2="6"/><line x1="3" y1="18" x2="21" y2="18"/></svg>{{end}}
{{define "icon-sun"}}<svg class="icon sun" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" aria-hidden="true"><circle cx="12" cy="12" r="5"/><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"/></svg>{{end}}
{{define "icon-moon"}}<svg class="icon moon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true"><path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"/></svg>{{end}}
{{define "social-icon"}}{{if eq . "github"}}{{template "icon-github"}}{{else if eq . "linkedin"}}{{template "icon-linkedin"}}{{else}}{{template "icon-mail"}}{{end}}{{end}}
```

Note: `icon-open` emits an unclosed `<svg>`; html/template may reject an unbalanced partial. If `New` errors with a context/branch error, replace each `{{template "icon-open"}}` with the literal `<svg class="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">` and delete the `icon-open` define.

- [ ] **Step 5: `templates/partials/socials.html`** (dot is `[]content.Social`)

```html
{{define "socials"}}<ul class="socials">{{range .}}<li><a href="{{.URL}}" aria-label="{{.Label}}"{{if not (hasPrefix .URL "mailto:")}} target="_blank" rel="noopener"{{end}}>{{template "social-icon" .Icon}}</a></li>{{end}}</ul>{{end}}
```

- [ ] **Step 6: `templates/partials/nav.html`**

```html
{{define "nav"}}<header class="nav glass">
  <a class="brand" href="/">{{.Site.Profile.Name}}<span class="dot">.</span></a>
  <nav id="nav-links" class="nav-links" aria-label="Main">
    <a href="/#about">About</a>
    <a href="/#experience">Experience</a>
    <a href="/#work">Work</a>
    <a href="/#resume">Resume</a>
    <a href="/blog/">Blog</a>
    <a href="/#contact">Contact</a>
  </nav>
  <button class="icon-btn theme-toggle" type="button" aria-label="Toggle color theme">{{template "icon-sun"}}{{template "icon-moon"}}</button>
  <button class="icon-btn nav-toggle" type="button" aria-label="Open menu" aria-expanded="false" aria-controls="nav-links">{{template "icon-menu"}}</button>
</header>{{end}}
```

- [ ] **Step 7: `templates/partials/footer.html`**

```html
{{define "footer"}}<footer class="footer">
  <p>© {{.Year}} {{.Site.Profile.Name}} · Built with Go, served as static HTML.</p>
  {{template "socials" .Site.Profile.Socials}}
</footer>{{end}}
```

- [ ] **Step 8: `templates/partials/blog-card.html`** (dot is `content.BlogEntry`)

```html
{{define "blog-card"}}<article class="card glass blog-card" data-tags="{{join .Tags " "}}">
  <div class="meta"><span class="badge">{{if .IsExternal}}{{or .Source "Article"}}{{else}}Post{{end}}</span><time datetime="{{.Date.ISO}}">{{.Date.Human}}</time></div>
  <h3><a href="{{.Href}}"{{if .IsExternal}} target="_blank" rel="noopener"{{end}}>{{.Title}}{{if .IsExternal}} {{template "icon-external"}}{{end}}</a></h3>
  {{with .Summary}}<p>{{.}}</p>{{end}}
  {{with .Tags}}<ul class="tags">{{range .}}<li>#{{.}}</li>{{end}}</ul>{{end}}
</article>{{end}}
```

- [ ] **Step 9: `templates/pages/home.html`**

```html
{{define "main"}}
{{with .Site.Profile}}
<section class="hero" id="top">
  <div class="hero-photo"><img src="{{asset .Photo}}" alt="Photo of {{.Name}}" width="220" height="220" fetchpriority="high"></div>
  <div class="hero-text">
    <p class="eyebrow mono">Hi, I'm</p>
    <h1>{{.Name}}</h1>
    <p class="role">{{.Title}}</p>
    <p class="pitch">{{.Pitch}}</p>
    <div class="actions">
      <a class="btn primary" href="#work">View work {{template "icon-arrow"}}</a>
      <a class="btn" href="{{asset .Resume}}" download>Download resume {{template "icon-download"}}</a>
    </div>
    {{template "socials" .Socials}}
  </div>
</section>

<section id="about" class="section reveal">
  <h2 class="section-title">About</h2>
  <div class="about glass card">
    <div class="about-text">{{range .About}}<p>{{.}}</p>{{end}}</div>
    {{with .Facts}}<ul class="facts">{{range .}}<li><strong>{{.Value}}</strong><span>{{.Label}}</span></li>{{end}}</ul>{{end}}
  </div>
</section>
{{end}}

<section id="experience" class="section reveal">
  <h2 class="section-title">Experience</h2>
  {{range .Site.Experience}}
  <div class="company">
    <h3 class="company-name">{{.Name}}{{with .Location}} <span class="muted">· {{.}}</span>{{end}}</h3>
    <ol class="timeline">
      {{range .Roles}}
      <li class="glass card role-card">
        <div class="role-head"><h4>{{.Title}}</h4><span class="mono muted">{{.Start}} — {{or .End "Present"}}</span></div>
        {{with .Highlights}}<ul class="chips">{{range .}}<li>{{.}}</li>{{end}}</ul>{{end}}
        {{with .Bullets}}<ul class="bullets">{{range .}}<li>{{.}}</li>{{end}}</ul>{{end}}
      </li>
      {{end}}
    </ol>
  </div>
  {{end}}
</section>

<section id="work" class="section reveal">
  <h2 class="section-title">Case studies</h2>
  <div class="grid">
    {{range .Site.CaseStudies}}
    <a class="card glass case-card" href="/case-studies/{{.Slug}}/">
      <h3>{{.Title}}</h3>
      <p>{{.Summary}}</p>
      {{with .Metrics}}<ul class="chips">{{range .}}<li>{{.}}</li>{{end}}</ul>{{end}}
      <span class="more">Read case study {{template "icon-arrow"}}</span>
    </a>
    {{end}}
  </div>
</section>

<section id="skills" class="section reveal">
  <h2 class="section-title">Skills</h2>
  <div class="grid skills">
    {{range .Site.Skills}}
    <div class="glass card"><h3>{{.Name}}</h3><ul class="tags">{{range .Items}}<li>{{.}}</li>{{end}}</ul></div>
    {{end}}
  </div>
</section>

{{with .Site.Profile}}
<section id="resume" class="section reveal">
  <h2 class="section-title">Resume</h2>
  <div class="glass card resume">
    {{with .ResumePreview}}<img class="resume-preview" src="{{asset .}}" alt="Preview of the first page of the resume" width="220" height="285" loading="lazy">{{end}}
    <div>
      <h3>{{.Name}} — Resume</h3>
      <p class="muted">A one-page summary of my experience, skills and achievements.</p>
      {{if not .ResumeUpdated.IsZero}}<p class="mono muted">Last updated {{.ResumeUpdated.Human}}</p>{{end}}
      <div class="actions">
        <a class="btn primary" href="{{asset .Resume}}" download>Download PDF {{template "icon-download"}}</a>
        <a class="btn" href="{{asset .Resume}}" target="_blank" rel="noopener">View online {{template "icon-external"}}</a>
      </div>
    </div>
  </div>
</section>
{{end}}

<section id="blog" class="section reveal">
  <div class="section-head"><h2 class="section-title">Writing</h2>{{if .Site.Blog}}<a class="link" href="/blog/">See all {{template "icon-arrow"}}</a>{{end}}</div>
  {{if .Site.Blog}}<div class="grid">{{range take 6 .Site.Blog}}{{template "blog-card" .}}{{end}}</div>
  {{else}}<p class="glass card empty">Articles are on the way. Check back soon.</p>{{end}}
</section>

{{with .Site.Projects}}
<section id="projects" class="section reveal">
  <h2 class="section-title">Projects</h2>
  <div class="grid">
    {{range .}}
    <article class="card glass project-card">
      <h3>{{.Title}}</h3>
      <p>{{.Description}}</p>
      {{with .Tech}}<ul class="tags">{{range .}}<li>{{.}}</li>{{end}}</ul>{{end}}
      {{if or .Repo .Live}}<p class="links">{{with .Repo}}<a href="{{.}}" target="_blank" rel="noopener">{{template "icon-github"}} Code</a>{{end}}{{with .Live}}<a href="{{.}}" target="_blank" rel="noopener">{{template "icon-external"}} Live</a>{{end}}</p>{{end}}
    </article>
    {{end}}
  </div>
</section>
{{end}}

{{with .Site.Profile}}
<section id="contact" class="section reveal">
  <h2 class="section-title">Contact</h2>
  <div class="glass card contact">
    <div>
      <h3>Let's talk</h3>
      <p class="muted">Open to senior backend roles and good engineering conversations. Send a message or email me directly.</p>
      <p><a class="link" href="mailto:{{.Email}}">{{template "icon-mail"}} {{.Email}}</a></p>
    </div>
    <form class="contact-form" action="https://api.web3forms.com/submit" method="POST" data-contact data-configured="{{if .Web3FormsKey}}true{{end}}">
      <input type="hidden" name="access_key" value="{{.Web3FormsKey}}">
      <input type="hidden" name="subject" value="New message from your portfolio">
      <input type="hidden" name="redirect" value="{{$.SiteRoot}}/#contact">
      <input type="checkbox" name="botcheck" class="hp" tabindex="-1" autocomplete="off" aria-hidden="true">
      <label>Name<input name="name" required maxlength="100" autocomplete="name"></label>
      <label>Email<input name="email" type="email" required maxlength="200" autocomplete="email"></label>
      <label>Message<textarea name="message" required minlength="10" maxlength="5000" rows="5"></textarea></label>
      <button class="btn primary" type="submit">Send message</button>
      <p class="form-status" role="status" aria-live="polite"></p>
    </form>
  </div>
</section>
{{end}}
{{end}}
```

- [ ] **Step 10: `templates/pages/blog.html`**

```html
{{define "main"}}
<section class="section page-head">
  <h1>Writing</h1>
  <p class="lead">Articles on Go, backend architecture and building fintech systems.</p>
  {{with .Tags}}<div class="tag-filter" role="group" aria-label="Filter by tag">
    <button class="tag active" type="button" data-tag="" aria-pressed="true">All</button>
    {{range .}}<button class="tag" type="button" data-tag="{{.}}" aria-pressed="false">#{{.}}</button>{{end}}
  </div>{{end}}
</section>
<section class="section">
  {{if .Site.Blog}}<div class="grid" data-blog-list>{{range .Site.Blog}}{{template "blog-card" .}}{{end}}</div>
  {{else}}<p class="glass card empty">Articles are on the way. Check back soon.</p>{{end}}
</section>
{{end}}
```

- [ ] **Step 11: `templates/pages/case-study.html`**

```html
{{define "main"}}
<article class="prose-page">
  <a class="link back" href="/#work">← All case studies</a>
  <header class="page-head">
    <h1>{{.CaseStudy.Title}}</h1>
    <p class="lead">{{.CaseStudy.Summary}}</p>
    {{with .CaseStudy.Metrics}}<ul class="chips">{{range .}}<li>{{.}}</li>{{end}}</ul>{{end}}
    {{with .CaseStudy.Tags}}<ul class="tags">{{range .}}<li>#{{.}}</li>{{end}}</ul>{{end}}
  </header>
  <div class="section"><div class="glass card prose">{{.Body}}</div></div>
</article>
{{end}}
```

- [ ] **Step 12: `templates/pages/post.html`**

```html
{{define "main"}}
<article class="prose-page">
  <a class="link back" href="/blog/">← All writing</a>
  <header class="page-head">
    <h1>{{.Post.Title}}</h1>
    <p class="mono muted"><time datetime="{{.Post.Date.ISO}}">{{.Post.Date.Human}}</time></p>
    {{with .Post.Tags}}<ul class="tags">{{range .}}<li>#{{.}}</li>{{end}}</ul>{{end}}
  </header>
  <div class="section"><div class="glass card prose">{{.Body}}</div></div>
</article>
{{end}}
```

- [ ] **Step 13: `templates/pages/404.html`**

```html
{{define "main"}}
<section class="glass card notfound">
  <p class="mono muted">404</p>
  <h1>This page wandered off.</h1>
  <p class="muted">The link may be broken, or the page may have moved.</p>
  <a class="btn primary" href="/">Back home {{template "icon-arrow"}}</a>
</section>
{{end}}
```

- [ ] **Step 14: Run** `go test ./internal/render/ -run RealTemplatesContent -v` → PASS. Fix template text until every `mustContain` string matches.
- [ ] **Step 15: Generate golden files** `go test ./internal/render/ -run RealTemplatesGolden -update`, then `go test ./internal/render/` → PASS. Skim `testdata/golden/home.html` for obvious escaping problems (e.g. `ZgotmplZ`, `&lt;svg`).
- [ ] **Step 16: Commit** — `git add templates internal/render && git commit -m "feat(templates): add glass portfolio templates with golden tests"`
