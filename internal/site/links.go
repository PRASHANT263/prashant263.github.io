package site

import (
	"errors"
	"fmt"
	"html"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// linkAttr matches href or src attributes with double-quoted values.
// \b prevents matches like nohref (where the attribute name follows a word
// character) but does not exclude data-src or data-href, because the hyphen
// before the name is a non-word character and still triggers the boundary.
// Limits: double-quoted attributes only; srcset and absolute same-site URLs
// are not checked.
var linkAttr = regexp.MustCompile(`\b(?:href|src)="([^"]*)"`)

// cssURL matches root-relative url() references in CSS files.
var cssURL = regexp.MustCompile(`url\(["']?(/[^"')]+)["']?\)`)

// CheckLinks reports every root-relative href, src, or CSS url() in dir's
// HTML and CSS files that doesn't resolve to a file in dir. External, mailto
// and fragment-only links are skipped.
// Limits: HTML scanning covers double-quoted attributes only; srcset and
// absolute same-site URLs are not checked.
func CheckLinks(dir string) error {
	var broken []error
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(dir, path)
		ext := filepath.Ext(path)
		switch ext {
		case ".html":
			for _, m := range linkAttr.FindAllSubmatch(data, -1) {
				link := html.UnescapeString(string(m[1]))
				if !strings.HasPrefix(link, "/") || strings.HasPrefix(link, "//") {
					continue
				}
				if !resolves(dir, link) {
					broken = append(broken, fmt.Errorf("%s: broken link %s", rel, link))
				}
			}
		case ".css":
			for _, m := range cssURL.FindAllSubmatch(data, -1) {
				link := string(m[1])
				if strings.HasPrefix(link, "//") {
					continue
				}
				if !resolves(dir, link) {
					broken = append(broken, fmt.Errorf("%s: broken link %s", rel, link))
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return errors.Join(broken...)
}

func resolves(dir, link string) bool {
	if i := strings.IndexAny(link, "?#"); i >= 0 {
		link = link[:i]
	}
	p := filepath.Join(dir, filepath.FromSlash(link))
	info, err := os.Stat(p)
	if err != nil {
		return false
	}
	if info.IsDir() {
		_, err = os.Stat(filepath.Join(p, "index.html"))
		return err == nil
	}
	return true
}
