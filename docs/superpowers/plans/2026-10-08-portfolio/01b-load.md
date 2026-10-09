# Part 1b — Content types and loading

### Task 3: Types + Load

**Files:** Create `internal/content/types.go`, `frontmatter.go`, `load.go`, `load_test.go`

- [ ] **Step 1: `types.go`** (types only, no behaviour beyond two URL helpers)

```go
package content

// Site is everything the generator knows about, loaded from content/.
type Site struct {
	Profile     Profile
	Experience  []Company
	Skills      []SkillGroup
	Projects    []Project
	Blog        []BlogEntry // newest first
	CaseStudies []CaseStudy // by Order, then Slug
}

type Profile struct {
	Name          string   `yaml:"name"`
	Title         string   `yaml:"title"`
	Pitch         string   `yaml:"pitch"`
	About         []string `yaml:"about"`
	Location      string   `yaml:"location"`
	Email         string   `yaml:"email"`
	Photo         string   `yaml:"photo"`          // path under static/
	Resume        string   `yaml:"resume"`         // path under static/
	ResumePreview string   `yaml:"resume_preview"` // optional, path under static/
	ResumeUpdated Date     `yaml:"resume_updated"`
	Web3FormsKey  string   `yaml:"web3forms_key"`
	SiteURL       string   `yaml:"site_url"`
	Facts         []Fact   `yaml:"facts"`
	Socials       []Social `yaml:"socials"`
}

type Fact struct {
	Value string `yaml:"value"`
	Label string `yaml:"label"`
}

type Social struct {
	Label string `yaml:"label"`
	Icon  string `yaml:"icon"` // github | linkedin | email
	URL   string `yaml:"url"`
}

type Company struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	Roles    []Role `yaml:"roles"`
}

type Role struct {
	Title      string   `yaml:"title"`
	Start      string   `yaml:"start"`
	End        string   `yaml:"end"` // empty means Present
	Highlights []string `yaml:"highlights"`
	Bullets    []string `yaml:"bullets"`
}

type SkillGroup struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items"`
}

type Project struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Tech        []string `yaml:"tech"`
	Repo        string   `yaml:"repo"`
	Live        string   `yaml:"live"`
}

// BlogEntry is either an external link (URL set) or a local post (Slug set,
// body read from content/posts/<slug>.md).
type BlogEntry struct {
	Title   string   `yaml:"title"`
	Date    Date     `yaml:"date"`
	URL     string   `yaml:"url"`
	Slug    string   `yaml:"slug"`
	Source  string   `yaml:"source"`
	Tags    []string `yaml:"tags"`
	Summary string   `yaml:"summary"`
	Body    string   `yaml:"-"` // Markdown, local posts only
}

func (e BlogEntry) IsExternal() bool { return e.URL != "" }

// Href is where the entry's link points.
func (e BlogEntry) Href() string {
	if e.IsExternal() {
		return e.URL
	}
	return "/blog/" + e.Slug + "/"
}

type CaseStudy struct {
	Slug    string   `yaml:"-"` // from the file name
	Title   string   `yaml:"title"`
	Summary string   `yaml:"summary"`
	Order   int      `yaml:"order"`
	Tags    []string `yaml:"tags"`
	Metrics []string `yaml:"metrics"`
	Body    string   `yaml:"-"` // Markdown
}
```

- [ ] **Step 2: Failing tests `load_test.go`**

```go
package content

import (
	"strings"
	"testing"
	"testing/fstest"
)

const profileYAML = `
name: Ada
title: Backend Engineer
pitch: Builds reliable systems.
email: ada@example.com
photo: img/ada.jpg
resume: resume.pdf
site_url: https://ada.example.com
`

const experienceYAML = `
- name: Acme
  roles:
    - title: Engineer
      start: Jan 2020
      bullets: [Built things.]
`

const skillsYAML = `
- name: Languages
  items: [Go, SQL]
`

func minimalFS() fstest.MapFS {
	return fstest.MapFS{
		"profile.yaml":    {Data: []byte(profileYAML)},
		"experience.yaml": {Data: []byte(experienceYAML)},
		"skills.yaml":     {Data: []byte(skillsYAML)},
	}
}

func TestLoadMinimal(t *testing.T) {
	s, err := Load(minimalFS())
	if err != nil {
		t.Fatal(err)
	}
	if s.Profile.Name != "Ada" {
		t.Errorf("name = %q", s.Profile.Name)
	}
	if len(s.Experience) != 1 || s.Experience[0].Roles[0].Title != "Engineer" {
		t.Errorf("experience = %+v", s.Experience)
	}
	if len(s.Skills) != 1 || len(s.Skills[0].Items) != 2 {
		t.Errorf("skills = %+v", s.Skills)
	}
	if len(s.Blog)+len(s.Projects)+len(s.CaseStudies) != 0 {
		t.Errorf("optional content should be empty")
	}
}

func TestLoadCommentOnlyFileIsEmpty(t *testing.T) {
	fsys := minimalFS()
	fsys["blog.yaml"] = &fstest.MapFile{Data: []byte("# nothing yet\n")}
	s, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Blog) != 0 {
		t.Errorf("blog = %+v", s.Blog)
	}
}

func TestLoadBlogSortedNewestFirst(t *testing.T) {
	fsys := minimalFS()
	fsys["blog.yaml"] = &fstest.MapFile{Data: []byte(`
- {title: Old, date: 2025-01-01, url: "https://a.example/old"}
- {title: New, date: 2026-05-01, url: "https://a.example/new"}
`)}
	s, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if s.Blog[0].Title != "New" || s.Blog[1].Title != "Old" {
		t.Errorf("order = %q, %q", s.Blog[0].Title, s.Blog[1].Title)
	}
}

func TestLoadLocalPost(t *testing.T) {
	fsys := minimalFS()
	fsys["blog.yaml"] = &fstest.MapFile{Data: []byte("- {title: Mine, date: 2026-01-01, slug: mine}\n")}
	fsys["posts/mine.md"] = &fstest.MapFile{Data: []byte("# Hello")}
	s, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	e := s.Blog[0]
	if e.Body != "# Hello" || e.IsExternal() || e.Href() != "/blog/mine/" {
		t.Errorf("entry = %+v href=%q", e, e.Href())
	}
}

func TestLoadCaseStudiesSortedByOrder(t *testing.T) {
	fsys := minimalFS()
	fsys["case-studies/a.md"] = &fstest.MapFile{Data: []byte("---\ntitle: A\nsummary: S\norder: 2\n---\nBody A\n")}
	fsys["case-studies/b.md"] = &fstest.MapFile{Data: []byte("---\ntitle: B\nsummary: S\norder: 1\nmetrics: [fast]\n---\n## Problem\n")}
	s, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if len(s.CaseStudies) != 2 || s.CaseStudies[0].Slug != "b" || s.CaseStudies[1].Slug != "a" {
		t.Fatalf("case studies = %+v", s.CaseStudies)
	}
	b := s.CaseStudies[0]
	if b.Title != "B" || b.Metrics[0] != "fast" || b.Body != "## Problem\n" {
		t.Errorf("b = %+v", b)
	}
}

func TestLoadErrors(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(fstest.MapFS)
		want   []string
	}{
		{"missing profile", func(f fstest.MapFS) { delete(f, "profile.yaml") }, []string{"profile.yaml"}},
		{"unknown field", func(f fstest.MapFS) {
			f["profile.yaml"] = &fstest.MapFile{Data: []byte(profileYAML + "nmae: x\n")}
		}, []string{"profile.yaml", "nmae"}},
		{"bad date", func(f fstest.MapFS) {
			f["blog.yaml"] = &fstest.MapFile{Data: []byte("- {title: X, date: 01/02/2026, url: \"https://x.example\"}\n")}
		}, []string{"blog.yaml", "YYYY-MM-DD"}},
		{"missing post file", func(f fstest.MapFS) {
			f["blog.yaml"] = &fstest.MapFile{Data: []byte("- {title: X, date: 2026-01-01, slug: gone}\n")}
		}, []string{"blog.yaml", "[0].slug"}},
		{"case study without front matter", func(f fstest.MapFS) {
			f["case-studies/x.md"] = &fstest.MapFile{Data: []byte("# just markdown\n")}
		}, []string{"case-studies/x.md", "front matter"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := minimalFS()
			tt.mutate(fsys)
			_, err := Load(fsys)
			if err == nil {
				t.Fatal("expected error")
			}
			for _, w := range tt.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("error %q missing %q", err, w)
				}
			}
		})
	}
}

func TestSplitFrontMatter(t *testing.T) {
	tests := []struct {
		name, in, front, body string
		wantErr               bool
	}{
		{"basic", "---\na: 1\n---\nbody\n", "a: 1\n", "body\n", false},
		{"crlf", "---\r\na: 1\r\n---\r\nbody\r\n", "a: 1\n", "body\n", false},
		{"empty front", "---\n---\nbody", "", "body", false},
		{"no opening", "body", "", "", true},
		{"unclosed", "---\na: 1\nbody", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			front, body, err := splitFrontMatter([]byte(tt.in))
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v", err)
			}
			if string(front) != tt.front || body != tt.body {
				t.Errorf("front=%q body=%q", front, body)
			}
		})
	}
}
```

- [ ] **Step 3: Run** `go test ./internal/content/` → FAIL (undefined: Load, splitFrontMatter)

- [ ] **Step 4: `frontmatter.go`**

```go
package content

import (
	"errors"
	"strings"
)

// splitFrontMatter separates a leading "---" YAML block from a Markdown body.
func splitFrontMatter(data []byte) (front []byte, body string, err error) {
	s := strings.ReplaceAll(string(data), "\r\n", "\n")
	if !strings.HasPrefix(s, "---\n") {
		return nil, "", errors.New("must start with a --- front matter block")
	}
	s = s[len("---\n"):]
	if strings.HasPrefix(s, "---\n") {
		return nil, s[len("---\n"):], nil
	}
	end := strings.Index(s, "\n---\n")
	if end < 0 {
		return nil, "", errors.New("front matter is missing its closing ---")
	}
	return []byte(s[:end+1]), s[end+len("\n---\n"):], nil
}
```

- [ ] **Step 5: `load.go`**

```go
package content

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Load reads every content file from fsys (the content/ directory).
// profile.yaml, experience.yaml and skills.yaml are required; projects.yaml,
// blog.yaml, posts and case studies are optional.
func Load(fsys fs.FS) (*Site, error) {
	s := &Site{}
	files := []struct {
		name     string
		out      any
		required bool
	}{
		{"profile.yaml", &s.Profile, true},
		{"experience.yaml", &s.Experience, true},
		{"skills.yaml", &s.Skills, true},
		{"projects.yaml", &s.Projects, false},
		{"blog.yaml", &s.Blog, false},
	}
	for _, f := range files {
		if err := readYAML(fsys, f.name, f.out, f.required); err != nil {
			return nil, err
		}
	}
	if err := loadPosts(fsys, s.Blog); err != nil {
		return nil, err
	}
	cs, err := loadCaseStudies(fsys)
	if err != nil {
		return nil, err
	}
	s.CaseStudies = cs
	sort.SliceStable(s.Blog, func(i, j int) bool { return s.Blog[i].Date.After(s.Blog[j].Date.Time) })
	return s, nil
}

func readYAML(fsys fs.FS, name string, out any, required bool) error {
	data, err := fs.ReadFile(fsys, name)
	if errors.Is(err, fs.ErrNotExist) && !required {
		return nil
	}
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	if err := decodeStrict(data, out); err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return nil
}

// decodeStrict rejects unknown fields so typos fail the build. An empty or
// comment-only document leaves out untouched.
func decodeStrict(data []byte, out any) error {
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(out); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func loadPosts(fsys fs.FS, entries []BlogEntry) error {
	for i := range entries {
		if entries[i].Slug == "" {
			continue
		}
		data, err := fs.ReadFile(fsys, path.Join("posts", entries[i].Slug+".md"))
		if err != nil {
			return fmt.Errorf("blog.yaml: [%d].slug: %w", i, err)
		}
		entries[i].Body = string(data)
	}
	return nil
}

func loadCaseStudies(fsys fs.FS) ([]CaseStudy, error) {
	names, err := fs.Glob(fsys, "case-studies/*.md")
	if err != nil {
		return nil, err
	}
	var out []CaseStudy
	for _, name := range names {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		front, body, err := splitFrontMatter(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		var cs CaseStudy
		if err := decodeStrict(front, &cs); err != nil {
			return nil, fmt.Errorf("%s: %w", name, err)
		}
		cs.Slug = strings.TrimSuffix(path.Base(name), ".md")
		cs.Body = body
		out = append(out, cs)
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Order != out[j].Order {
			return out[i].Order < out[j].Order
		}
		return out[i].Slug < out[j].Slug
	})
	return out, nil
}
```

- [ ] **Step 6: Run** `go test ./internal/content/` → PASS
- [ ] **Step 7: Commit** — `git add internal/content && git commit -m "feat(content): load YAML and Markdown content"`
