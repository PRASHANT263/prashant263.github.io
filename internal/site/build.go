// Package site builds the static site: it is the only package that touches
// the filesystem.
package site

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"portfolio/internal/content"
	"portfolio/internal/render"
)

// Build loads and validates content, then replaces OutDir with the rendered
// site plus static files, and checks internal links. Content and template
// errors are reported before OutDir is touched.
func Build(cfg Config) error {
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
	if err := resetDir(cfg.OutDir); err != nil {
		return err
	}
	if err := copyDir(cfg.StaticDir, cfg.OutDir); err != nil {
		return err
	}
	if err := writePages(r, s, cfg.OutDir); err != nil {
		return err
	}
	return CheckLinks(cfg.OutDir)
}

func writePages(r *render.Renderer, s *content.Site, out string) error {
	type page struct {
		file   string
		render func() ([]byte, error)
	}
	pages := []page{
		{"index.html", func() ([]byte, error) { return r.Home(s) }},
		{"blog/index.html", func() ([]byte, error) { return r.BlogIndex(s) }},
		{"404.html", func() ([]byte, error) { return r.NotFound(s) }},
	}
	urls := []string{"/", "/blog/"}
	for _, cs := range s.CaseStudies {
		cs := cs // capture loop variable
		pages = append(pages, page{"case-studies/" + cs.Slug + "/index.html", func() ([]byte, error) { return r.CaseStudy(s, cs) }})
		urls = append(urls, "/case-studies/"+cs.Slug+"/")
	}
	for _, e := range s.Blog {
		e := e // capture loop variable
		if e.IsExternal() {
			continue
		}
		pages = append(pages, page{"blog/" + e.Slug + "/index.html", func() ([]byte, error) { return r.Post(s, e) }})
		urls = append(urls, e.Href())
	}
	for _, p := range pages {
		html, err := p.render()
		if err != nil {
			return err
		}
		if err := writeFile(filepath.Join(out, p.file), html); err != nil {
			return err
		}
	}
	if err := writeFile(filepath.Join(out, "sitemap.xml"), render.Sitemap(s.Profile.SiteURL, urls)); err != nil {
		return err
	}
	return writeFile(filepath.Join(out, "robots.txt"), render.Robots(s.Profile.SiteURL))
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
