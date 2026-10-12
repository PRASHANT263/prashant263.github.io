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
		{"recognition.yaml", &s.Recognition, false},
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
		if entries[i].Slug == "" || !slugPattern.MatchString(entries[i].Slug) {
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
