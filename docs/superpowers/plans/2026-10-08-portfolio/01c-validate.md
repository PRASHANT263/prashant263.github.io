# Part 1c — Validation

### Task 4: Validate

**Files:** Create `internal/content/validate.go`, `internal/content/validate_test.go`

- [ ] **Step 1: Failing tests**

```go
package content

import (
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func mustDate(s string) Date {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		panic(err)
	}
	return Date{t}
}

func validSite() *Site {
	return &Site{
		Profile: Profile{
			Name: "Ada", Title: "Engineer", Pitch: "Builds.", Email: "ada@example.com",
			Photo: "img/ada.jpg", Resume: "resume.pdf", SiteURL: "https://ada.example.com",
			Socials: []Social{
				{Label: "GitHub", Icon: "github", URL: "https://github.com/ada"},
				{Label: "Email", Icon: "email", URL: "mailto:ada@example.com"},
			},
		},
		Experience: []Company{{Name: "Acme", Roles: []Role{{Title: "Engineer", Start: "Jan 2020"}}}},
		Skills:     []SkillGroup{{Name: "Languages", Items: []string{"Go"}}},
		Projects:   []Project{{Title: "Tool", Description: "A tool.", Repo: "https://github.com/ada/tool"}},
		Blog: []BlogEntry{
			{Title: "Post", Date: mustDate("2026-01-01"), URL: "https://medium.com/@ada/post"},
			{Title: "Local", Date: mustDate("2026-01-02"), Slug: "local"},
		},
		CaseStudies: []CaseStudy{{Slug: "migration", Title: "Migration", Summary: "Moved."}},
	}
}

var testStatic = fstest.MapFS{"img/ada.jpg": {}, "resume.pdf": {}}

func TestValidateValidSite(t *testing.T) {
	if err := Validate(validSite(), testStatic); err != nil {
		t.Fatal(err)
	}
}

func TestValidateProblems(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Site)
		want   string
	}{
		{"missing name", func(s *Site) { s.Profile.Name = "" }, "profile.yaml: name: is required"},
		{"relative site url", func(s *Site) { s.Profile.SiteURL = "ada.example.com" }, "profile.yaml: site_url:"},
		{"bad email", func(s *Site) { s.Profile.Email = "ada" }, "profile.yaml: email:"},
		{"missing photo file", func(s *Site) { s.Profile.Photo = "img/none.jpg" }, "static/img/none.jpg does not exist"},
		{"bad social url", func(s *Site) { s.Profile.Socials[0].URL = "github.com/ada" }, "profile.yaml: socials[0].url:"},
		{"no companies", func(s *Site) { s.Experience = nil }, "experience.yaml: companies: at least one"},
		{"role without title", func(s *Site) { s.Experience[0].Roles[0].Title = "" }, "experience.yaml: [0].roles[0].title: is required"},
		{"empty skill group", func(s *Site) { s.Skills[0].Items = nil }, "skills.yaml: [0].items:"},
		{"bad project repo", func(s *Site) { s.Projects[0].Repo = "nope" }, "projects.yaml: [0].repo:"},
		{"blog missing date", func(s *Site) { s.Blog[0].Date = Date{} }, "blog.yaml: [0].date: is required"},
		{"blog url and slug", func(s *Site) { s.Blog[0].Slug = "x" }, "blog.yaml: [0].url: set only one of url or slug"},
		{"blog neither", func(s *Site) { s.Blog[0].URL = "" }, "blog.yaml: [0].url: set either url"},
		{"blog bad slug", func(s *Site) { s.Blog[1].Slug = "Bad Slug" }, "blog.yaml: [1].slug:"},
		{"blog duplicate slug", func(s *Site) { s.Blog[0].URL = ""; s.Blog[0].Slug = "local" }, "blog.yaml: [1].slug: duplicate"},
		{"case study missing summary", func(s *Site) { s.CaseStudies[0].Summary = "" }, "case-studies/migration.md: summary: is required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validSite()
			tt.mutate(s)
			err := Validate(s, testStatic)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestValidateReportsEveryProblem(t *testing.T) {
	s := validSite()
	s.Profile.Name, s.Profile.Title = "", ""
	err := Validate(s, testStatic)
	for _, want := range []string{"name: is required", "title: is required"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v, missing %q", err, want)
		}
	}
}
```

- [ ] **Step 2: Run** `go test ./internal/content/` → FAIL (undefined: Validate)

- [ ] **Step 3: `validate.go`**

```go
package content

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Validate checks content for mistakes that would produce a broken site.
// static is the static/ directory, used to confirm referenced files exist.
// Every problem is reported, each as "file: field: message".
func Validate(s *Site, static fs.FS) error {
	var p problems
	p.profile(s.Profile, static)
	p.experience(s.Experience)
	p.skills(s.Skills)
	p.projects(s.Projects)
	p.blog(s.Blog)
	p.caseStudies(s.CaseStudies)
	return errors.Join(p...)
}

type problems []error

func (p *problems) add(file, field, msg string) {
	*p = append(*p, fmt.Errorf("%s: %s: %s", file, field, msg))
}

func (p *problems) required(file, field, value string) {
	if strings.TrimSpace(value) == "" {
		p.add(file, field, "is required")
	}
}

func (p *problems) absURL(file, field, value string) {
	if value != "" && !isAbsURL(value) {
		p.add(file, field, fmt.Sprintf("%q must be an absolute http(s) URL", value))
	}
}

func (p *problems) staticFile(static fs.FS, file, field, value string) {
	name := strings.TrimPrefix(value, "/")
	if name == "" {
		return
	}
	if _, err := fs.Stat(static, name); err != nil {
		p.add(file, field, fmt.Sprintf("static/%s does not exist", name))
	}
}

func isAbsURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func (p *problems) profile(pr Profile, static fs.FS) {
	const f = "profile.yaml"
	p.required(f, "name", pr.Name)
	p.required(f, "title", pr.Title)
	p.required(f, "pitch", pr.Pitch)
	p.required(f, "photo", pr.Photo)
	p.required(f, "resume", pr.Resume)
	p.required(f, "site_url", pr.SiteURL)
	p.absURL(f, "site_url", pr.SiteURL)
	if !strings.Contains(pr.Email, "@") {
		p.add(f, "email", fmt.Sprintf("%q is not an email address", pr.Email))
	}
	p.staticFile(static, f, "photo", pr.Photo)
	p.staticFile(static, f, "resume", pr.Resume)
	p.staticFile(static, f, "resume_preview", pr.ResumePreview)
	for i, s := range pr.Socials {
		field := fmt.Sprintf("socials[%d].url", i)
		p.required(f, field, s.URL)
		if !strings.HasPrefix(s.URL, "mailto:") {
			p.absURL(f, field, s.URL)
		}
	}
}

func (p *problems) experience(companies []Company) {
	const f = "experience.yaml"
	if len(companies) == 0 {
		p.add(f, "companies", "at least one company is required")
	}
	for i, c := range companies {
		p.required(f, fmt.Sprintf("[%d].name", i), c.Name)
		if len(c.Roles) == 0 {
			p.add(f, fmt.Sprintf("[%d].roles", i), "at least one role is required")
		}
		for j, r := range c.Roles {
			p.required(f, fmt.Sprintf("[%d].roles[%d].title", i, j), r.Title)
			p.required(f, fmt.Sprintf("[%d].roles[%d].start", i, j), r.Start)
		}
	}
}

func (p *problems) skills(groups []SkillGroup) {
	const f = "skills.yaml"
	for i, g := range groups {
		p.required(f, fmt.Sprintf("[%d].name", i), g.Name)
		if len(g.Items) == 0 {
			p.add(f, fmt.Sprintf("[%d].items", i), "at least one skill is required")
		}
	}
}

func (p *problems) projects(projects []Project) {
	const f = "projects.yaml"
	for i, pr := range projects {
		p.required(f, fmt.Sprintf("[%d].title", i), pr.Title)
		p.required(f, fmt.Sprintf("[%d].description", i), pr.Description)
		p.absURL(f, fmt.Sprintf("[%d].repo", i), pr.Repo)
		p.absURL(f, fmt.Sprintf("[%d].live", i), pr.Live)
	}
}

func (p *problems) blog(entries []BlogEntry) {
	const f = "blog.yaml"
	seen := map[string]bool{}
	for i, e := range entries {
		field := func(name string) string { return fmt.Sprintf("[%d].%s", i, name) }
		p.required(f, field("title"), e.Title)
		if e.Date.IsZero() {
			p.add(f, field("date"), "is required")
		}
		switch {
		case e.URL == "" && e.Slug == "":
			p.add(f, field("url"), "set either url (external link) or slug (local post)")
		case e.URL != "" && e.Slug != "":
			p.add(f, field("url"), "set only one of url or slug")
		case e.URL != "":
			p.absURL(f, field("url"), e.URL)
		default:
			if !slugPattern.MatchString(e.Slug) {
				p.add(f, field("slug"), fmt.Sprintf("%q must be lowercase letters, digits and hyphens", e.Slug))
			}
			if seen[e.Slug] {
				p.add(f, field("slug"), fmt.Sprintf("duplicate slug %q", e.Slug))
			}
			seen[e.Slug] = true
		}
	}
}

func (p *problems) caseStudies(studies []CaseStudy) {
	for _, cs := range studies {
		f := "case-studies/" + cs.Slug + ".md"
		if !slugPattern.MatchString(cs.Slug) {
			p.add(f, "file name", "must be lowercase letters, digits and hyphens")
		}
		p.required(f, "title", cs.Title)
		p.required(f, "summary", cs.Summary)
	}
}
```

- [ ] **Step 4: Run** `go test ./internal/content/` → PASS
- [ ] **Step 5: Commit** — `git add internal/content && git commit -m "feat(content): validate content with file/field errors"`
