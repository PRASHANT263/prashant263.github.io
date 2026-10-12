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
		{"eof no newline", "---\na: 1\n---", "a: 1\n", "", false},
		{"empty front eof", "---\n---", "", "", false},
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

func TestLoadBadSlugSkipsFileAndValidateReports(t *testing.T) {
	fsys := minimalFS()
	fsys["blog.yaml"] = &fstest.MapFile{Data: []byte(`
- {title: Bad, date: 2026-01-01, slug: "Bad Slug"}
- {title: Traversal, date: 2026-01-02, slug: "../profile"}
`)}
	s, err := Load(fsys)
	if err != nil {
		t.Fatalf("Load should succeed for bad slug (no file read attempted): %v", err)
	}
	err = Validate(s, testStatic)
	if err == nil {
		t.Fatal("Validate should report slug errors")
	}
	for _, want := range []string{"[0].slug:", "must be lowercase"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error %q missing %q", err, want)
		}
	}
}
