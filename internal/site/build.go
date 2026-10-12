// Package site builds the static site: it is the only package that writes
// the built site.
package site

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"portfolio/internal/content"
	"portfolio/internal/render"
)

// Build loads and validates content, renders every page, then replaces OutDir
// with the static files and rendered pages, and checks internal links.
// Content, template-parse and template-execution errors are all reported
// before OutDir is touched.
func Build(cfg Config) error {
	if err := checkOutDir(cfg); err != nil {
		return err
	}
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
	// Render all pages BEFORE touching OutDir so that template execution
	// errors leave the previous output untouched.
	rendered, err := renderAll(r, s)
	if err != nil {
		return err
	}
	if err := resetDir(cfg.OutDir); err != nil {
		return err
	}
	if err := copyDir(cfg.StaticDir, cfg.OutDir); err != nil {
		return err
	}
	for _, p := range rendered {
		if err := writeFile(filepath.Join(cfg.OutDir, p.file), p.data); err != nil {
			return err
		}
	}
	return CheckLinks(cfg.OutDir)
}

type renderedPage struct {
	file string
	data []byte
}

func renderAll(r *render.Renderer, s *content.Site) ([]renderedPage, error) {
	type spec struct {
		file   string
		render func() ([]byte, error)
	}
	specs := []spec{
		{"index.html", func() ([]byte, error) { return r.Home(s) }},
		{"blog/index.html", func() ([]byte, error) { return r.BlogIndex(s) }},
		{"404.html", func() ([]byte, error) { return r.NotFound(s) }},
	}
	urls := []string{"/", "/blog/"}
	for _, cs := range s.CaseStudies {
		specs = append(specs, spec{"case-studies/" + cs.Slug + "/index.html", func() ([]byte, error) { return r.CaseStudy(s, cs) }})
		urls = append(urls, "/case-studies/"+cs.Slug+"/")
	}
	for _, e := range s.Blog {
		if e.IsExternal() {
			continue
		}
		specs = append(specs, spec{"blog/" + e.Slug + "/index.html", func() ([]byte, error) { return r.Post(s, e) }})
		urls = append(urls, e.Href())
	}

	pages := make([]renderedPage, 0, len(specs)+2)
	for _, sp := range specs {
		data, err := sp.render()
		if err != nil {
			return nil, err
		}
		pages = append(pages, renderedPage{sp.file, data})
	}
	pages = append(pages, renderedPage{"sitemap.xml", render.Sitemap(s.Profile.SiteURL, urls)})
	pages = append(pages, renderedPage{"robots.txt", render.Robots(s.Profile.SiteURL)})
	return pages, nil
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

// checkOutDir resolves all Config paths to absolute form and refuses
// configurations where OutDir would overlap with a source directory or the
// current working directory.
func checkOutDir(cfg Config) error {
	out, err := filepath.Abs(cfg.OutDir)
	if err != nil {
		return fmt.Errorf("out dir: %w", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	sources := []struct {
		name string
		dir  string
	}{
		{"content", cfg.ContentDir},
		{"templates", cfg.TemplatesDir},
		{"static", cfg.StaticDir},
	}
	for _, src := range sources {
		abs, err := filepath.Abs(src.dir)
		if err != nil {
			return fmt.Errorf("%s dir: %w", src.name, err)
		}
		if within(abs, out) || within(out, abs) {
			return fmt.Errorf("out dir %q overlaps with %s dir %q", out, src.name, abs)
		}
	}
	// Refuse when OutDir equals or contains the working directory.
	if within(out, cwd) {
		return fmt.Errorf("out dir %q equals or contains the working directory %q", out, cwd)
	}
	return nil
}

// within reports whether child is equal to or nested inside parent.
// A child is outside only when filepath.Rel returns ".." (the parent itself is
// the child's parent) or a path that starts with "../" (a true ancestor
// escape). A relative path like "..foo" names a file whose name begins with
// ".." but lives inside the parent, so it counts as inside.
func within(parent, child string) bool {
	rel, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}
	sep := string(filepath.Separator)
	return rel != ".." && !strings.HasPrefix(rel, ".."+sep)
}
